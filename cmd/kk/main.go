package main

import (
	"context"
	"fmt"
	"os"

	kevinkit "github.com/kevinle128/skills"
	"github.com/kevinle128/skills/internal/buildinfo"
	"github.com/kevinle128/skills/internal/cli"
)

func main() {
	app, err := cli.New(cli.Config{
		Source:        kevinkit.SkillFS,
		SkillRoot:     kevinkit.SkillRoot,
		Version:       buildinfo.Version,
		Commit:        buildinfo.Commit,
		BuildDate:     buildinfo.BuildDate,
		ReleaseAPIURL: buildinfo.ReleaseAPIURL,
		Stdout:        os.Stdout,
		Stderr:        os.Stderr,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: load embedded skills: %v\n", err)
		os.Exit(cli.ExitError)
	}
	os.Exit(app.Run(context.Background(), os.Args[1:]))
}
