package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/kevinle128/skills/internal/catalog"
	"github.com/kevinle128/skills/internal/lifecycle"
	"github.com/kevinle128/skills/internal/update"
)

const (
	ExitOK           = 0
	ExitError        = 1
	ExitUsage        = 2
	ExitConfirmation = 3
)

type Config struct {
	Source        fs.FS
	SkillRoot     string
	Version       string
	Commit        string
	BuildDate     string
	ReleaseAPIURL string
	Stdout        io.Writer
	Stderr        io.Writer
	HomeDir       func() (string, error)
	Executable    func() (string, error)
	Env           []string
}

type App struct {
	config  Config
	catalog catalog.Catalog
}

func New(config Config) (*App, error) {
	payload, err := catalog.Load(config.Source, config.SkillRoot)
	if err != nil {
		return nil, err
	}
	if config.Version == "" {
		config.Version = "dev"
	}
	if config.Stdout == nil {
		config.Stdout = io.Discard
	}
	if config.Stderr == nil {
		config.Stderr = io.Discard
	}
	if config.HomeDir == nil {
		config.HomeDir = os.UserHomeDir
	}
	if config.Executable == nil {
		config.Executable = os.Executable
	}
	if config.Env == nil {
		config.Env = os.Environ()
	}
	return &App{config: config, catalog: payload}, nil
}

func (a *App) Run(ctx context.Context, args []string) int {
	if len(args) == 0 {
		a.printRootHelp(a.config.Stdout)
		return ExitOK
	}
	if len(args) == 2 && validCommand(args[0]) && helpOnly(args[1:]) {
		a.printCommandHelp(a.config.Stdout, args[0])
		return ExitOK
	}
	switch args[0] {
	case "help":
		return a.help(args[1:])
	case "-h", "--help":
		a.printRootHelp(a.config.Stdout)
		return ExitOK
	case "version":
		return a.version(args[1:])
	case "install":
		return a.install(args[1:])
	case "status":
		return a.status(args[1:])
	case "uninstall":
		return a.uninstall(args[1:])
	case "update":
		return a.update(ctx, args[1:])
	default:
		fmt.Fprintf(a.config.Stderr, "error: unknown command %q\n\n", args[0])
		a.printRootHelp(a.config.Stderr)
		return ExitUsage
	}
}

func (a *App) help(args []string) int {
	if len(args) == 0 {
		a.printRootHelp(a.config.Stdout)
		return ExitOK
	}
	if len(args) > 1 || !validCommand(args[0]) {
		fmt.Fprintln(a.config.Stderr, "error: usage: kk help [command]")
		return ExitUsage
	}
	a.printCommandHelp(a.config.Stdout, args[0])
	return ExitOK
}

func (a *App) version(args []string) int {
	if helpOnly(args) {
		a.printCommandHelp(a.config.Stdout, "version")
		return ExitOK
	}
	if len(args) != 0 {
		return a.usageError("version", fmt.Errorf("version does not accept arguments"))
	}
	fmt.Fprintf(a.config.Stdout, "kk %s\ncommit: %s\nbuilt: %s\nembedded kit: %s\n", a.config.Version, valueOrUnknown(a.config.Commit), valueOrUnknown(a.config.BuildDate), a.config.Version)
	return ExitOK
}

func (a *App) install(args []string) int {
	options, code := a.parseMutationFlags("install", args, true)
	if code != ExitOK {
		return code
	}
	if options.help {
		a.printCommandHelp(a.config.Stdout, "install")
		return ExitOK
	}
	if options.force && !options.dryRun && !options.yes {
		return a.confirmationError("install --force")
	}
	manager, err := a.manager()
	if err != nil {
		return a.runtimeError(err)
	}
	result, err := manager.Sync(lifecycle.Options{
		Target:    options.target,
		Force:     options.force,
		DryRun:    options.dryRun,
		LockToken: envValue(a.config.Env, "KEVINKIT_LOCK_TOKEN"),
	})
	if err != nil {
		return a.runtimeError(err)
	}
	a.printChanges(result, options.dryRun)
	return ExitOK
}

