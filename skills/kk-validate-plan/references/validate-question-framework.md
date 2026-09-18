# Validation Question Framework

## Question Categories

| Category | Keywords to detect |
|----------|-------------------|
| **Architecture** | "approach", "pattern", "design", "structure", "database", "API" |
| **Assumptions** | "assume", "expect", "should", "will", "must", "default" |
| **Tradeoffs** | "tradeoff", "vs", "alternative", "option", "choice", "either/or" |
| **Risks** | "risk", "might", "could fail", "dependency", "blocker", "concern" |
| **Scope** | "phase", "MVP", "future", "out of scope", "nice to have" |
| **Runtime Surface** | "trigger", "route", "UI", "command", "cron", "listener", "event", "callback", "scheduler", "consumer" |

## Question Format Rules

- Each general question must have 2-4 concrete options
- Each missing-runtime-surface question must have 2-3 mutually exclusive, codebase-grounded options
- Put the recommended option first and mark it with the "(Recommended)" suffix
- "Other" option is automatic
- Questions should surface implicit decisions
- A `NEEDS_DECISION` Runtime Flow Proof row always requires `ask_user capability`, even when the configured question count is zero or already reached
- Do not print a required HITL question as prose and then continue the workflow

## Example Questions

Category: Architecture
Question: "How should the validation results be persisted?"
Options:
1. Save to plan.md frontmatter (Recommended)
2. Create validation-answers.md
3. Don't persist

Category: Assumptions
Question: "The plan assumes API rate limiting is not needed. Is this correct?"
Options:
1. Yes, not needed for MVP
2. No, add basic rate limiting now (Recommended)
3. Defer to Phase 2

Category: Runtime Surface
Question: "The requested model change has no runtime entry point. Which surface should own it?"
Options:
1. Add the existing session API pattern and test it through HTTP (Recommended)
2. Add a CLI command and test the registered command
3. Remove user-controlled model changes from the accepted scope

## Validation Log Format

```markdown
## Validation Log

### Session {N} — {YYYY-MM-DD}
**Trigger:** {what prompted this validation}
**Questions asked:** {count}

### Runtime Flow Proof Matrix

| Feature | Actor | Runtime trigger | Entry point | Internal path | Observable result | End-to-end test | External mocks | Prepared data | Status |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| ... | ... | ... | ... | ... | ... | ... | ... | ... | PASSED / FAILED / NEEDS_DECISION |

- **Gate status:** PASSED / FAILED / NEEDS_DECISION

#### Questions & Answers

1. **[{Category}]** {full question text}
   - Options: {A} | {B} | {C}
   - **Answer:** {user's choice}
   - **Custom input:** {verbatim "Other" text if applicable}
   - **Rationale:** {why this decision matters}

#### Confirmed Decisions
- {decision}: {choice} — {brief why}

#### Action Items
- [ ] {specific change needed}

#### Impact on Phases
- Phase {N}: {what needs updating and why}
```

## Recording Rules

- **Full question text**: exact question, not summary
- **All options**: every option presented
- **Verbatim custom input**: record "Other" text exactly
- **Rationale**: explain why decision affects implementation
- **Session numbering**: increment from last session
- **Trigger**: state what prompted validation

## Section Mapping for Phase Propagation

| Change Type | Target Section |
|-------------|----------------|
| Requirements | Requirements |
| Architecture | Architecture |
| Scope | Overview / Implementation Steps |
| Risk | Risk Assessment |
| Runtime Surface | Requirements / Architecture / Implementation Steps / Success Criteria |
| End-to-End Proof | Implementation Steps / Success Criteria / Risk Assessment |
| Unknown | Key Insights (new subsection) |
