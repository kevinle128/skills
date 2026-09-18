# Kevin Kit

Kevin Kit is my personal collection of reusable AI-agent skills for vibe coding.

It stores the workflows I use to research codebases, design solutions, create implementation plans, validate decisions, and move from an idea to working software.

The current core skills are `kevinle128-skills:plan`, `kevinle128-skills:validate-plan`, and `kevinle128-skills:implement`.

## Guide

Start with [The Kevin Kit Guide](./docs/guide/README.md) for workflow selection, practical examples, and the complete skill catalog.

## `kevinle128-skills:plan`

`kevinle128-skills:plan` creates implementation-ready technical plans from a task description.

It can inspect an existing codebase, research unfamiliar areas, compare solution approaches, organize work into phases, define test coverage, review risks, and prepare a handoff to implementation.

The generated Markdown files are the source of truth.

AgentKit maintains a rebuildable plan index around those files.

### Requirements

- A compatible AI-agent runtime that can load skills from this plugin.
- The AgentKit `ak` CLI for plan scaffolding, indexing, and phase status operations.
- `gh` authentication when using `--github`.
- AgentWiki CLI authentication or AgentWiki MCP access when using `--wiki`.

Before it changes plan state, the skill checks the live `ak plan` help output instead of assuming CLI syntax.

## Basic Usage

Create a plan from a task description:

```text
/kevinle128-skills:plan Add session-based authentication to the API
```

When no mode flag is supplied, `kevinle128-skills:plan` selects a mode from the task scope and codebase evidence.

You can select a mode explicitly:

```text
/kevinle128-skills:plan Add session-based authentication to the API --hard
```

Flags can be combined:

```text
/kevinle128-skills:plan Refactor the payment workflow --deep --tdd --github
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
/kevinle128-skills:plan Replace the cache implementation --hard --tdd
```

### `--no-tasks`

Skips task hydration after the plan files are written.

Use this when you want the plan documents but do not want phases mirrored into the available task-management system.

```text
/kevinle128-skills:plan Explore a new event-processing architecture --two --no-tasks
```

### `--html`

Creates a self-contained `plan.html` as the primary user-facing artifact.

The HTML plan includes phase summaries, full phase details, an implementation workflow diagram, responsive interaction, and UI mockups when the planned work affects a user interface.

The file contains inline CSS and JavaScript and does not require a build step or network assets.

```text
/kevinle128-skills:plan Redesign the project dashboard --deep --html
```

When `--html` is combined with `--github`, a concise `plan.md` index is also kept so the GitHub issue has a stable Markdown link.

### `--github`

Creates or updates a GitHub issue after the plan review gates complete.

The issue contains the branch, plan summary, repository-relative plan links, open questions, acceptance criteria, and the `ready to review` label.

The local plan files remain canonical.

If the repository has no GitHub remote or `gh` is not authenticated, publishing is skipped without invalidating the plan.

```text
/kevinle128-skills:plan Add organization-level permissions --hard --github
```

### `--wiki`

Publishes the final reviewed plan to AgentWiki when AgentWiki CLI or MCP access is available.

The default behavior uses a private or workspace share.

Public publication happens only when the user explicitly requests it.

If AgentWiki is unavailable or authentication fails, publishing is skipped without blocking plan creation.

```text
/kevinle128-skills:plan Document the new deployment architecture --hard --wiki
```

### `--advice`

Runs the planning workflow with `kongming` as an advisory supervisor.

The supervisor reviews major planning checkpoints, difficult blockers, high-stakes design decisions, and the final implementation when the workflow reaches a pull request.

It advises the primary agent but does not replace approval, validation, or review gates.

```text
/kevinle128-skills:plan Migrate the authorization model --deep --advice
```

### `--yagni`

Explicitly allows the skill to challenge and remove scope that is not required for the stated outcome.

Without this flag, the skill preserves the full requested scope.

The flag is forwarded to planning subagents and downstream workflows so the scope decision remains active.

```text
/kevinle128-skills:plan Simplify the notification service --hard --yagni
```

### `--skip-journal`

Skips the optional journal step.

This is useful for temporary planning sessions or when the plan does not need a chronological project record.

```text
/kevinle128-skills:plan archive ./plans/260918-authentication --skip-journal
```

### `--global`

Uses the configured global plans root instead of the current project's plan directory.