func (a *App) status(args []string) int {
	options, code := a.parseTargetFlags("status", args)
	if code != ExitOK {
		return code
	}
	if options.help {
		a.printCommandHelp(a.config.Stdout, "status")
		return ExitOK
	}
	manager, err := a.manager()
	if err != nil {
		return a.runtimeError(err)
	}
	result, err := manager.Status(options.target)
	if err != nil {
		return a.runtimeError(err)
	}
	fmt.Fprintf(a.config.Stdout, "KevinKit %s\n", a.config.Version)
	for _, target := range result.Targets {
		if target.Installed {
			fmt.Fprintf(a.config.Stdout, "%s: installed (%d managed files)\n", target.ID, target.Owned)
		} else {
			fmt.Fprintf(a.config.Stdout, "%s: not installed\n", target.ID)
		}
	}
	if result.Clean() {
		fmt.Fprintln(a.config.Stdout, "status: clean")
		return ExitOK
	}
	fmt.Fprintf(a.config.Stdout, "status: drift (%d)\n", len(result.Drift))
	for _, drift := range result.Drift {
		if drift.Path == "" {
			fmt.Fprintf(a.config.Stdout, "- %s: %s\n", drift.Target, drift.Kind)
		} else {
			fmt.Fprintf(a.config.Stdout, "- %s:%s: %s\n", drift.Target, drift.Path, drift.Kind)
		}
	}
	return ExitError
}

func (a *App) uninstall(args []string) int {
	options, code := a.parseMutationFlags("uninstall", args, false)
	if code != ExitOK {
		return code
	}
	if options.help {
		a.printCommandHelp(a.config.Stdout, "uninstall")
		return ExitOK
	}
	if !options.dryRun && !options.yes {
		return a.confirmationError("uninstall")
	}
	manager, err := a.manager()
	if err != nil {
		return a.runtimeError(err)
	}
	result, err := manager.Uninstall(lifecycle.Options{
		Target:    options.target,
		DryRun:    options.dryRun,
		LockToken: envValue(a.config.Env, "KEVINKIT_LOCK_TOKEN"),
	})
	if err != nil {
		return a.runtimeError(err)
	}
	a.printChanges(result, options.dryRun)
	return ExitOK
}

func (a *App) update(ctx context.Context, args []string) int {
	options, code := a.parseUpdateFlags(args)
	if code != ExitOK {
		return code
	}
	if options.help {
		a.printCommandHelp(a.config.Stdout, "update")
		return ExitOK
	}
	if err := validateTarget(options.target); err != nil {
		return a.usageError("update", err)
	}
	executable, err := a.config.Executable()
	if err != nil {
		return a.runtimeError(fmt.Errorf("locate current executable: %w", err))
	}
	if resolved, resolveErr := filepath.EvalSymlinks(executable); resolveErr == nil {
		executable = resolved
	}
	home, err := a.config.HomeDir()
	if err != nil {
		return a.runtimeError(fmt.Errorf("resolve user home: %w", err))
	}
	updater, err := update.New(update.Config{
		APIURL:         a.config.ReleaseAPIURL,
		CurrentVersion: a.config.Version,
		Executable:     executable,
		StateDir:       filepath.Join(home, ".kevinkit"),
		Target:         options.target,
		Force:          options.force,
		Stdout:         a.config.Stdout,
		Stderr:         a.config.Stderr,
		Env:            a.config.Env,
	})
	if err != nil {
		return a.runtimeError(err)
	}
	release, available, err := updater.Check(ctx)
	if err != nil {
		return a.runtimeError(err)
	}
	if !available {
		fmt.Fprintf(a.config.Stdout, "KevinKit %s is current.\n", a.config.Version)
		return ExitOK
	}
	if options.check {
		fmt.Fprintf(a.config.Stdout, "Update available: %s -> %s (%s)\n", a.config.Version, release.Version, release.Archive.Name)
		return ExitOK
	}
	if options.dryRun {
		fmt.Fprintf(a.config.Stdout, "Would update: %s -> %s (%s)\n", a.config.Version, release.Version, release.Archive.Name)
		return ExitOK
	}
	if !options.yes {
		fmt.Fprintf(a.config.Stdout, "Update available: %s -> %s (%s)\n", a.config.Version, release.Version, release.Archive.Name)
		return a.confirmationError("update")
	}
	manager, err := a.manager()
	if err != nil {
		return a.runtimeError(err)
	}
	lock, err := manager.AcquireLock("")
	if err != nil {
		return a.runtimeError(err)
	}
	updater, err = update.New(update.Config{
		APIURL:         a.config.ReleaseAPIURL,
		CurrentVersion: a.config.Version,
		Executable:     executable,
		StateDir:       manager.StateDir,
		Target:         options.target,
		Force:          options.force,
		LockToken:      lock.Token(),
		Stdout:         a.config.Stdout,
		Stderr:         a.config.Stderr,
		Env:            a.config.Env,
	})
	if err != nil {
		_ = lock.Release()
		return a.runtimeError(err)
	}
	result, err := updater.Apply(ctx, release)
	releaseErr := lock.Release()
	if err != nil {
		if releaseErr != nil {
			return a.runtimeError(fmt.Errorf("%v; release lifecycle lock: %w", err, releaseErr))
		}
		return a.runtimeError(err)
	}
	if releaseErr != nil {
		return a.runtimeError(releaseErr)
	}
	fmt.Fprintf(a.config.Stdout, "KevinKit updated to %s.\nbackup: %s\n", result.Version, result.Backup)
	return ExitOK
}

