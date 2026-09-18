# Test, Review, and Ship

Compilation is not proof that a feature works.

Kevin Kit separates execution evidence, code review, pull-request review, and release operations so each gate has a clear owner.

## Run the Right Tests

Use [`kevinle128-skills:test`](../../skills/test/SKILL.md) for unit, integration, end-to-end, UI, coverage, build, and visual verification.

```text
/kevinle128-skills:test Run the API integration and end-to-end tests for session setting updates.
```

For a web application, provide the URL and the user behavior that must work.

```text
/kevinle128-skills:test ui http://localhost:3000
Sign in, change the active model, reload the page, and verify the selected value is still active.
```

Use [`kevinle128-skills:web-testing`](../../skills/web-testing/SKILL.md) when you need Playwright, Vitest, k6, cross-browser checks, accessibility, load testing, or visual regression.

```text
/kevinle128-skills:web-testing e2e http://localhost:3000/settings
```

## Build an Evidence Ladder

Use the narrowest test that proves the changed logic while developing.

Before completion, broaden the evidence to the affected contract.

1. Run the focused unit or component test.
2. Run integration tests through real internal services.
3. Run the public API, UI, command, event, or scheduler flow end to end.
4. Run shared lint, type, build, and regression suites when common contracts changed.

Mock third parties at their boundary.

Do not mock the internal path that the feature must prove.

## Review Local Changes

Use [`kevinle128-skills:code-review`](../../skills/code-review/SKILL.md) for pending changes, a commit, a pull request, or a codebase scan.

```text
/kevinle128-skills:code-review --pending
```

```text
/kevinle128-skills:code-review abc1234
```

```text
/kevinle128-skills:code-review #123
```

The reviewer prioritizes correctness, regressions, reliability, maintainability, public-contract changes, and missing verification.

Findings should cite file and line evidence.

## Review a GitHub Pull Request

Use [`kevinle128-skills:review-pr`](../../skills/review-pr/SKILL.md) when the review target is already a GitHub PR.

```text
/kevinle128-skills:review-pr 123
```

Available options include:

- `--fix` repairs actionable findings and reruns the review loop.
- `--reply` posts the final review to GitHub.
- `--merge` merges only after review and CI gates pass.
- `--advice` adds advisory supervision around the verdict, repair loop, reply, and merge decision.

```text
/kevinle128-skills:review-pr 123 --fix --reply
```

## Use Git Operations Directly

Use [`kevinle128-skills:git`](../../skills/git/SKILL.md) for conventional commits, pushes, PR creation, merges, and stacked PRs.

```text
/kevinle128-skills:git cm
```

```text
/kevinle128-skills:git pr
```

The skill scans for secrets and can split unrelated changes into focused commits.

## Create a Pull Request

Use [`kevinle128-skills:git`](../../skills/git/SKILL.md) after implementation and review are complete.

```text
/kevinle128-skills:git cp
```

This command creates focused commits and pushes the current branch.

Create the pull request after the push succeeds.

```text
/kevinle128-skills:git pr
```

`kevinle128-skills:git` scans staged changes for secrets before it commits.

## Keep Plan and Project State Synchronized

[`kevinle128-skills:project-management`](../../skills/project-management/SKILL.md) hydrates tasks, reports status, synchronizes completed work back to plan files, and prepares cross-session handoffs.

[`kevinle128-skills:journal`](../../skills/journal/SKILL.md) records chronological work and decisions when a durable session history is useful.

Journals do not replace current architecture or product documentation.

Next: [Use specialized workflows](./06-specialized-workflows.md).
