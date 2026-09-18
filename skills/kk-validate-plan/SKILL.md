---
name: kk:validate-plan
description: "Validate an implementation plan against the real codebase, prove each feature through its runtime trigger and end-to-end test, resolve material decisions with the user, and reconcile the whole plan before implementation."
user-invocable: true
when_to_use: "Invoke after planning and before implementation when a plan needs evidence-backed validation."
category: utilities
keywords: [plan, validation, interview, verification, consistency, runtime-flow, end-to-end]
argument-hint: "[plan-directory-or-plan.md-path]"
license: MIT
metadata:
  author: KevinKit
  version: "1.2.0"
---

# Validate Plan

Interview the user with critical questions to validate assumptions, confirm decisions, and surface potential issues in an implementation plan before coding begins.

## Plan Resolution

1. If `$ARGUMENTS` provided → Use that path
2. Else check `## Plan Context` section → Use active plan path
3. If no plan found → Ask user to specify path or run `/kk:plan --hard` first

## Configuration

Check `## Plan Context` section for validation settings:
- `mode` - Controls auto/prompt/off behavior
- `questions` - Range like `3-8` (min-max)

## Workflow

### Step 1: Read Plan Files
- `plan.md` - Overview and phases list
- `phase-*.md` - All phase files
- Look for decision points, assumptions, risks, tradeoffs

### Step 2: Extract Question Topics
Load: `references/validate-question-framework.md`

### Step 2.5: Runtime Flow Proof and Verification Pass

Before interviewing the user, verify plan accuracy against the actual codebase.
Load: `references/verification-roles.md`

#### Mandatory Runtime Flow Proof Gate

Run the `Mandatory Runtime Flow Proof Gate` for every plan.
This gate is not sampled and cannot be skipped because the plan is small or already has red-team evidence.

1. Extract every distinct feature and acceptance criterion promised by the plan.
2. For each feature, identify the real actor and the existing or planned runtime trigger.
3. Trace the complete planned runtime path from that trigger through real internal components to an observable result.
4. Verify existing path segments against source and verify that every new segment has a code-level implementation task with a clear owner.
5. Verify that the plan includes an end-to-end test which enters through the same runtime boundary.
6. Verify that the test prepares realistic data, keeps internal services real, and mocks only third-party systems at their boundary.
7. Record one row per feature in the `Runtime Flow Proof Matrix` defined in `references/verification-roles.md`.

Do not accept a helper-level, handler-level, repository-level, or component-only test as end-to-end proof when a higher runtime boundary exists.
One scenario can cover several features only when its steps and assertions prove each feature explicitly.

If neither the existing system nor the plan provides a runtime trigger, public surface, scheduler registration, event ingress, or owning caller for a required feature, stop and open a HITL question with `ask_user capability`.
Also open HITL when the requirement does not choose a surface and the plan selected one without a recorded user decision.
Offer two or three codebase-grounded options, put the recommended option first, and suffix it with `(Recommended)`.
Do not invent a test-only trigger, silently reduce scope, or select a surface for the user.
Record the result as `NEEDS_DECISION` until the user answers.

After the answer, update the affected plan phases at code level.
Include the selected entry surface, owner, registration or routing, internal call path, data setup, third-party mocks, end-to-end test steps, assertions, and acceptance criteria.

#### Auto-Scaled Supporting Verification

**Red-team reuse guard:** If `plan.md` already contains a `## Red Team Review` with verification evidence, reuse that evidence and only resolve remaining `[UNVERIFIED]` claims.
The red-team guard never skips the Mandatory Runtime Flow Proof Gate.

1. **Tier detection** — Count phases in the plan:
   - 1-2 phases → Light (Fact Checker only, 5 claims/phase)
   - 3-4 phases → Standard (Fact Checker + Contract Verifier, 10 claims/phase)
   - 5+ phases → Full (all 4 roles, 15+ claims/phase)
