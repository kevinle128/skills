---
phase: 2
title: "Safe Skill Lifecycle"
status: completed
priority: P1
effort: 2d
dependencies: [1]
---

# Phase 2: Safe Skill Lifecycle

## Overview

Implement install, status, refresh, and uninstall through one manifest-backed sync engine.
The engine must never overwrite or delete an unknown or user-modified file without explicit force confirmation.

## Requirements

- Functional: install all embedded skills into both targets by default.
- Functional: support `--target agents`, `--target claude-code`, and `--target all`.
- Functional: record the SHA-256 of every file that KevinKit writes.
- Functional: remove stale owned files only when their current hash still matches the prior manifest.
- Functional: preserve modified and unowned files and report every preserved path.
- Functional: make `status` read-only and return a non-zero drift result when managed state differs.
- Non-functional: make an interrupted run safe to repeat.
- Non-functional: serialize lifecycle mutations with an atomic lock directory.

## Architecture

Store all lifecycle state under `~/.kevinkit` so users have one predictable location on every supported platform.
Use `~/.kevinkit/install-manifest.json` for ownership, `~/.kevinkit/backups/<timestamp>/` for recovery copies, and `~/.kevinkit/lifecycle.lock` for mutation serialization.
The manifest contains a schema version, kit version, target identifier, destination root, and sorted `{path, sha256}` entries.
One planner classifies every desired path as unchanged, create, replace-clean, conflict-modified, or conflict-unowned.
It also classifies prior owned paths that are absent from the new payload as remove-clean or preserve-modified.
The executor backs up every replaced or removed file, writes new content to a sibling temporary file, and uses rename for the final file switch.
The manifest is written last through the same temporary-file and rename pattern.
If a prior run stopped after a file switch, a rerun adopts a file whose current hash already equals the desired hash and converges without force.

## Related Code Files

- Create: `/Users/dale/Desktop/workspace/opensources/kevinle128-skills/internal/lifecycle/manifest.go`.
- Create: `/Users/dale/Desktop/workspace/opensources/kevinle128-skills/internal/lifecycle/plan.go`.
- Create: `/Users/dale/Desktop/workspace/opensources/kevinle128-skills/internal/lifecycle/apply.go`.
- Create: `/Users/dale/Desktop/workspace/opensources/kevinle128-skills/internal/lifecycle/status.go`.
- Create: `/Users/dale/Desktop/workspace/opensources/kevinle128-skills/internal/lifecycle/lifecycle_test.go`.
- Modify: `/Users/dale/Desktop/workspace/opensources/kevinle128-skills/internal/cli/cli.go`.

## Implementation Steps

1. Define manifest schema version 1 and reject unknown future schema versions without writing.
2. Resolve the target roots from the user home directory and refuse a path that escapes the selected skill root.
3. Acquire `~/.kevinkit/lifecycle.lock` and record enough owner data to diagnose a stale lock.
4. Build the complete change plan before writing and print it during `--dry-run`.
5. For normal install and update, skip modified or unowned conflicts and leave their manifest ownership unchanged or unclaimed.
6. For `--force`, require `--yes` in non-interactive use, back up the conflicting bytes, and then replace them.
7. Write new files with source mode bits limited to safe regular-file permissions.
8. Delete empty `kk-*` directories only after their owned files are removed, and never remove a directory that contains unknown content.
9. On uninstall, delete only files whose hash matches the manifest and release ownership of every preserved path.
10. Store recovery copies under `~/.kevinkit/backups/<timestamp>/` and print that location with the operation summary.

## Tests Before

- Create failing table tests for first install, idempotent install, clean update, stale removal, user edit, unowned collision, interrupted-run recovery, and uninstall.
- Run every case against temporary home and config directories.

## Tests After

- Verify that Codex and Devin cause one write to the shared `agents` target.
- Verify that a user-edited `SKILL.md` survives update and uninstall byte for byte.
- Verify that an unowned file inside a `kk-*` directory survives update and prevents unsafe directory removal.
- Verify that a failed apply keeps the prior manifest readable and that a rerun converges.
- Run `go test ./... -race` and `go vet ./...`.

## Success Criteria

- [x] Install is idempotent for both runtime roots.
- [x] Status identifies every missing, modified, and stale owned path without mutation.
- [x] Update changes only clean owned files unless force is explicitly confirmed.
- [x] Uninstall never deletes user-modified or unowned content.
- [x] Every destructive write has a recoverable backup.

## Risk Assessment

The critical risk is data loss caused by incorrect ownership classification.
Any test that observes a modified or unowned file change must stop implementation and require a redesign of the shared planner before proceeding.

## Security Considerations

Use `Lstat` and reject symlink traversal at every managed path component.
Create state, lock, temporary, and backup files with user-only permissions.
Do not trust absolute roots stored in a manifest without confirming that they match the current target resolution.
