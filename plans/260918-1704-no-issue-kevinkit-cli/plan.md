---
title: "KevinKit CLI"
description: "Add an open-source kk CLI that safely installs and updates KevinKit skill copies from verified GitHub Releases."
status: completed
priority: P1
effort: 5d
issue: null
branch: main
tags: [feature, cli, release, infra]
blockedBy: []
blocks: []
created: 2026-09-18
---

# KevinKit CLI

## Overview

Build a public `kk` Go binary that embeds the KevinKit skills and manages their installed copies for Claude Code, Codex, and Devin.
The CLI must preserve user changes, recover from interrupted writes, and update itself and the embedded skills from GitHub Releases without any license, login, entitlement, or paid registry layer.

## Accepted Decisions

- GitHub Releases are the distribution source and trust root.
- The Go binary embeds every `skills/kk-*` file with `go:embed`.
- Codex and Devin share the `agents` target at `~/.agents/skills`.
- Claude Code uses the `claude-code` target at `~/.claude/skills`.
- KevinKit stores only lifecycle state under `~/.kevinkit`.
- The CLI uses only the Go standard library at runtime.
- A SHA-256 manifest defines KevinKit ownership and protects user-modified files.

## Command Contract

| Command | Result |
|---|---|
| `kk install` | Install the embedded skills into all selected targets. |
| `kk update` | Verify the latest GitHub Release, update the binary, and synchronize its embedded skills. |
| `kk status` | Report version, installed targets, missing files, modified files, and stale owned files without writing. |
| `kk uninstall` | Remove only files that still match the install manifest. |
| `kk version` | Print the CLI and embedded kit versions. |

## Non-Goals

- No license checks, accounts, remote registry, analytics, TUI, plugin mode, or database.
- No project-local installation mode in the first release.
- No automatic edit of agent configuration files.

## Architecture

`GitHub Release -> verified platform binary -> embedded skill catalog -> safe sync engine -> runtime skill roots -> ~/.kevinkit manifest and backups`.

## Phases

| # | Phase | Status |
|---|---|---|
| 1 | [CLI Core and Embedded Payload](./phase-01-start.md) | Completed |
| 2 | [Safe Skill Lifecycle](./phase-02-safe-skill-lifecycle.md) | Completed |
| 3 | [GitHub Release Update](./phase-03-github-release-update.md) | Completed |
| 4 | [End-to-End Release Proof](./phase-04-end-to-end-release-proof.md) | Completed |

## Dependencies

- Go 1.24 or later for builds.
- GitHub Actions and GitHub Releases for public artifacts.
- Existing `skills/kk-*` directories as the canonical payload.
- [Xia source and decision report](./reports/xia-analysis.md).

## Success Criteria

- [x] A new user can install `kk`, run `kk install`, and see all 34 skills in both runtime roots.
- [x] `kk update` verifies the release checksum before it executes or replaces downloaded bytes.
- [x] Clean owned files update and stale clean owned files disappear.
- [x] Modified and unknown files remain unchanged unless the user explicitly passes `--force`.
- [x] An interrupted operation is safe to rerun and converges to one consistent manifest.
- [x] The full E2E test runs the CLI as a user from install through update, status, and uninstall.
- [x] Release CI tests the real binary and produces deterministic platform archives plus `checksums.txt`.

<!-- slug: no-issue-kevinkit-cli -->
