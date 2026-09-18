# Recipes and Pitfalls

These prompts are starting points.

Replace the examples with your real outcome, constraints, and finish condition.

## Plan a Production API Feature

```text
/kevinle128-skills:plan Add a public API that lets a user change the active model, thinking mode, and fast mode.
Preserve existing session behavior.
Done means an end-to-end API test updates the values, the live worker applies them after commit and wake, and stale callbacks cannot overwrite the new state.
Use --hard --tdd.
```

## Plan and Implement a Web Feature

```text
/kevinle128-skills:plan Add model controls to the session settings page --hard --tdd
Done means a browser test signs in, changes each control, reloads the page, and proves the live session uses the saved values.
```

After plan approval:

```text
/kevinle128-skills:validate-plan /absolute/path/to/plan.md
```

After validation passes:

```text
/kevinle128-skills:implement /absolute/path/to/plan.md --tdd
```

## Fix a Production Regression

```text
/kevinle128-skills:fix Users receive duplicate notifications after a retry.
Reproduce through the public event flow first.
Prove the shared root cause, preserve normal delivery, and add a regression test for the full retry sequence.
```

## Diagnose Without Editing

```text
/kevinle128-skills:debug Explain why an idle worker applies an older session generation after wake.
Do not change code.
Return the exact event order, state owner, stale-write point, and evidence for the root cause.
```

## Port a Feature from Another Repository

```text
/kevinle128-skills:feature-lens owner/source-repository "retry strategy" --port
```

Feature Lens treats the source repository as untrusted evidence, maps the feature to local architecture, challenges incompatible assumptions, and hands the approved result to `kevinle128-skills:plan`.

## Execute Independent Phases in Parallel

```text
/kevinle128-skills:plan Add the API, admin UI, and audit export --parallel --tdd
```

```text
/kevinle128-skills:implement /absolute/path/to/plan.md --parallel --tdd
```

The plan must assign exclusive files and order the integration phase after its producers.

## Run Autonomously with a Finish Condition

```text
/kevinle128-skills:implement /absolute/path/to/plan.md --auto
Continue until every acceptance criterion has direct evidence, all affected tests pass, and independent review has no blocking findings.
Stop if a product decision or irreversible action is required.
```

Autonomous execution is safe only when the finish condition can pass or fail objectively.

## Test a Logged-In Browser Flow

```text
/kevinle128-skills:chrome-profile Use my work profile to open the settings page and confirm the saved model is shown.
```

Use `kevinle128-skills:agent-browser` instead when the test should start from a clean browser state.

## Review and Create a Pull Request

```text
/kevinle128-skills:code-review --pending
```

```text
/kevinle128-skills:git cp
```

After the push succeeds:

```text
/kevinle128-skills:git pr
```

## Explain an Architecture Visually

```text
/kevinle128-skills:mermaidjs-v11 Draw the user request, commit, worker wake, dispatch, and callback validation sequence.
```

Use `kevinle128-skills:tech-graph` when the result must be a polished SVG and PNG artifact.

Use `kevinle128-skills:preview` when you want an interactive explanation or visual review.

## Pitfalls

### Starting Implementation Without an Observable Trigger

A feature cannot be proven end to end if no user, API, command, event, listener, or scheduler can trigger it.

Stop during planning and resolve the missing production surface instead of inventing a test-only path.

### Testing Components but Not the Joined Flow

Unit tests may prove each helper while the complete feature still fails at wiring, ordering, persistence, or ownership boundaries.

Keep focused tests, then add one test that acts as the user or system and traverses the whole flow.

### Mocking Internal Services

Mock external third parties at their boundary.

Use real internal services and realistic prepared data for integration and end-to-end proof.

### Using `--fast` to Reduce Scope

`--fast` reduces research and process depth.

It does not remove requested behavior.

Only `--yagni` or an explicit user decision authorizes scope reduction.

### Treating `--auto` as Permission for Unsafe Actions

`--auto` skips routine approval pauses.

It does not bypass security, test, review, public-contract, or irreversible-action gates.

### Running Parallel Agents on Shared Files

Parallel mode needs exclusive ownership and explicit dependencies.

Do not use it when workers would edit the same file, migration sequence, generated artifact, or shared configuration.

### Skipping Tests Silently

`--no-test` and `--skip-tests` are explicit risk decisions.

The final report must state what was not verified and why.

### Confusing Skills with the AgentKit CLI

Kevin Kit skills use the `kevinle128-skills:*` namespace.

The AgentKit executable remains `ak`, so commands such as `ak plan status` and `ak config prefs resolve` do not change names.

### Calling Every Supporting Skill Manually

Choose the owning workflow and let it route supporting skills.

Call a supporting skill directly only when its output is the requested deliverable or when you intentionally override the default route.

Back to [The Kevin Kit Guide](./README.md).