2. **For each active role at the current tier:**
   - Sample N claims per phase (per tier budget)
   - Run grep/glob to verify file paths, symbols, endpoints
   - Collect findings: VERIFIED | FAILED | UNVERIFIED
3. **Handle failures:**
   - Surface ALL failures as additional interview questions in Step 4 (with glob-suggested alternatives as "(Recommended)" options)
   - Never auto-correct plan files — all corrections require user confirmation via interview
4. **Check `[UNVERIFIED]` tags** — Scan plan for planner-tagged unverified claims, attempt to resolve
5. **Append results** to `## Validation Log` after the Runtime Flow Proof Matrix:
   ```
   ### Verification Results
   - Claims checked: N
   - Verified: N | Failed: N | Unverified: N
   - Tier: Light|Standard|Full
   - Failures: [list with file:line evidence]
   ```

### Step 3: Generate Questions
For each detected topic, formulate a concrete question with 2-4 options.
Mark recommended option with "(Recommended)" suffix.
Runtime-surface questions are mandatory and do not count against the configured question range.

### Step 4: Interview User
Use `ask_user capability` tool.
- Use question count from `## Plan Context` validation settings
- Group related questions (max 4 per tool call)
- Focus on: missing runtime surfaces, assumptions, risks, tradeoffs, architecture
- Use a popup-capable `ask_user capability` call for every `NEEDS_DECISION` runtime-flow proof row
- Do not replace required HITL with a prose question followed by continued execution

### Step 5: Document Answers
Add or append `## Validation Log` section in `plan.md`.
Load: `references/validate-question-framework.md` for recording format.

### Step 6: Propagate Changes to Phases
Auto-propagate validation decisions to affected phase files.
Add marker: `<!-- Updated: Validation Session N - {change} -->`

### Step 7: Whole-Plan Consistency Sweep
Load: `references/verification-roles.md` → "Whole-Plan Consistency Sweep".

After propagation, re-read `plan.md` and every `phase-*.md` file. Check the whole plan for stale or contradictory claims caused by the validation decisions.

Required checks:
- Search all plan files for old terms, renamed fields/APIs/files, rejected assumptions, and superseded validation decisions.
- Reconcile `plan.md` overview, phase summaries, implementation steps, success criteria, risk notes, and validation logs.
- If the same SQL/query/API/body/contract appears as both prose and embedded draft, update both copies or mark the unresolved conflict.
- Append `### Whole-Plan Consistency Sweep` to the current `## Validation Log`.
- If any unresolved contradiction remains, ask the user before recommending implementation.

## Output
- Number of questions asked
- Key decisions confirmed
- Runtime Flow Proof Matrix and gate status
- Phase propagation results
- Whole-plan consistency sweep results
- Recommendation: proceed or revise

## Next Steps
Present user-choice next steps with the absolute path:
> **Best Practice:** Run `/clear` before implementing to start with fresh context.
> If the user chooses implementation, run:
> ```
> /kk:implement {ABSOLUTE_PATH_TO_PLAN_DIR}/plan.md
> ```
> **Flag selection:** The plan is eligible for implementation only when Runtime Flow Proof is `PASSED`, verification shows `Failed: 0`, no blocking claim is `UNVERIFIED`, no HITL decision is open, and the consistency sweep has zero unresolved contradictions. Ask the user before proceeding. Add `--auto` only when the user explicitly asks for autonomous implementation.
> **Why absolute path?** After `/clear`, the new session loses previous context.
> Fresh context helps Claude focus solely on implementation without planning context pollution.

## Important Notes
- Only ask about genuine decision points
- If plan is simple, fewer than min questions is okay
- Prioritize questions that could change implementation significantly
- Never waive a missing runtime trigger because the feature is internal or has no direct user interface
- For cron jobs, listeners, consumers, callbacks, and other system-triggered features, treat the real scheduler, event ingress, queue, or owning caller as the actor boundary
- Never recommend implementation until Runtime Flow Proof passes and the whole-plan consistency sweep has no unresolved contradictions