Use it only when the plan is intentionally shared across projects.

Global scope is also allowed when no project context exists.

```text
/kevinle128-skills:plan Standardize release checks across repositories --global --hard
```

## `kevinle128-skills:validate-plan`

Validates an existing plan against the real codebase and proves every feature through its production trigger and end-to-end test.

If a production surface is missing, it opens a HITL question with codebase-grounded options and a recommendation.
It records confirmed decisions, propagates them to affected phases, and runs a whole-plan consistency sweep before recommending implementation.

```text
/kevinle128-skills:validate-plan /absolute/path/to/plans/260918-authentication
```

## `kevinle128-skills:plan` Subcommands

### `red-team`

Runs an adversarial review of an existing plan.

Independent reviewers examine security, assumptions, failure modes, scope, and complexity.

Only findings with codebase evidence are considered, and the user decides which accepted findings are applied.

```text
/kevinle128-skills:plan red-team /absolute/path/to/plans/260918-authentication
```

### `archive`

Archives one or more plans in the AgentKit index after user confirmation.

Archiving changes index visibility and does not delete or move the canonical Markdown files.

```text
/kevinle128-skills:plan archive /absolute/path/to/plans/260918-authentication
```

Use `--skip-journal` with this subcommand when no journal entry is required.

## Common Recipes

Create a quick plan for a well-understood change:

```text
/kevinle128-skills:plan Add a health-check endpoint --fast
```

Plan a production feature with tests-first implementation phases:

```text
/kevinle128-skills:plan Add passwordless login --hard --tdd
```

Plan a large refactor and produce an interactive review artifact:

```text
/kevinle128-skills:plan Replace the job scheduler --deep --tdd --html
```

Prepare independent work for parallel implementation:

```text
/kevinle128-skills:plan Build the API, admin UI, and audit pipeline --parallel --tdd
```

Compare two architectures before committing to one:

```text
/kevinle128-skills:plan Introduce multi-region data replication --two --advice
```

Publish a reviewed plan to GitHub and AgentWiki:

```text
/kevinle128-skills:plan Add tenant isolation --deep --github --wiki
```

## Plan Output

Project plans are stored in the project's configured plan directory, which defaults to `plans/`.

Each standard plan contains:

- `plan.md` for the overview, dependencies, phases, decisions, and acceptance criteria.
- `phase-NN-*.md` files for implementation details, related files, validation, risks, and rollback information.
- `plan.html` when `--html` is active.
- `wiki-publish.md` when AgentWiki needs a combined Markdown publication artifact.

The generated files remain editable and are the source of truth even when AgentKit task state, GitHub issues, or AgentWiki pages are also created.

## `kevinle128-skills:implement`

`kevinle128-skills:implement` executes an accepted plan or implements a clearly defined task through a structured delivery workflow.

It supports interactive approval gates, fast execution, parallel agents, autonomous execution, tests-first development, testing, code review, plan synchronization, and final Git operations.

Use it directly with a task:

```text
/kevinle128-skills:implement "Add user authentication to the app" --interactive
```

Use it with a plan created by `kevinle128-skills:plan`:

```text
/kevinle128-skills:implement /absolute/path/to/plan-directory/plan.md
```

The implementation modes are `--interactive`, `--fast`, `--parallel`, `--auto`, and `--no-test`.

The composable flags are `--tdd`, `--advice`, `--yagni`, and `--skip-journal`.

## Implementation Handoff

After reviewing and approving a plan, start implementation with the absolute plan path shown by `kevinle128-skills:plan`:

```text
/kevinle128-skills:implement /absolute/path/to/plan-directory/plan.md
```

For a parallel plan, use:

```text
/kevinle128-skills:implement --parallel /absolute/path/to/plan-directory/plan.md
```

When the plan was created with `--tdd`, keep `--tdd` in the implementation handoff.

The skill recommends clearing the planning context before implementation so the implementation session starts with the plan as its focused source of truth.

## Repository Layout

- `.codex-plugin/plugin.json` contains the plugin metadata.
- `skills/plan/SKILL.md` contains the main planning workflow.
- `skills/plan/references/` contains the supporting workflow contracts.
- `skills/implement/SKILL.md` contains the implementation workflow.
- `skills/implement/references/` contains its routing, review, and execution contracts.

More Kevin Kit skills can be added under `skills/` as the collection grows.
