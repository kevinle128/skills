# Skill Catalog

This catalog lists every skill currently shipped in Kevin Kit.

The linked `SKILL.md` file is the authority for arguments, gates, and detailed behavior.

## Core Delivery

| Skill | Use it when |
| --- | --- |
| [`kevinle128-skills:bootstrap`](../../skills/bootstrap/SKILL.md) | Start a new project from requirements and carry it through planning, implementation, testing, and setup. |
| [`kevinle128-skills:brainstorm`](../../skills/brainstorm/SKILL.md) | Turn unclear intent into an outcome, constraints, non-goals, acceptance criteria, and a chosen direction. |
| [`kevinle128-skills:plan`](../../skills/plan/SKILL.md) | Create an architecture or implementation plan with phases, dependencies, review, and validation. |
| [`kevinle128-skills:validate-plan`](../../skills/validate-plan/SKILL.md) | Verify an existing plan against the codebase, resolve material questions, and reconcile the full plan before implementation. |
| [`kevinle128-skills:implement`](../../skills/implement/SKILL.md) | Execute an accepted plan or clearly defined feature through implementation, testing, review, and finalization. |
| [`kevinle128-skills:project-management`](../../skills/project-management/SKILL.md) | Hydrate tasks, report progress, synchronize plan state, or prepare a handoff. |
| [`kevinle128-skills:journal`](../../skills/journal/SKILL.md) | Record a chronological technical work log or session reflection. |

## Investigation and Repair

| Skill | Use it when |
| --- | --- |
| [`kevinle128-skills:scout`](../../skills/scout/SKILL.md) | Locate files, symbols, owners, callers, tests, or repository patterns quickly. |
| [`kevinle128-skills:debug`](../../skills/debug/SKILL.md) | Prove a root cause for a bug, failed test, unexpected behavior, performance problem, or CI failure. |
| [`kevinle128-skills:fix`](../../skills/fix/SKILL.md) | Repair a concrete bug or failure through reproduction, diagnosis, implementation, and verification. |
| [`kevinle128-skills:feature-lens`](../../skills/feature-lens/SKILL.md) | Compare, copy, improve, or idiomatically port a feature from another repository. |
| [`kevinle128-skills:problem-solving`](../../skills/problem-solving/SKILL.md) | Reframe a problem after complexity, assumptions, or repeated failed attempts block progress. |
| [`kevinle128-skills:sequential-thinking`](../../skills/sequential-thinking/SKILL.md) | Work through a complex problem with explicit steps, hypothesis checks, and revisions. |
| [`kevinle128-skills:docs-seeker`](../../skills/docs-seeker/SKILL.md) | Retrieve current library, framework, API, or repository documentation. |
| [`kevinle128-skills:repomix`](../../skills/repomix/SKILL.md) | Pack a repository into an AI-friendly analysis artifact. |
| [`kevinle128-skills:context-engineering`](../../skills/context-engineering/SKILL.md) | Diagnose or design context budgets, memory systems, token use, or multi-agent context flow. |

## Testing, Review, and Delivery

| Skill | Use it when |
| --- | --- |
| [`kevinle128-skills:test`](../../skills/test/SKILL.md) | Run or design unit, integration, end-to-end, UI, coverage, build, or QA checks. |
| [`kevinle128-skills:web-testing`](../../skills/web-testing/SKILL.md) | Build Playwright, Vitest, k6, accessibility, visual, load, or cross-browser tests. |
| [`kevinle128-skills:code-review`](../../skills/code-review/SKILL.md) | Review pending changes, a commit, a PR, or a codebase for bugs and regressions. |
| [`kevinle128-skills:review-pr`](../../skills/review-pr/SKILL.md) | Review a GitHub PR and optionally fix findings, post the review, or merge after CI. |
| [`kevinle128-skills:git`](../../skills/git/SKILL.md) | Create focused conventional commits, push branches, open PRs, merge, or manage stacked PRs. |

## Browser, UI, and Media

| Skill | Use it when |
| --- | --- |
| [`kevinle128-skills:agent-browser`](../../skills/agent-browser/SKILL.md) | Automate a clean browser, cloud browser, Electron app, form, screenshot, scrape, or exploratory QA flow. |
| [`kevinle128-skills:chrome-profile`](../../skills/chrome-profile/SKILL.md) | Automate a real Chrome profile with the user's cookies, accounts, or authenticated tabs. |
| [`kevinle128-skills:frontend-design`](../../skills/frontend-design/SKILL.md) | Build or reproduce a polished frontend with strong visual fidelity. |
| [`kevinle128-skills:ui-ux-pro-max`](../../skills/ui-ux-pro-max/SKILL.md) | Make UX, design-system, accessibility, typography, layout, interaction, or responsive design decisions. |
| [`kevinle128-skills:threejs`](../../skills/threejs/SKILL.md) | Build a Three.js, WebGL, WebGPU, GLTF, physics, animation, VR, or XR experience. |
| [`kevinle128-skills:ai-multimodal`](../../skills/ai-multimodal/SKILL.md) | Analyze or generate images, audio, video, and documents, including OCR and transcription. |
| [`kevinle128-skills:media-processing`](../../skills/media-processing/SKILL.md) | Transform media with FFmpeg, ImageMagick, background removal, batch jobs, or streaming formats. |

## Documentation and Visualization

| Skill | Use it when |
| --- | --- |
| [`kevinle128-skills:docs`](../../skills/docs/SKILL.md) | Create, update, summarize, or audit project documentation or agent context. |
| [`kevinle128-skills:folder-context`](../../skills/folder-context/SKILL.md) | Add compact agent instructions for one subfolder. |
| [`kevinle128-skills:project-organization`](../../skills/project-organization/SKILL.md) | Choose paths, organize assets, or standardize a project layout. |
| [`kevinle128-skills:preview`](../../skills/preview/SKILL.md) | Preview files or generate visual explanations, slides, diagrams, recaps, or plan reviews. |
| [`kevinle128-skills:mermaidjs-v11`](../../skills/mermaidjs-v11/SKILL.md) | Create editable inline diagrams with Mermaid v11. |
| [`kevinle128-skills:tech-graph`](../../skills/tech-graph/SKILL.md) | Create publish-grade SVG and PNG architecture or flow diagrams. |

## Quick Selection Rules

- Use `kevinle128-skills:plan` when the work needs phases or architecture.
- Use `kevinle128-skills:validate-plan` when an existing plan needs evidence-backed validation before implementation.
- Use `kevinle128-skills:implement` when scope is accepted and code must change.
- Use `kevinle128-skills:fix` when a concrete failure must be repaired.
- Use `kevinle128-skills:debug` when diagnosis is the requested deliverable.
- Use `kevinle128-skills:test` when verification is the requested deliverable.
- Use `kevinle128-skills:code-review` for a local diff and `kevinle128-skills:review-pr` for an existing GitHub PR.
- Use `kevinle128-skills:agent-browser` for clean automation and `kevinle128-skills:chrome-profile` for real logged-in browser state.
- Use `kevinle128-skills:mermaidjs-v11` for editable inline diagrams and `kevinle128-skills:tech-graph` for polished image output.

Next: [Use recipes and avoid pitfalls](./08-recipes-and-pitfalls.md).
