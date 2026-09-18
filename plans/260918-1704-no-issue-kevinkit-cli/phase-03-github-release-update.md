---
phase: 3
title: "GitHub Release Update"
status: completed
priority: P1
effort: 1.5d
dependencies: [1, 2]
---

# Phase 3: GitHub Release Update

## Overview

Make `kk update` fetch, verify, and activate the latest public GitHub Release for the current platform.
The downloaded binary is never executed or installed before its SHA-256 value matches the published checksum file.

## Requirements

- Functional: query the latest release for `kevinle128/skills` over HTTPS.
- Functional: select the asset for the current `GOOS` and `GOARCH`.
- Functional: verify the archive against `checksums.txt` before extraction.
- Functional: replace the current executable with the verified binary, then run that new binary to synchronize its embedded skills.
- Functional: support `kk update --check` and `kk update --dry-run` without mutation.
- Non-functional: use timeouts, bounded downloads, and explicit GitHub API headers.
- Non-functional: retain the prior executable until the replacement succeeds.

## Architecture

`internal/update` owns release discovery, download, checksum parsing, archive extraction, and executable replacement.
GitHub HTTPS and the release checksum file are the initial trust root.
The updater accepts internal HTTP and filesystem interfaces so tests use a local server and temporary executable path.
Release asset names follow one fixed contract for Darwin and Linux on `amd64` and `arm64`.
The updater backs up and replaces the old executable before it launches the new binary to synchronize the embedded skills.
If skill synchronization fails, the verified new binary remains runnable, the old binary remains in backup, and a repeated install can converge safely.

## Related Code Files

- Create: `/Users/dale/Desktop/workspace/opensources/kevinle128-skills/internal/update/update.go`.
- Create: `/Users/dale/Desktop/workspace/opensources/kevinle128-skills/internal/update/replace_unix.go`.
- Create: `/Users/dale/Desktop/workspace/opensources/kevinle128-skills/internal/update/replace_windows.go`.
- Create: `/Users/dale/Desktop/workspace/opensources/kevinle128-skills/internal/update/update_test.go`.
- Create: `/Users/dale/Desktop/workspace/opensources/kevinle128-skills/.goreleaser.yaml`.
- Create: `/Users/dale/Desktop/workspace/opensources/kevinle128-skills/.github/workflows/release.yml`.
- Modify: `/Users/dale/Desktop/workspace/opensources/kevinle128-skills/internal/cli/cli.go`.

## Implementation Steps

1. Define the release API URL, repository, user agent, asset pattern, maximum response size, and request timeout in one update package.
2. Parse only the release fields and assets required for update, and reject draft or prerelease data for the default stable channel.
3. Return without download when the release tag equals the running version.
4. Download the platform archive and `checksums.txt` into a private temporary directory.
5. Parse one exact checksum entry, compute SHA-256 while reading the archive, and reject missing, duplicate, or mismatched entries.
6. Extract only the expected `kk` or `kk.exe` entry and reject absolute paths, parent traversal, links, devices, and extra executable candidates.
7. Confirm that the current executable directory is writable before any installed state changes.
8. Back up the running executable, atomically place the new executable where the current one resides, and restore the backup on replacement failure.
9. Execute the installed new binary with an internal non-interactive install command that uses the phase 2 lifecycle engine.
10. Make the platform replacement behavior explicit in small build-tagged files only where the operating systems require different calls.
11. Configure GoReleaser to build supported platform archives and one `checksums.txt` asset from a version tag.

## Tests Before

- Add local HTTP server tests for no update, unknown platform, missing asset, timeout, oversized body, bad checksum, duplicate checksum, and path traversal.
- Add replacement tests for success, permission failure, and restoration after rename failure.

## Tests After

- Verify that corrupt bytes are never executed.
- Verify that `--check` and `--dry-run` do not change the binary, manifest, or installed skills.
- Verify that the installed new binary performs the skill sync before final success is reported.
- Run `go test ./... -race`, `go vet ./...`, and a local GoReleaser snapshot build.

## Success Criteria

- [x] `kk update` installs only a release archive that matches `checksums.txt`.
- [x] Network and replacement failures leave a runnable CLI and clear recovery instructions.
- [x] The new embedded skills are synchronized through the same safe lifecycle engine as normal install.
- [x] Release archives have stable names that the updater can resolve without heuristics.

## Risk Assessment

The critical risks are supply-chain substitution and failure during self-replacement.
A checksum mismatch, unexpected archive member, or failed rollback must fail closed and must not report a successful update.
If cross-platform self-replacement cannot be proven in CI, stop and narrow the published platform list instead of shipping an unverified path.
The first release therefore publishes only the four Darwin and Linux `amd64` and `arm64` targets that run the native release gate.

## Security Considerations

Never execute unverified downloaded bytes.
Do not pass tokens because the repository and releases are public.
Do not log response bodies, local home paths, or temporary binary contents in normal output.
