---
phase: 4
title: "End-to-End Release Proof"
status: completed
priority: P1
effort: 0.5d
dependencies: [1, 2, 3]
---

# Phase 4: End-to-End Release Proof

## Overview

Prove the complete user workflow with real compiled binaries, isolated home directories, a local release server, and the real runtime target paths.
Update the public documentation only after this proof passes.

## Requirements

- Functional: drive install, status, update, and uninstall through the `kk` process boundary.
- Functional: prove clean update, stale removal, modified-file preservation, unknown-file preservation, and binary version replacement in one flow.
- Functional: verify both `~/.agents/skills` and `~/.claude/skills` from the user point of view.
- Non-functional: run the E2E flow in the release workflow on every published operating system and architecture that GitHub supports directly.
- Non-functional: fail the release when unit, race, vet, E2E, archive, or checksum gates fail.

## Architecture

The E2E test creates a temporary user home and starts a local HTTP server that implements the small GitHub release response used by `kk update`.
It builds an old binary from payload A and a new release binary from payload B, then serves the new archive and checksum.
The test invokes the binaries with `os/exec` and inspects only public output, target directories, manifest state, and process exit codes.
The release workflow uses the same build and test commands before it publishes any asset.

## Related Code Files

- Create: `/Users/dale/Desktop/workspace/opensources/kevinle128-skills/internal/e2e/lifecycle_test.go`.
- Create: `/Users/dale/Desktop/workspace/opensources/kevinle128-skills/internal/e2e/testdata/README.md` if fixture generation needs a durable explanation.
- Modify: `/Users/dale/Desktop/workspace/opensources/kevinle128-skills/.github/workflows/release.yml`.
- Modify: `/Users/dale/Desktop/workspace/opensources/kevinle128-skills/README.md`.
- Modify: `/Users/dale/Desktop/workspace/opensources/kevinle128-skills/docs/guide/README.md`.
- Modify: `/Users/dale/Desktop/workspace/opensources/kevinle128-skills/docs/guide/07-skill-catalog.md` only if installation guidance belongs there after navigation review.

## Implementation Steps

1. Build payload A with one file that changes, one clean file that is later removed, and one file that remains stable.
2. Run the old binary as a user with temporary `HOME` and platform config variables.
3. Assert that all embedded skills appear in the `agents` and `claude-code` roots and that `kk status` is clean.
4. Modify one installed owned file and add one unknown file under an installed skill directory.
5. Build payload B and serve its platform archive, release JSON, and checksum from the local HTTP server.
6. Run `kk update` through the old binary and assert that the process activates payload B and reports all preserved conflicts.
7. Assert the exact order of observable results: verified release, executable replacement, safe skill sync, clean-or-drift status summary.
8. Assert that the clean changed file updates, the clean stale file disappears, and both user files remain byte-identical.
9. Run `kk uninstall` and assert that clean owned files disappear while the modified and unknown files remain.
10. Add release CI jobs for unit tests, race tests where supported, vet, E2E, snapshot packaging, archive inspection, and checksum verification.
11. Replace the README `cp -R` instructions with the CLI bootstrap, install, update, status, uninstall, target selection, restart, and recovery commands.
12. Document that KevinKit is open source and has no license enforcement, login, or paid registry.

## Test Scenario Matrix

| Scenario | Trigger | Observable result |
|---|---|---|
| Fresh install | `kk install --yes` | Both target roots contain the full embedded catalog. |
| No-op install | Repeat install | No file or manifest content changes. |
| User edit | Edit one managed file, then update | The edited bytes remain and status reports drift. |
| Stale clean file | Remove a file from payload B | Update removes the old clean file. |
| Unknown file | Add a file inside `kk-*` | Update and uninstall preserve it. |
| Corrupt release | Serve a bad archive | Update fails before execution or replacement. |
| Full removal | `kk uninstall --yes` | Only clean owned files disappear. |

## Success Criteria

- [x] One E2E test proves the full user flow through real process calls.
- [x] The release job cannot publish an artifact that did not pass that flow.
- [x] README commands match the released binary and work from a clean temporary home.
- [x] No test reads or writes the developer's real runtime skill directories.

## Risk Assessment

The main risk is a shallow test that calls internal functions but misses wiring errors.
The gate is a process-level test that starts from the public `kk` command and verifies the real filesystem result.

## Security Considerations

Use only local fixture servers and temporary credentials-free directories in tests.
Scan release logs and documentation so they do not expose absolute developer paths or environment values.
