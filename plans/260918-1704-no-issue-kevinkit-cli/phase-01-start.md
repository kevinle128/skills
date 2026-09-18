---
phase: 1
title: "CLI Core and Embedded Payload"
status: completed
priority: P1
effort: 1d
dependencies: []
---

# Phase 1: CLI Core and Embedded Payload

## Overview

Create a small Go command that can inventory the canonical KevinKit skills from an embedded filesystem.
This phase defines the stable command, version, target, and payload contracts used by every later phase.

## Requirements

- Functional: expose `install`, `update`, `status`, `uninstall`, `version`, and help entry points.
- Functional: embed all files under `skills/kk-*`, including dotfiles.
- Functional: map `agents` to `~/.agents/skills` and `claude-code` to `~/.claude/skills`.
- Non-functional: use the Go standard library for runtime code.
- Non-functional: reject invalid embedded paths before any disk write.

## Architecture

The root `kevinkit` package owns the `go:embed` declaration because Go patterns cannot traverse parent directories.
Use `//go:embed all:skills/kk-*` so the payload includes dotfiles.
`cmd/kk` delegates argument parsing to `internal/cli`, and `internal/catalog` converts the embedded tree into a sorted list of skill roots and regular files.
Build variables provide the CLI version, commit, and build date without changing source files during release.

## Related Code Files

- Create: `/Users/dale/Desktop/workspace/opensources/kevinle128-skills/go.mod`.
- Create: `/Users/dale/Desktop/workspace/opensources/kevinle128-skills/content.go`.
- Create: `/Users/dale/Desktop/workspace/opensources/kevinle128-skills/content_test.go`.
- Create: `/Users/dale/Desktop/workspace/opensources/kevinle128-skills/cmd/kk/main.go`.
- Create: `/Users/dale/Desktop/workspace/opensources/kevinle128-skills/internal/cli/cli.go`.
- Create: `/Users/dale/Desktop/workspace/opensources/kevinle128-skills/internal/cli/cli_test.go`.
- Create: `/Users/dale/Desktop/workspace/opensources/kevinle128-skills/internal/catalog/catalog.go`.
- Create: `/Users/dale/Desktop/workspace/opensources/kevinle128-skills/internal/catalog/catalog_test.go`.

## Implementation Steps

1. Initialize module `github.com/kevinle128/skills` and set the minimum supported Go version.
2. Embed `all:skills/kk-*` in `content.go` and expose it as a read-only `fs.FS`.
3. Walk the embedded tree in lexical order and require each top-level directory to match `kk-*` and contain `SKILL.md`.
4. Reject absolute paths, `..`, non-regular payload entries, and duplicate relative paths.
5. Add a standard-library subcommand parser with shared `--target`, `--dry-run`, `--force`, `--yes`, and `--help` handling only where the command supports them.
6. Keep stdout for results and stderr for diagnostics so tests and scripts can distinguish them.
7. Add build variables and make `kk version` report both CLI and embedded kit versions.

## Tests Before

- Add a catalog test that fails when dotfiles are omitted from the embedded payload.
- Add parser table tests for valid commands, unknown commands, missing flag values, and conflicting flags.

## Tests After

- Assert that the catalog contains exactly the source skill directories and files.
- Assert that `kk version` and help do not read or write the user home directory.
- Run `go test ./...` and `go vet ./...`.

## Success Criteria

- [x] `go build ./cmd/kk` produces one binary with no runtime package dependency.
- [x] The embedded catalog matches all canonical skill files, including hidden files.
- [x] Every public command has deterministic help and exit behavior.

## Risk Assessment

The main risk is silent payload omission by `go:embed`.
The catalog parity test is the observable signal and must block the build when the embedded file set differs from `skills/kk-*`.

## Security Considerations

Treat every destination path as untrusted input even though the embedded source is trusted.
Do not follow symlinks while resolving or writing a managed path.
