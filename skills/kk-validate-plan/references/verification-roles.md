# Verification Roles

Language-agnostic roles for verifying plan accuracy against the actual codebase.

**Loaded by:** `kk:plan` for self-verification, `kk:validate-plan` for external verification, and `kk:plan red-team` for evidence-backed adversarial review.

**Principle:** `user asks → scout → planner writes → audit-verify → report` instead of `user asks → planner writes → report done`.

## Mandatory Runtime Flow Proof Gate

**Purpose:** Prove that the plan delivers each requested feature through a real runtime trigger and that its test proves the connected behavior, not isolated parts.

Runtime means the real operational path by which a user, system, scheduler, event source, or owning caller invokes the behavior.
It does not mean a deployed production environment.

This gate runs for every plan before tiered spot checks.
It is complete coverage, not a sample.

### Build the Feature Inventory

Extract each distinct feature from the requested outcome, plan scope, acceptance criteria, phase requirements, and success criteria.
Do not collapse separate controls, modes, operations, or failure behaviors into one vague row.

For each feature, answer these questions from the accepted requirements, plan, and source:

1. Who or what triggers the feature at runtime?
2. What existing or planned runtime boundary receives that trigger?
3. Which real internal components carry the request or event, and which components must the plan add or change?
4. What externally observable result proves the feature worked?
5. Which test enters through that runtime boundary and observes that result?

The feature can use a new runtime boundary that does not exist yet.
The accepted requirement or a recorded user decision must choose its surface.
The plan must then name the code owner, registration or routing change, contract, internal wiring, and end-to-end test.
Verify existing path segments against source and mark new segments as planned work.

### Recognized Runtime Boundaries

| Feature type | Required test entry |
| --- | --- |
| HTTP or RPC API | Start the real application test server and call the registered route through the protocol boundary. |
| Website | Open the running website in a browser and operate the visible UI as a user. |
| CLI or TUI | Invoke the registered command or interactive entry as a user would. |
| Cron or scheduled job | Trigger the registered scheduler or job boundary with controlled time and assert the external effect. |
| Listener, consumer, webhook, or event handler | Deliver the event through the real ingress, broker adapter, or registered listener and assert the external effect. |
| Internal library or callback | Invoke the nearest real runtime caller or public module boundary and prove the observable consumer result. |

A direct call to a helper, handler method, repository, reducer, or component is not end-to-end proof when a higher runtime boundary exists.
A unit or component test can supplement the proof but cannot replace it.

### Full-Flow Test Contract

The planned test must state:

- the actor and runtime trigger;
- the exact application or system entry point;
- realistic prepared data and required initial state;
- the ordered internal path through real services;
- third-party boundaries that will be mocked;
- the fault, crash, retry, or timing action when the acceptance case requires it;
- the observable response, persisted state, emitted event, rendered UI, or runtime behavior;
- the assertions that prove every feature covered by the scenario.

Mock third-party systems at their boundary.
Do not mock internal services that own part of the feature path.
Use their real integration path and prepare the required data.

### Missing Trigger or Surface

Search the source before concluding that a surface is missing.
Check routes, command registration, UI actions, scheduler registration, event subscriptions, queue consumers, callbacks, and public module callers.

Open HITL when any of these conditions is true:

- neither the current system nor the plan provides a runtime trigger for a required feature;
- the requested outcome needs a user or system surface, but the plan only describes internal helpers;
- the requirement leaves the surface open and the plan selected one without a recorded user decision;
- the specification and the proposed surface conflict;
- several credible surfaces would produce materially different public contracts.

When HITL is required:

1. Mark the matrix row `NEEDS_DECISION`.
2. Open a HITL question with `ask_user capability` before editing the plan or continuing to readiness.
3. Offer two or three mutually exclusive, codebase-grounded options.
4. Put the recommended option first and add the `(Recommended)` suffix.
5. Explain which runtime owner and end-to-end test each option creates or reuses.
6. Do not create a test-only route or infer that the feature is intentionally unreachable.

If the specification requires user behavior but forbids every usable surface, present that conflict as the HITL question.
The user must select the intended contract.

### Runtime Flow Proof Matrix

Append this matrix to the current `## Validation Log`:

```markdown
### Runtime Flow Proof Matrix

| Feature | Actor | Runtime trigger | Entry point | Internal path | Observable result | End-to-end test | External mocks | Prepared data | Status |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| ... | ... | ... | ... | ... | ... | ... | ... | ... | PASSED / FAILED / NEEDS_DECISION |

- **Gate status:** PASSED / FAILED / NEEDS_DECISION
```

`PASSED` requires a codebase-grounded runtime path, explicit implementation tasks for every new segment, and a planned end-to-end test that satisfies the Full-Flow Test Contract.
`FAILED` means the plan has a concrete logic, wiring, or test gap that can be repaired without a product decision.
`NEEDS_DECISION` means the missing or conflicting runtime surface requires the user to choose.

The plan cannot proceed to implementation unless every row and the gate are `PASSED`.

## Tiering (Auto-Scale by Plan Size)

Count phases in the plan to determine the supporting verification tier.
The Mandatory Runtime Flow Proof Gate remains active at every tier.

| Phases | Tier | Active Roles | Spot-Check Budget |
|--------|------|-------------|-------------------|
| 1-2 | Light | Runtime Flow Verifier + Fact Checker | Complete proof matrix + 5 claims/phase |
| 3-4 | Standard | Runtime Flow Verifier + Fact Checker + Contract Verifier | Complete proof matrix + 10 claims/phase |
| 5+ | Full | Runtime Flow Verifier + all 4 supporting roles | Complete proof matrix + 15+ claims/phase |

## Role: Fact Checker

**Purpose:** Verify every file path, symbol, endpoint, and config key cited in the plan actually exists.