type targetFlags struct {
	target string
	help   bool
}

type mutationFlags struct {
	target string
	dryRun bool
	force  bool
	yes    bool
	help   bool
}

type updateFlags struct {
	mutationFlags
	check bool
}

func (a *App) parseTargetFlags(command string, args []string) (targetFlags, int) {
	var options targetFlags
	set := newFlagSet(command)
	set.StringVar(&options.target, "target", "all", "target: all, agents, or claude-code")
	set.BoolVar(&options.help, "help", false, "show help")
	if err := set.Parse(args); err != nil {
		return options, a.usageError(command, err)
	}
	if set.NArg() != 0 {
		return options, a.usageError(command, fmt.Errorf("unexpected argument %q", set.Arg(0)))
	}
	if !options.help {
		if err := validateTarget(options.target); err != nil {
			return options, a.usageError(command, err)
		}
	}
	return options, ExitOK
}

func (a *App) parseMutationFlags(command string, args []string, allowForce bool) (mutationFlags, int) {
	var options mutationFlags
	set := newFlagSet(command)
	set.StringVar(&options.target, "target", "all", "target: all, agents, or claude-code")
	set.BoolVar(&options.dryRun, "dry-run", false, "show changes without writing")
	if allowForce {
		set.BoolVar(&options.force, "force", false, "replace conflicting files after backup")
	}
	set.BoolVar(&options.yes, "yes", false, "confirm the operation")
	set.BoolVar(&options.help, "help", false, "show help")
	if err := set.Parse(args); err != nil {
		return options, a.usageError(command, err)
	}
	if set.NArg() != 0 {
		return options, a.usageError(command, fmt.Errorf("unexpected argument %q", set.Arg(0)))
	}
	if !options.help {
		if err := validateTarget(options.target); err != nil {
			return options, a.usageError(command, err)
		}
	}
	return options, ExitOK
}

func (a *App) parseUpdateFlags(args []string) (updateFlags, int) {
	var options updateFlags
	set := newFlagSet("update")
	set.StringVar(&options.target, "target", "all", "target: all, agents, or claude-code")
	set.BoolVar(&options.dryRun, "dry-run", false, "show the release without writing")
	set.BoolVar(&options.force, "force", false, "replace conflicting skill files after backup")
	set.BoolVar(&options.yes, "yes", false, "confirm the update")
	set.BoolVar(&options.check, "check", false, "check for a new release")
	set.BoolVar(&options.help, "help", false, "show help")
	if err := set.Parse(args); err != nil {
		return options, a.usageError("update", err)
	}
	if set.NArg() != 0 {
		return options, a.usageError("update", fmt.Errorf("unexpected argument %q", set.Arg(0)))
	}
	if options.check && options.dryRun {
		return options, a.usageError("update", fmt.Errorf("--check and --dry-run cannot be used together"))
	}
	return options, ExitOK
}

