# Choose a Workflow

Kevin Kit has several entry skills because a new feature, a production defect, and a completed branch need different safety gates.

Start from the user outcome, not from a list of tools.

## Choose the Front Door

| Your situation | Start with | What it owns |
| --- | --- | --- |
| You are starting a new application or service. | [`kevinle128-skills:bootstrap`](../../skills/bootstrap/SKILL.md) | Research, stack selection, planning, implementation, testing, and initial delivery. |
| The desired behavior is still unclear or has meaningful design choices. | [`kevinle128-skills:brainstorm`](../../skills/brainstorm/SKILL.md) | Outcome, constraints, non-goals, acceptance criteria, and approach selection. |
| You need an implementation plan or architecture roadmap. | [`kevinle128-skills:plan`](../../skills/plan/SKILL.md) | Codebase research, phases, dependencies, validation, and implementation handoff. |
| You need to verify an existing plan before implementation. | [`kevinle128-skills:validate-plan`](../../skills/validate-plan/SKILL.md) | Codebase verification, material decisions, answer propagation, and whole-plan consistency. |
| You already have an accepted plan or clear implementation contract. | [`kevinle128-skills:implement`](../../skills/implement/SKILL.md) | Code changes, tests, review, progress tracking, and finalization. |
| You have a reproducible bug, failed test, or CI failure. | [`kevinle128-skills:fix`](../../skills/fix/SKILL.md) | Reproduction, diagnosis, root-cause repair, regression checks, and review. |
| You need diagnosis but do not want code changes yet. | [`kevinle128-skills:debug`](../../skills/debug/SKILL.md) | Evidence gathering, hypothesis testing, and root-cause reporting. |
| The implementation is finished and the branch must become a PR. | [`kevinle128-skills:git`](../../skills/git/SKILL.md) | Focused commits, push, pull request creation, and merge operations. |

## Give the Workflow a Contract

A strong request answers four questions.

- **Outcome:** What user-visible or operational result must exist?
- **Constraints:** What compatibility, safety, technology, or ownership boundaries apply?
- **Non-goals:** What nearby work must remain outside this change?
- **Acceptance criteria:** What observable evidence proves the work is complete?

You can write these in plain language.

```text
/kevinle128-skills:plan Add an API that lets a user change the active model, thinking mode, and fast mode.
Keep existing session APIs compatible.
Do not redesign provider configuration.
Done means an end-to-end API test performs the change and proves the live worker uses the new values.
```

## Common Routes

### New Feature

```text
kevinle128-skills:brainstorm -> kevinle128-skills:scout -> kevinle128-skills:plan -> kevinle128-skills:validate-plan -> kevinle128-skills:implement -> kevinle128-skills:test -> kevinle128-skills:code-review -> kevinle128-skills:git
```

Start directly at `kevinle128-skills:plan` when the outcome and acceptance criteria are already clear.

Start directly at `kevinle128-skills:implement` when an accepted plan already exists.

### Bug or Regression

```text
intent frame -> kevinle128-skills:scout -> kevinle128-skills:debug -> kevinle128-skills:fix -> kevinle128-skills:test -> kevinle128-skills:code-review
```

`kevinle128-skills:fix` owns this route and can call the supporting skills itself.

Use `kevinle128-skills:debug` alone when the requested outcome is diagnosis rather than repair.

### GitHub Issue to Delivery

```text
kevinle128-skills:scout -> kevinle128-skills:plan -> kevinle128-skills:plan red-team -> kevinle128-skills:validate-plan -> kevinle128-skills:implement -> kevinle128-skills:review-pr -> kevinle128-skills:git
```

The issue is a coordination surface.

The plan files remain the implementation source of truth.

### New Project

```text
kevinle128-skills:bootstrap -> kevinle128-skills:plan -> kevinle128-skills:implement -> kevinle128-skills:test -> kevinle128-skills:code-review -> kevinle128-skills:docs
```

`kevinle128-skills:bootstrap` chooses the correct planning and implementation modes from its own flags.

### Completed Branch

```text
kevinle128-skills:test -> kevinle128-skills:code-review -> kevinle128-skills:git
```

Use `kevinle128-skills:review-pr` instead of `kevinle128-skills:code-review` when the primary artifact is an existing GitHub pull request.

## When to Invoke a Supporting Skill Directly

Most supporting skills are called by an owning workflow.

Invoke one directly when the supporting operation is the complete task.

Examples include using `kevinle128-skills:scout` to locate an owner, `kevinle128-skills:test` to verify a branch, `kevinle128-skills:docs` to update documentation, or `kevinle128-skills:preview` to explain an architecture visually.

## Pitfall: Writing the Whole Workflow in the Prompt

Avoid prompts such as "run scout, then debug, then three reviewers, then ship."

That sequence may bypass an owning skill's gates or put steps in the wrong order.

State the outcome and evidence, select the front door, and override only the behavior that you intentionally want to change.

Next: [Plan the work](./02-plan-the-work.md).