**Method:**
- Sample N claims per phase (per tier budget)
- `grep -rn "{symbol}" .` to verify symbols exist
- `glob "{path}"` to verify file paths
- For endpoints: grep route definitions
- For config keys: grep env files, config objects

**Red flags:**
- Wrapper/validator/manager/handler names that grep returns nothing for
- Centralized packages that are actually scattered across the codebase
- Paths from scout reports that were renamed or moved since scouting

**Output per claim:** `VERIFIED (file:line)` | `FAILED (not found)` | `UNVERIFIED (ambiguous)`

## Role: Flow Tracer

**Purpose:** Verify behavioral claims ("X triggers Y", "A calls B before C", "middleware runs before handler").

**Method:**
- Start from the claimed entry point
- Read the actual code path: entry → guards → branching → target
- List all early returns, middleware chains, event listeners in the path
- For async code: check ordering guarantees (await, Promise.then, callbacks)
- Verify causality (A actually invokes B) vs correlation (both exist in same file)

**Red flags:**
- "X triggers Y" but X and Y share no call path
- Missing intermediate steps (A calls C which calls B, not A calls B directly)
- Async ordering assumed synchronous

**Output:** Traced path with file:line citations, or FAILED with explanation of actual flow.

## Role: Scope Auditor

**Purpose:** Verify state additions (new fields, context values, singletons, env vars) respect lifetime boundaries.

**Method:**
- Grep the target struct/class/object for ALL instantiation sites
- Determine lifetime: request-scoped, session-scoped, process-global
- Check for shared-state leaks across isolation boundaries
- Verify no existing state already serves the same purpose (grep for similar field names)

**Red flags:**
- "Adding field to X" when X is a singleton shared across requests
- New state duplicating existing state under a different name
- Module-level variables in request-handling code

**Output:** Lifetime classification with instantiation sites, or FAILED with leak description.

## Role: Contract Verifier

**Purpose:** Verify interface changes (API endpoints, function signatures, config schemas, exports) account for ALL consumers.

**Method:**
- `grep -rn "{function_name}" .` to enumerate ALL callers — list explicitly
- Never write "update all callers" — always state the count and list them
- If count > 10: list first 10 with file:line, state total count
- Check downstream: tests that call the function, imports, re-exports
- Check upstream: config files, env vars, CI scripts, CLI help text

**Red flags:**
- Plan says "3 callers" but grep finds 7
- Missing test file updates
- Re-exported types not updated at barrel files
- CLI help text referencing old parameter names

**Output:** Caller list with file:line, type compatibility assessment, or FAILED with missing callers.

## Verification Output Format

Append to plan's `## Validation Log` or include in red-team findings:

```markdown
### Verification Results
- **Tier:** Light|Standard|Full
- **Claims checked:** N
- **Verified:** N | **Failed:** N | **Unverified:** N

#### Failures
1. [Fact Checker] `src/utils/auth.ts` — path not found, actual: `src/lib/auth.ts`
2. [Contract Verifier] `parseConfig()` — plan says 3 callers, found 7
```

## Whole-Plan Consistency Sweep

**Purpose:** Prevent iterative validate/red-team edits from fixing one phase while leaving stale claims elsewhere.

Run this after any validation or red-team change that edits `plan.md` or any `phase-*.md` file.

### Required Inputs

- `plan.md`
- Every `phase-*.md` file in the plan directory
- New decisions or accepted findings from the current validation/red-team session

### Sweep Method

1. Re-read `plan.md` and all `phase-*.md` files after applying edits.
2. Build a short decision delta list from the current session:
   - renamed fields, APIs, files, tags, timestamps, scopes, or workflows
   - changed validation decisions or rejected assumptions
   - changed phase order, dependencies, ownership, or success criteria
3. Search all plan files for old terms, superseded assumptions, and duplicate embedded drafts from each delta.
4. Reconcile affected sections across files, not only the file that triggered the finding.
5. Check `plan.md` summary, phases table text, phase requirements, implementation steps, success criteria, risk notes, and validation/red-team logs for contradictions.
6. Re-run affected Runtime Flow Proof Matrix rows when a decision changes a trigger, path, test, or observable result.
7. If a conflict cannot be resolved with current evidence, add it to unresolved questions and do not recommend implementation yet.

### Output Format

Append to the current `## Validation Log` or `## Red Team Review` section:

```markdown
### Whole-Plan Consistency Sweep
- Files reread: plan.md, phase-01-..., phase-02-...
- Decision deltas checked: N
- Reconciled stale references: N
- Unresolved contradictions: N
```

If `Unresolved contradictions` is greater than zero, list each conflict with the affected files and ask the user before implementation.

## Integration Points

- **Planner readiness gate:** Run the complete Mandatory Runtime Flow Proof Gate after drafting the plan and before red-team review, task hydration, publication, or implementation handoff. Repair `FAILED` rows and resolve every `NEEDS_DECISION` row through popup HITL.
- **Planner fact checks:** Apply Fact Checker inline while writing each phase. Tag `[UNVERIFIED]` for claims that cannot be confirmed.
- **`kk:validate-plan` Step 2.5:** Re-run the complete Mandatory Runtime Flow Proof Gate, then run tier-appropriate supporting roles before general interview questions.
  FAILED findings become additional interview topics.
- **`kk:validate-plan` HITL:** Resolve every `NEEDS_DECISION` runtime-flow row with the user before readiness can pass.
- **`kk:validate-plan` after propagation:** Run the Whole-Plan Consistency Sweep before recommending implementation.
- **Red-team reviewers:** Each reviewer carries their adversarial lens + assigned verification role. Findings without grep evidence = auto-rejected during adjudication.
- **Red-team after accepted edits:** Run Whole-Plan Consistency Sweep before presenting next steps.