func newFlagSet(name string) *flag.FlagSet {
	set := flag.NewFlagSet(name, flag.ContinueOnError)
	set.SetOutput(io.Discard)
	return set
}

func (a *App) manager() (*lifecycle.Manager, error) {
	home, err := a.config.HomeDir()
	if err != nil {
		return nil, fmt.Errorf("resolve user home: %w", err)
	}
	return lifecycle.NewManager(a.catalog, a.config.Version, home)
}

func (a *App) printChanges(result lifecycle.Result, dryRun bool) {
	label := "Applied"
	if dryRun {
		label = "Planned"
	}
	fmt.Fprintf(a.config.Stdout, "%s changes: create=%d update=%d remove=%d preserve=%d unchanged=%d\n",
		label,
		result.Count(lifecycle.ActionCreate),
		result.Count(lifecycle.ActionUpdate),
		result.Count(lifecycle.ActionRemove),
		result.Count(lifecycle.ActionPreserve),
		result.Count(lifecycle.ActionUnchanged),
	)
	for _, change := range result.Changes {
		if dryRun || change.Action == lifecycle.ActionPreserve {
			fmt.Fprintf(a.config.Stdout, "- %s %s:%s (%s)\n", change.Action, change.Target, change.Path, change.Reason)
		}
	}
	if result.BackupDir != "" {
		fmt.Fprintf(a.config.Stdout, "backup: %s\n", result.BackupDir)
	}
}

func (a *App) runtimeError(err error) int {
	fmt.Fprintf(a.config.Stderr, "error: %v\n", err)
	return ExitError
}

func (a *App) usageError(command string, err error) int {
	fmt.Fprintf(a.config.Stderr, "error: %v\n\n", err)
	a.printCommandHelp(a.config.Stderr, command)
	return ExitUsage
}

func (a *App) confirmationError(operation string) int {
	fmt.Fprintf(a.config.Stderr, "confirmation required: rerun %s with --yes\n", operation)
	return ExitConfirmation
}

func (a *App) printRootHelp(writer io.Writer) {
	fmt.Fprintln(writer, `KevinKit manages installed copies of the embedded KevinKit skills.

Usage:
  kk <command> [flags]

Commands:
  install     Install or synchronize embedded skills
  update      Update kk from a verified GitHub Release
  status      Report installed targets and drift
  uninstall   Remove clean managed files
  version     Print build and embedded kit versions
  help        Show help for a command`)
}

func (a *App) printCommandHelp(writer io.Writer, command string) {
	usage := map[string]string{
		"install":   "kk install [--target all|agents|claude-code] [--dry-run] [--force --yes]",
		"update":    "kk update [--target all|agents|claude-code] [--check|--dry-run] [--force] [--yes]",
		"status":    "kk status [--target all|agents|claude-code]",
		"uninstall": "kk uninstall [--target all|agents|claude-code] [--dry-run] [--yes]",
		"version":   "kk version",
	}
	fmt.Fprintf(writer, "Usage:\n  %s\n", usage[command])
}

func validateTarget(target string) error {
	switch target {
	case "all", "agents", "claude-code":
		return nil
	default:
		return fmt.Errorf("unknown target %q; want all, agents, or claude-code", target)
	}
}

func validCommand(command string) bool {
	switch command {
	case "install", "status", "uninstall", "update", "version":
		return true
	default:
		return false
	}
}

func helpOnly(args []string) bool {
	return len(args) == 1 && (args[0] == "-h" || args[0] == "--help")
}

func valueOrUnknown(value string) string {
	if value == "" {
		return "unknown"
	}
	return value
}

func envValue(environment []string, key string) string {
	prefix := key + "="
	for index := len(environment) - 1; index >= 0; index-- {
		if strings.HasPrefix(environment[index], prefix) {
			return strings.TrimPrefix(environment[index], prefix)
		}
	}
	return ""
}
