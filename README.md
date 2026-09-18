# KevinKit

KevinKit is my personal collection of reusable AI-agent skills for vibe coding.

It stores the workflows I use to research codebases, design solutions, create implementation plans, validate decisions, and move from an idea to working software.

KevinKit is open source and has no license enforcement, account, login, analytics, paid registry, or entitlement check.

The current core skills are `kk:plan`, `kk:validate-plan`, and `kk:implement`.

## Set Up the `kk` CLI

KevinKit ships as one `kk` binary with all skills embedded.

Release archives support macOS and Linux on `amd64` and `arm64`.

Install the binary in a user-writable directory so `kk update` can replace it without `sudo`.

### 1. Select the Release Archive

Download the latest release archive and `checksums.txt` from [GitHub Releases](https://github.com/kevinle128/skills/releases).

Use this table to select the archive suffix for your computer.

| Computer | Archive suffix |
| --- | --- |
| Apple silicon Mac | `darwin_arm64.tar.gz` |
| Intel Mac | `darwin_amd64.tar.gz` |
| ARM64 Linux | `linux_arm64.tar.gz` |
| AMD64 or x86-64 Linux | `linux_amd64.tar.gz` |

The complete archive name is `kk_<version>_<os>_<arch>.tar.gz`.

You can also download a specific version from a terminal.

Set `VERSION`, `OS`, and `ARCH` to match the release and your computer.

```bash
VERSION="<version-without-v>"
OS="darwin"
ARCH="arm64"
ARCHIVE="kk_${VERSION}_${OS}_${ARCH}.tar.gz"

curl -fLO "https://github.com/kevinle128/skills/releases/download/v${VERSION}/${ARCHIVE}"
curl -fLO "https://github.com/kevinle128/skills/releases/download/v${VERSION}/checksums.txt"
```

### 2. Verify the Download

Verify the archive before you extract or run it.

On macOS, run:

```bash
set -o pipefail
grep "  ${ARCHIVE}$" checksums.txt | shasum -a 256 --check
```

On Linux, run:

```bash
set -o pipefail
grep "  ${ARCHIVE}$" checksums.txt | sha256sum --check
```

The command must report the archive as `OK`.

Do not install the archive when verification fails or no matching checksum exists.

### 3. Install the Binary

Extract the archive and install `kk` in `~/.local/bin`.

```bash
tar -xzf "$ARCHIVE"
install -d "$HOME/.local/bin"
install -m 0755 kk "$HOME/.local/bin/kk"
```

Add this line to `~/.zshrc`, `~/.bashrc`, or the startup file for your shell when `~/.local/bin` is not already on `PATH`.

```bash
export PATH="$HOME/.local/bin:$PATH"
```

Start a new terminal, then verify the installation.

```bash
kk version
kk --help
```

### Install from Source

Go 1.24 or later can install a development build.

```bash
go install github.com/kevinle128/skills/cmd/kk@latest
export PATH="$(go env GOPATH)/bin:$PATH"
kk version
```

A source-installed binary reports the version as `dev` until `kk update --yes` replaces it with a published release.

## Install the Skills

Preview the initial installation without writing files.

```bash
kk install --dry-run
```

Install all embedded skills for Codex, Claude Code, and Devin.

```bash
kk install
```

The default target is `all`.

| Target | Runtime | Skill root |
| --- | --- | --- |
| `agents` | Codex and Devin | `~/.agents/skills` |
| `claude-code` | Claude Code | `~/.claude/skills` |
| `all` | Codex, Devin, and Claude Code | Both roots |

Install only one target when required.

```bash
kk install --target agents
kk install --target claude-code
```

Start a new agent session after installation so the runtime reloads its skill catalog.

## Use the CLI

### Check the Installation

```bash
kk status
```

`status: clean` means the installed files match the KevinKit ownership manifest.

`status: drift` lists missing, modified, or stale managed files and returns exit code `1`.

Use `--target` to inspect one runtime root.

```bash
kk status --target agents
```

### Synchronize the Embedded Skills

Run `install` again after moving the binary or when you need to restore missing clean files.

```bash
kk install
```

This command synchronizes the skills embedded in the current binary.

It does not download a newer release.

### Update KevinKit

Check for a newer stable release without changing local state.

```bash
kk update --check
```

Preview the release that would be installed.

```bash
kk update --dry-run
```

Verify the published checksum, replace the current binary, and synchronize the new embedded skills.

```bash
kk update --yes
```

Limit the skill synchronization to one target when required.

```bash
kk update --target agents --yes
```

The updater keeps a backup of the previous binary under `~/.kevinkit/backups/`.

### Protect Local Changes

KevinKit records the SHA-256 value of each file that it installs.

Normal install, update, and uninstall operations preserve user-modified and unknown files.

Use `status` to review these files before you decide what to do.

```bash
kk status
```

Use force only when you want to replace conflicting files with the embedded copies.

```bash
kk install --force --yes
```

You can also allow an update to replace conflicts after it downloads the verified release.

```bash
kk update --force --yes
```

KevinKit creates recovery copies under `~/.kevinkit/backups/` before destructive writes.

### Uninstall the Skills

Preview removal first.

```bash
kk uninstall --dry-run
```

Remove clean managed files from all targets.

```bash
kk uninstall --yes
```

Remove files from one target only.

```bash
kk uninstall --target claude-code --yes
```

Uninstall preserves modified and unknown files.

Remove the `kk` binary separately when you no longer need the CLI.

```bash
rm "$HOME/.local/bin/kk"
```

## CLI Reference

| Command | Result |
| --- | --- |
| `kk install` | Install or synchronize the skills embedded in the current binary. |
| `kk update --yes` | Install a checksum-verified release and synchronize its embedded skills. |
| `kk update --check` | Report whether a newer stable release is available without writing. |
| `kk status` | Report installed targets and managed-file drift. |
| `kk uninstall --yes` | Remove clean managed files and preserve user files. |
| `kk version` | Print CLI build metadata and the embedded kit version. |
| `kk help <command>` | Print usage for one command. |

Run `kk help <command>` for the exact flags accepted by a command.

## Troubleshooting

### `kk: command not found`

Confirm that the directory containing `kk` is on `PATH`.

```bash
command -v kk
printf '%s\n' "$PATH"
```

Add `~/.local/bin` or `$(go env GOPATH)/bin` to your shell startup file, then start a new terminal.

### `status: drift`

Run `kk status` and review every reported path.

Keep the files unchanged when they contain work that you want to preserve.

Run `kk install --force --yes` only when you want to discard those changes and restore the embedded copies.

### Update Cannot Replace the Binary

The directory that contains the running `kk` binary must be writable by the current user.

Install `kk` under `~/.local/bin` instead of running `kk update` with `sudo`.

Running the lifecycle commands with `sudo` would resolve a different home directory and target the wrong skill roots.

### Inspect Local State

KevinKit keeps its ownership manifest, lifecycle lock, and backups under `~/.kevinkit`.

Do not edit `install-manifest.json` manually.

## Guide

Start with [The KevinKit Guide](./docs/guide/README.md) for workflow selection, practical examples, and the complete skill catalog.

## `kk:plan`

`kk:plan` creates implementation-ready technical plans from a task description.

It can inspect an existing codebase, research unfamiliar areas, compare solution approaches, organize work into phases, define test coverage, review risks, and prepare a handoff to implementation.

The generated Markdown files are the source of truth.

AgentKit maintains a rebuildable plan index around those files.

### Requirements

- A compatible AI-agent runtime that can load standalone skills.
- The AgentKit `ak` CLI for plan scaffolding, indexing, and phase status operations.
- `gh` authentication when using `--github`.
- AgentWiki CLI authentication or AgentWiki MCP access when using `--wiki`.

Before it changes plan state, the skill checks the live `ak plan` help output instead of assuming CLI syntax.

## Basic Usage

Create a plan from a task description:

```text
/kk:plan Add session-based authentication to the API
```

When no mode flag is supplied, `kk:plan` selects a mode from the task scope and codebase evidence.

You can select a mode explicitly:

```text
/kk:plan Add session-based authentication to the API --hard
```

Flags can be combined:

```text
/kk:plan Refactor the payment workflow --deep --tdd --github
```

When invoked without a task or subcommand, the skill asks whether to create, archive, or red-team a plan.

## Planning Modes

Only one planning mode should normally be selected for one invocation.

| Flag | Use it when | Behavior |
| --- | --- | --- |
| `--auto` | You want the skill to choose the planning depth. | Selects `fast`, `hard`, `deep`, `parallel`, or `two` from task complexity and codebase signals. |
| `--fast` | The task is small, clear, and already understood. | Skips dedicated research and review gates, then creates the plan and hydrates phase-level tasks. |
| `--hard` | The task is complex or touches an unfamiliar area. | Uses up to two researchers, inspects the codebase, creates the plan, runs red-team review, and offers validation. |
| `--deep` | A major refactor touches at least five areas or has architectural debt. | Adds per-phase scouting, file inventories, test matrices, interface checklists, dependency maps, red-team review, and validation. |
| `--parallel` | Three or more independent features, layers, or modules can be implemented concurrently. | Defines exclusive file ownership, a dependency graph, parallel groups, conflict prevention, review, and validation. |
| `--two` | The problem has multiple credible implementation approaches. | Produces two approaches with trade-offs and a recommendation, waits for a selection, then reviews and validates the selected approach. |

### Automatic Mode Selection

The default selection uses these signals:

| Signal | Selected mode |
| --- | --- |
| Simple task, clear scope, no material unknowns | `fast` |
| Complex task, unfamiliar domain, or new technology | `hard` |
| Major refactor across at least five areas | `deep` |
| At least three independent work areas | `parallel` |
| Ambiguous design with multiple valid approaches | `two` |

If the evidence is unclear, the skill asks the user before selecting a mode.

`--fast` reduces planning depth only.

It does not reduce the requested feature scope.

## Composable Flags

These flags can be combined with a planning mode.

### `--tdd`

Adds a tests-first structure to every implementation phase.

Each phase identifies tests for existing behavior, the protected implementation work, tests for new behavior, and a regression gate containing the required test and compile or type-check commands.

```text
/kk:plan Replace the cache implementation --hard --tdd
```

### `--no-tasks`

Skips task hydration after the plan files are written.

Use this when you want the plan documents but do not want phases mirrored into the available task-management system.

```text
/kk:plan Explore a new event-processing architecture --two --no-tasks
```

### `--html`

Creates a self-contained `plan.html` as the primary user-facing artifact.

The HTML plan includes phase summaries, full phase details, an implementation workflow diagram, responsive interaction, and UI mockups when the planned work affects a user interface.

The file contains inline CSS and JavaScript and does not require a build step or network assets.

```text
/kk:plan Redesign the project dashboard --deep --html
```

When `--html` is combined with `--github`, a concise `plan.md` index is also kept so the GitHub issue has a stable Markdown link.

### `--github`

Creates or updates a GitHub issue after the plan review gates complete.

The issue contains the branch, plan summary, repository-relative plan links, open questions, acceptance criteria, and the `ready to review` label.

The local plan files remain canonical.

If the repository has no GitHub remote or `gh` is not authenticated, publishing is skipped without invalidating the plan.

```text
/kk:plan Add organization-level permissions --hard --github
```

### `--wiki`

Publishes the final reviewed plan to AgentWiki when AgentWiki CLI or MCP access is available.

The default behavior uses a private or workspace share.

Public publication happens only when the user explicitly requests it.

If AgentWiki is unavailable or authentication fails, publishing is skipped without blocking plan creation.

```text
/kk:plan Document the new deployment architecture --hard --wiki
```

### `--advice`

Runs the planning workflow with `kongming` as an advisory supervisor.

The supervisor reviews major planning checkpoints, difficult blockers, high-stakes design decisions, and the final implementation when the workflow reaches a pull request.

It advises the primary agent but does not replace approval, validation, or review gates.

```text
/kk:plan Migrate the authorization model --deep --advice
```

### `--yagni`

Explicitly allows the skill to challenge and remove scope that is not required for the stated outcome.

Without this flag, the skill preserves the full requested scope.

The flag is forwarded to planning subagents and downstream workflows so the scope decision remains active.

```text
/kk:plan Simplify the notification service --hard --yagni
```

### `--skip-journal`

Skips the optional journal step.

This is useful for temporary planning sessions or when the plan does not need a chronological project record.

```text
/kk:plan archive ./plans/260918-authentication --skip-journal
```

### `--global`

Uses the configured global plans root instead of the current project's plan directory.

Use it only when the plan is intentionally shared across projects.

Global scope is also allowed when no project context exists.

```text
/kk:plan Standardize release checks across repositories --global --hard
```

## `kk:validate-plan`

Validates an existing plan against the real codebase and rechecks every feature through its runtime trigger and end-to-end test.

`kk:plan` runs the same Runtime Flow Proof Gate before it marks a new plan ready, including in `--fast` mode.

If a runtime surface is missing, either skill opens a HITL question with codebase-grounded options and a recommendation.
It records confirmed decisions, propagates them to affected phases, and runs a whole-plan consistency sweep before recommending implementation.

```text
/kk:validate-plan /absolute/path/to/plans/260918-authentication
```

## `kk:plan` Subcommands

### `red-team`

Runs an adversarial review of an existing plan.

Independent reviewers examine security, assumptions, failure modes, scope, and complexity.

Only findings with codebase evidence are considered, and the user decides which accepted findings are applied.

```text
/kk:plan red-team /absolute/path/to/plans/260918-authentication
```

### `archive`

Archives one or more plans in the AgentKit index after user confirmation.

Archiving changes index visibility and does not delete or move the canonical Markdown files.

```text
/kk:plan archive /absolute/path/to/plans/260918-authentication
```

Use `--skip-journal` with this subcommand when no journal entry is required.

## Common Recipes

Create a quick plan for a well-understood change:

```text
/kk:plan Add a health-check endpoint --fast
```

Plan a production feature with tests-first implementation phases:

```text
/kk:plan Add passwordless login --hard --tdd
```

Plan a large refactor and produce an interactive review artifact:

```text
/kk:plan Replace the job scheduler --deep --tdd --html
```

Prepare independent work for parallel implementation:

```text
/kk:plan Build the API, admin UI, and audit pipeline --parallel --tdd
```

Compare two architectures before committing to one:

```text
/kk:plan Introduce multi-region data replication --two --advice
```

Publish a reviewed plan to GitHub and AgentWiki:

```text
/kk:plan Add tenant isolation --deep --github --wiki
```

## Plan Output

Project plans are stored in the project's configured plan directory, which defaults to `plans/`.

Each standard plan contains:

- `plan.md` for the overview, dependencies, phases, decisions, and acceptance criteria.
- `phase-NN-*.md` files for implementation details, related files, validation, risks, and rollback information.
- `plan.html` when `--html` is active.
- `wiki-publish.md` when AgentWiki needs a combined Markdown publication artifact.

The generated files remain editable and are the source of truth even when AgentKit task state, GitHub issues, or AgentWiki pages are also created.

## `kk:implement`

`kk:implement` executes an accepted plan or implements a clearly defined task through a structured delivery workflow.

It supports interactive approval gates, fast execution, parallel agents, autonomous execution, tests-first development, testing, code review, plan synchronization, and final Git operations.

Use it directly with a task:

```text
/kk:implement "Add user authentication to the app" --interactive
```

Use it with a plan created by `kk:plan`:

```text
/kk:implement /absolute/path/to/plan-directory/plan.md
```

The implementation modes are `--interactive`, `--fast`, `--parallel`, `--auto`, and `--no-test`.

The composable flags are `--tdd`, `--advice`, `--yagni`, and `--skip-journal`.

## Implementation Handoff

After reviewing and approving a plan, start implementation with the absolute plan path shown by `kk:plan`:

```text
/kk:implement /absolute/path/to/plan-directory/plan.md
```

For a parallel plan, use:

```text
/kk:implement --parallel /absolute/path/to/plan-directory/plan.md
```

When the plan was created with `--tdd`, keep `--tdd` in the implementation handoff.

The skill recommends clearing the planning context before implementation so the implementation session starts with the plan as its focused source of truth.

## Repository Layout

- `skills/kk-plan/SKILL.md` contains the main planning workflow.
- `skills/kk-plan/references/` contains the supporting workflow contracts.
- `skills/kk-implement/SKILL.md` contains the implementation workflow.
- `skills/kk-implement/references/` contains its routing, review, and execution contracts.

More KevinKit skills can be added under `skills/` as the collection grows.
