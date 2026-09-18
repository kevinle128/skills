# Xia Analysis: AgentKit Lifecycle to KevinKit CLI

## Source Manifest

- Source feature: installed AgentKit CLI lifecycle behavior.
- Executable: `/Users/dale/.local/bin/ak`.
- Version: `2.13.0-beta.38`.
- Build revision: `6c7de83f4961164d87a4a4d811dca914dff3aea2`.
- Build form: static Go binary with build tag `embeddedkits`.
- Source repository reference: `github.com/bestagentkits/agentkit/apps/cli`.
- Source code access: unavailable in the local workspace and unavailable through the current GitHub credentials.
- Behavioral evidence: live `ak kit install`, `ak kit refresh`, `ak kit uninstall`, `ak audit`, `ak update`, and `ak self-update` help contracts.
- State evidence: AgentKit install manifests contain sorted file paths and per-file SHA-256 values.

## Local Manifest

- Repository: `git@github.com:kevinle128/skills.git`.
- Branch: `main`.
- Inspected commit: `093e1b9e5974976acac9399ee6b98b15ad2bc1cb`.
- Payload: 34 `skills/kk-*` directories, 407 files, and about 4.5 MB.
- Installed copies: exact directory copies under `~/.agents/skills` and `~/.claude/skills`.
- Current installer: manual `cp -R` loop in `README.md`.
- Existing CLI or build manifest: none.

## Source Anatomy

| Layer | AgentKit behavior | KevinKit adoption |
|---|---|---|
| Payload | Kits embedded in a Go binary and available from a registry. | Embed the public `skills/kk-*` tree in one Go binary. |
| Targets | Adapter-specific project and global destinations. | Two global destinations with one shared Codex and Devin target. |
| Ownership | Install manifest with per-file SHA-256 values. | Adopt the same invariant with a smaller manifest. |
| Update | Self-update binary first, then refresh kit output. | Use one `kk update` user command with the same order. |
| Safety | Snapshot before mutation and preserve modified or unknown files. | Back up affected files and preserve conflicts by default. |
| Audit | Compare installed bytes against install-time fingerprints. | Provide the smaller read-only `kk status` command. |
| Distribution | Signed private release and paid registry channels. | Public GitHub Releases with checksums and no auth or license layer. |

## Dependency Matrix

| Component | Local state | Decision |
|---|---|---|
| Canonical skill payload | EXISTS | Reuse `skills/kk-*`. |
| Runtime target paths | EXISTS | Move README knowledge into tested target mapping code. |
| Go module and CLI entry point | NEW | Add a small standard-library command. |
| Embedded payload catalog | NEW | Add root `go:embed` package and parity tests. |
| Manifest and sync engine | NEW | Add one shared planner and executor for all lifecycle commands. |
| Backup and mutation lock | NEW | Add config-local recovery data and serialization. |
| GitHub release updater | NEW | Add HTTPS discovery, checksum verification, and replacement. |
| Release pipeline | NEW | Add GoReleaser and GitHub Actions. |
| License, auth, registry, TUI, database | CONFLICT | Do not port because KevinKit is public and the use cases do not exist. |

## Challenge Results

| # | Question | Source answer | Local answer | Risk if wrong |
|---|---|---|---|---|
| 1 | Does a CLI need to exist? | AgentKit needs lifecycle management for many adapters and kits. | Yes, because manual copies have no ownership, drift, or update model. | Users keep stale or mixed skill versions. |
| 2 | What is the source of truth? | Embedded kits plus remote registry. | The selected GitHub Release binary and its embedded payload. | A local checkout becomes an accidental runtime dependency. |
| 3 | How are user edits protected? | Compare current bytes with install-time SHA-256. | Use the same invariant and preserve conflicts by default. | Update or uninstall deletes user work. |
| 4 | What does update mean? | Update binary, then refresh kit output. | One command performs both steps in that order. | CLI and installed skills report different versions. |
| 5 | How do shared targets behave? | Adapter ownership resolves each output. | Codex and Devin map once to `agents`. | The same directory is changed twice with inconsistent state. |
| 6 | Which AgentKit systems are necessary? | Registry, auth, licenses, TUI, and databases support its product. | None of them are necessary for a public skill bundle. | The port becomes larger and harder to secure without user value. |

## Decision Matrix

| Decision | AgentKit way | KevinKit way | Recommendation |
|---|---|---|---|
| Runtime | Go binary with external libraries. | Go binary with standard library runtime code. | Keep the binary and remove unnecessary dependencies. |
| Content | Embedded kits plus authenticated registry. | Embedded public skills in GitHub Release binaries. | Use GitHub Releases. |
| Integrity | Signed update manifest and SHA-256 fingerprints. | HTTPS GitHub trust root, release checksums, and install fingerprints. | Fail closed on every checksum mismatch. |
| Storage | Adapter manifests, snapshots, databases, and caches. | `~/.kevinkit` with one JSON manifest, lock, temporary area, and backups. | Keep only lifecycle state. |
| User interface | TUI, GUI, pretty, plain, and JSON modes. | Direct CLI commands with clear stdout and stderr. | Do not port TUI or GUI. |
| Scope | Project and global installation. | Global runtime skill roots only. | Add project scope only after a real use case appears. |

## Risk Score

Risk is medium because three assumptions can cause data loss, supply-chain execution, or a broken executable if implemented incorrectly.
The plan addresses them with manifest ownership tests, checksum verification before execution, executable backup, and a process-level E2E release gate.

## Approved Outcome

The user selected GitHub Releases instead of local checkout or Git pull as the update source.
KevinKit remains open source and contains no license enforcement, login, entitlement, or paid registry behavior.
