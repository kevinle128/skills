# Plan the Work

A useful plan tells an implementer what to change, why the sequence matters, and how to prove the complete user flow works.

`kk:plan` is the main planning skill.

It uses `kk:brainstorm`, `kk:scout`, research skills, `kk:validate-plan`, and adversarial review as the task requires.

## Start with the Delivery Contract

Use [`kk:brainstorm`](../../skills/kk-brainstorm/SKILL.md) when the request does not yet define a stable outcome or when multiple implementation directions are credible.

```text
/kk:brainstorm We need users to change runtime model settings without restarting their session.
```

The result must contain the outcome, constraints, non-goals, and observable acceptance criteria.

The brainstorm is not a substitute for codebase evidence.

If the request is a bug, the workflow scouts and diagnoses before choosing a repair.

## Find the Real Runtime Path

Use [`kk:scout`](../../skills/kk-scout/SKILL.md) to locate the entry surface, current owner, internal calls, tests, and public contracts.

```text
/kk:scout Find every entry point that changes a user setting and trace it to the observable result.
```

Use [`kk:docs-seeker`](../../skills/kk-docs-seeker/SKILL.md) when the plan depends on current third-party APIs or framework behavior.

Use [`kk:repomix`](../../skills/kk-repomix/SKILL.md) when a repository or external codebase must be packed into a focused analysis artifact.

## Select a Planning Mode

```text
/kk:plan <task> [mode] [composable flags]
```

| Mode | Use it when | Main behavior |
| --- | --- | --- |
| `--auto` | You want evidence-based mode selection. | Selects the planning depth from task and codebase signals. |
| `--fast` | The task is small, clear, and familiar. | Skips dedicated research but still creates a plan. |
| `--hard` | The task is complex or unfamiliar. | Uses research, scouting, red-team review, and optional validation. |
| `--deep` | A major refactor touches at least five areas. | Adds per-phase inventories, test matrices, interface checklists, and dependency maps. |
| `--parallel` | Independent work can run concurrently. | Adds exclusive file ownership and an execution dependency graph. |
| `--two` | Two credible architectures should be compared. | Produces two approaches and waits for a selection before finalizing. |

The main composable flags are:

- `--tdd` adds tests-first steps and regression gates to each phase.
- `--no-tasks` keeps the plan files without hydrating task management.
- `--html` creates an interactive self-contained plan artifact.
- `--github` projects the reviewed plan into a GitHub issue.
- `--wiki` publishes the reviewed artifact to AgentWiki when available.
- `--advice` adds advisory supervision at major decisions.
- `--yagni` explicitly permits removal of unnecessary scope.
- `--skip-journal` skips the optional journal step.
- `--global` places a cross-project plan in the configured global plan root.

The complete flag contract lives in [`kk:plan`](../../skills/kk-plan/SKILL.md).

## Write for the Implementer

Each phase should identify:

- the behavior it delivers;
- the files and interfaces it changes;
- the implementation steps at code level;
- dependencies and ownership boundaries;
- success criteria and runnable validation;
- risks, failure signals, and rollback behavior.

The plan must connect component-level work to an observable end-to-end result.

For an API feature, include a test that calls the real API boundary.

For a website, include a browser flow that uses the feature as a user would.

For a cron job or listener, use the real scheduler or event boundary and assert the externally observable effect.

Mock third-party systems at their boundary.

Call internal services through their real integration path and prepare realistic test data.

Every planning mode runs the Runtime Flow Proof Gate before the plan is ready.

If no user or system trigger exists for a required behavior, stop and ask the user to select a runtime surface before claiming the plan is implementable.

The plan must record one proof row per feature and cannot proceed while any row is `FAILED` or `NEEDS_DECISION`.

## Validate the Plan

```text
/kk:validate-plan /absolute/path/to/plan-directory
```

Validation rechecks the Runtime Flow Proof Matrix against the codebase, interviews the user about material decisions, propagates the answers into affected phases, and performs a whole-plan consistency sweep.

The complete validation contract lives in [`kk:validate-plan`](../../skills/kk-validate-plan/SKILL.md).

Failed verification claims must be revised before implementation.

## Red-Team the Plan

```text
/kk:plan red-team /absolute/path/to/plan-directory
```

Red-team review uses independent hostile lenses for security, assumptions, failure modes, scope, and complexity.

Findings require codebase evidence, and the user decides which accepted findings are applied.

Run red-team review before final validation because accepted findings can change the plan.

## Plan State

`plan.md` and `phase-NN-*.md` are canonical.

AgentKit's `ak plan` database is a rebuildable index over those files.

GitHub issues and AgentWiki pages are optional projections and never replace the local plan files.

After approval, hand the absolute plan path to `kk:implement`.

```text
/kk:implement /absolute/path/to/plan-directory/plan.md
```

Next: [Implement the plan](./03-implement-the-plan.md).
