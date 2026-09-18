---
name: bootstrap
description: "Bootstrap new projects with research, tech stack, design, planning, and implementation. Modes: full (default interactive), auto (explicit autonomous), fast (skip research), parallel (multi-agent)."
user-invocable: true
when_to_use: "Invoke to start a new project or full-stack setup from scratch."
category: utilities
keywords: [scaffold, project, setup, boilerplate]
license: MIT
argument-hint: "[requirements] [--full|--auto|--fast|--parallel] [--yagni] [--skip-journal]"
metadata:
  author: agentkit
  version: "1.0.0"
---

# Bootstrap - New Project Scaffolding

End-to-end project bootstrapping from idea to running code.

**Principles:** KISS, DRY | Full requested scope, nothing extra (`--yagni` to opt into scope-cutting) | Token efficiency | Concise reports

## Usage

```
/kevinle128-skills:bootstrap <user-requirements>
```

**Flags** (optional, default `--full`):

| Flag | Mode | Thinking | User Gates | Planning Skill | Cook Skill |
|------|------|----------|------------|----------------|------------|
| `--full` | Full interactive | Ultrathink | Every phase | `--hard` | (interactive) |
| `--auto` | Automatic explicit opt-in | Ultrathink | Design only | `--auto` | `--auto` |
| `--fast` | Quick | Think hard | Cook review gates | `--fast` | (interactive) |
| `--parallel` | Multi-agent | Ultrathink | Design only | `--parallel` | `--parallel` |

**Composable flags** (combine with any mode):

| Flag | Effect |
|------|--------|
| `--yagni` | Opt into YAGNI: challenge and cut scope not needed for the stated outcome (default: scaffold the full requested scope). Passed through to `kevinle128-skills:plan` and `kevinle128-skills:implement` |

**Example:**
```
/kevinle128-skills:bootstrap "Build a SaaS dashboard with auth" --fast
/kevinle128-skills:bootstrap "E-commerce platform with Stripe" --parallel
```

## Opening brainstorm gate (all modes)

Before Git initialization, research, design, planning, or scaffolding, capture:

- the intended product outcome;
- technology, safety, compatibility, and delivery constraints;
- explicit non-goals for this bootstrap;
- observable acceptance criteria for the running project.

Reuse an accepted brief or plan when it already contains these fields. Ask only
about a missing decision that would materially change the product or safety.
`--fast`, `--parallel`, and explicit `--auto` do not skip this gate; they only
change execution and approval behavior after the contract is concrete.

## Workflow Overview

```
[Brainstorm Contract] → [Git Init] → [Research?] → [Tech Stack?] → [Design?] → [Planning] → [Implementation] → [Test] → [Review] → [Docs] → [Onboard] → [Final]
```

Each mode loads a specific workflow reference + shared phases.

## Mode Detection

If no flag provided, default to `--full`.

Load the appropriate workflow reference:
- `--full`: Load `references/workflow-full.md`
- `--auto`: Load `references/workflow-auto.md` only when explicitly requested
- `--fast`: Load `references/workflow-fast.md`
- `--parallel`: Load `references/workflow-parallel.md`

All mode references inherit the opening brainstorm contract. Load
`references/shared-phases.md` for implementation through final report.

## Step 0: Git Init (ALL modes)

Check if Git initialized. If not:
- `--full`: Ask user if they want to init → `git-manager` subagent (`main` branch)
- Others: Auto-init via `git-manager` subagent (`main` branch)

## Skill Triggers (MANDATORY)

After early phases (research, tech stack, design), trigger downstream skills:

### Planning Phase
Activate **kevinle128-skills:plan** skill with mode-appropriate flag:
- `--full` → `/kevinle128-skills:plan --hard <requirements>` (thorough research + validation)
- `--auto` → `/kevinle128-skills:plan --auto <requirements>` (auto-detect complexity)
- `--fast` → `/kevinle128-skills:plan --fast <requirements>` (skip research)
- `--parallel` → `/kevinle128-skills:plan --parallel <requirements>` (file ownership + dependency graph)

Pass the brainstorm contract with the requirements so planning preserves the
accepted outcome, constraints, non-goals, and acceptance criteria.

Planning skill outputs a plan path. Pass this to cook.

### Implementation Phase
Activate **kevinle128-skills:implement** skill with the plan path and mode-appropriate flag:
- `--full` → `/kevinle128-skills:implement <plan-path>` (interactive review gates)
- `--auto` → `/kevinle128-skills:implement --auto <plan-path>` (explicit autonomous implementation)
- `--fast` → `/kevinle128-skills:implement <plan-path>` (skip extra research, keep cook review gates)
- `--parallel` → `/kevinle128-skills:implement --parallel <plan-path>` (multi-agent execution)

## Role

Elite software engineering expert specializing in system architecture and technical decisions. Brutally honest about feasibility and trade-offs.

## Critical Rules

- Activate relevant skills from catalog during the process
- Keep all research reports ≤150 lines
- All docs written to `./docs` directory
- Plans written to `./plans` directory using naming from `## Naming` section
- DO NOT implement code directly — delegate through planning + cook skills
- Sacrifice grammar for concision in reports
- List unresolved questions at end of reports
- Run `/kevinle128-skills:journal` to write a concise technical journal entry upon completion — unless the shared "Journal step — opt-out" below applies.

### Journal step — opt-out

Skip the automatic `/kevinle128-skills:journal` step when either applies:
- The invocation includes the `--skip-journal` flag, OR
- `ak config prefs resolve --json | jq -r 'if .prefs.journal.auto == false then "false" else "true" end'` returns `false`. If the command errors or prints anything other than the exact string `false`, treat as `true` (default) — corrupt or missing config never suppresses the automatic journal.

Precedence: flag > project config > user config > default (`true`).
When skipped, print one line:
- `journal skipped by --skip-journal` (flag), or
- `journal skipped by preference` (config).

Explicit `/kevinle128-skills:journal` and `ak journal create` are unaffected.

## References

- `references/workflow-full.md` - Full interactive workflow
- `references/workflow-auto.md` - Explicit auto workflow
- `references/workflow-fast.md` - Fast workflow
- `references/workflow-parallel.md` - Parallel workflow
- `references/shared-phases.md` - Common phases (implementation → final report)
