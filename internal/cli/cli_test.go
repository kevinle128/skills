package cli

import (
	"bytes"
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
)

func TestUserLifecycleThroughCLI(t *testing.T) {
	home := t.TempDir()
	app, stdout, stderr := testApp(t, home)

	if code := app.Run(context.Background(), []string{"install", "--dry-run"}); code != ExitOK {
		t.Fatalf("dry-run exit = %d, stderr = %s", code, stderr.String())
	}
	assertMissing(t, filepath.Join(home, ".agents"))
	if !strings.Contains(stdout.String(), "Planned changes: create=4") {
		t.Fatalf("unexpected dry-run output: %s", stdout.String())
	}

	stdout.Reset()
	stderr.Reset()
	if code := app.Run(context.Background(), []string{"install", "--yes"}); code != ExitOK {
		t.Fatalf("install exit = %d, stderr = %s", code, stderr.String())
	}
	agentsSkill := filepath.Join(home, ".agents", "skills", "kk-fixture", "SKILL.md")
	claudeSkill := filepath.Join(home, ".claude", "skills", "kk-fixture", "SKILL.md")
	assertFile(t, agentsSkill, "fixture")
	assertFile(t, claudeSkill, "fixture")

	stdout.Reset()
	if code := app.Run(context.Background(), []string{"status"}); code != ExitOK {
		t.Fatalf("clean status exit = %d, output = %s", code, stdout.String())
	}
	if !strings.Contains(stdout.String(), "status: clean") {
		t.Fatalf("unexpected clean status: %s", stdout.String())
	}

	if err := os.WriteFile(agentsSkill, []byte("user edit"), 0o644); err != nil {
		t.Fatal(err)
	}
	stdout.Reset()
	if code := app.Run(context.Background(), []string{"status"}); code != ExitError {
		t.Fatalf("drift status exit = %d, output = %s", code, stdout.String())
	}
	if !strings.Contains(stdout.String(), "agents:kk-fixture/SKILL.md: modified") {
		t.Fatalf("modified file missing from status: %s", stdout.String())
	}

	stderr.Reset()
	if code := app.Run(context.Background(), []string{"install", "--force"}); code != ExitConfirmation {
		t.Fatalf("force without confirmation exit = %d", code)
	}
	assertFile(t, agentsSkill, "user edit")
	if code := app.Run(context.Background(), []string{"install", "--force", "--yes"}); code != ExitOK {
		t.Fatalf("confirmed force exit = %d, stderr = %s", code, stderr.String())
	}
	assertFile(t, agentsSkill, "fixture")

	if code := app.Run(context.Background(), []string{"uninstall"}); code != ExitConfirmation {
		t.Fatalf("unconfirmed uninstall exit = %d", code)
	}
	if code := app.Run(context.Background(), []string{"uninstall", "--yes"}); code != ExitOK {
		t.Fatalf("uninstall exit = %d, stderr = %s", code, stderr.String())
	}
	assertMissing(t, agentsSkill)
	assertMissing(t, claudeSkill)
}

func TestHelpAndVersionDoNotResolveHome(t *testing.T) {
	called := false
	stdout := &bytes.Buffer{}
	app, err := New(Config{
		Source:        testSkillFS(),
		SkillRoot:     "skills",
		Version:       "1.2.3",
		Commit:        "abc123",
		BuildDate:     "2026-09-18",
		ReleaseAPIURL: "https://example.test/latest",
		Stdout:        stdout,
		HomeDir: func() (string, error) {
			called = true
			return "", errors.New("must not be called")
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"--help"}, {"install", "--help"}, {"version"}} {
		stdout.Reset()
		if code := app.Run(context.Background(), args); code != ExitOK {
			t.Fatalf("%v exit = %d", args, code)
		}
	}
	if called {
		t.Fatal("help or version resolved the user home")
	}
	if !strings.Contains(stdout.String(), "kk 1.2.3") {
		t.Fatalf("unexpected version output: %s", stdout.String())
	}
}

func TestUsageErrorsAreDeterministic(t *testing.T) {
	app, _, stderr := testApp(t, t.TempDir())
	tests := []struct {
		args []string
		want string
	}{
		{args: []string{"unknown"}, want: "unknown command"},
		{args: []string{"status", "--target", "other"}, want: "unknown target"},
		{args: []string{"install", "--target"}, want: "flag needs an argument"},
		{args: []string{"uninstall", "--force"}, want: "flag provided but not defined"},
		{args: []string{"update", "--check", "--dry-run"}, want: "cannot be used together"},
	}
	for _, test := range tests {
		stderr.Reset()
		if code := app.Run(context.Background(), test.args); code != ExitUsage {
			t.Fatalf("%v exit = %d", test.args, code)
		}
		if !strings.Contains(stderr.String(), test.want) {
			t.Fatalf("%v stderr = %q, want %q", test.args, stderr.String(), test.want)
		}
	}
}

func testApp(t *testing.T, home string) (*App, *bytes.Buffer, *bytes.Buffer) {
	t.Helper()
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	app, err := New(Config{
		Source:        testSkillFS(),
		SkillRoot:     "skills",
		Version:       "1.0.0",
		Commit:        "test",
		BuildDate:     "2026-09-18",
		ReleaseAPIURL: "https://example.test/latest",
		Stdout:        stdout,
		Stderr:        stderr,
		HomeDir:       func() (string, error) { return home, nil },
		Executable:    func() (string, error) { return filepath.Join(home, "bin", "kk"), nil },
		Env:           []string{"HOME=" + home},
	})
	if err != nil {
		t.Fatal(err)
	}
	return app, stdout, stderr
}

func testSkillFS() fs.FS {
	return fstest.MapFS{
		"skills/kk-fixture/SKILL.md": &fstest.MapFile{Data: []byte("fixture")},
		"skills/kk-fixture/file.txt": &fstest.MapFile{Data: []byte("content")},
	}
}

func assertFile(t *testing.T, name, want string) {
	t.Helper()
	data, err := os.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != want {
		t.Fatalf("%s = %q, want %q", name, data, want)
	}
}

func assertMissing(t *testing.T, name string) {
	t.Helper()
	if _, err := os.Lstat(name); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("%s exists or returned %v", name, err)
	}
}
