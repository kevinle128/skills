# Skill Catalog

This catalog lists every skill currently shipped in KevinKit.

The linked `SKILL.md` file is the authority for arguments, gates, and detailed behavior.

## Core Delivery

| Skill | Use it when |
| --- | --- |
| [`kk:bootstrap`](../../skills/kk-bootstrap/SKILL.md) | Start a new project from requirements and carry it through planning, implementation, testing, and setup. |
| [`kk:brainstorm`](../../skills/kk-brainstorm/SKILL.md) | Turn unclear intent into an outcome, constraints, non-goals, acceptance criteria, and a chosen direction. |
| [`kk:plan`](../../skills/kk-plan/SKILL.md) | Create an architecture or implementation plan with phases, dependencies, review, and validation. |
| [`kk:validate-plan`](../../skills/kk-validate-plan/SKILL.md) | Verify an existing plan against the codebase, resolve material questions, and reconcile the full plan before implementation. |
| [`kk:implement`](../../skills/kk-implement/SKILL.md) | Execute an accepted plan or clearly defined feature through implementation, testing, review, and finalization. |
| [`kk:project-management`](../../skills/kk-project-management/SKILL.md) | Hydrate tasks, report progress, synchronize plan state, or prepare a handoff. |
| [`kk:journal`](../../skills/kk-journal/SKILL.md) | Record a chronological technical work log or session reflection. |

## Investigation and Repair

| Skill | Use it when |
| --- | --- |
| [`kk:scout`](../../skills/kk-scout/SKILL.md) | Locate files, symbols, owners, callers, tests, or repository patterns quickly. |
| [`kk:debug`](../../skills/kk-debug/SKILL.md) | Prove a root cause for a bug, failed test, unexpected behavior, performance problem, or CI failure. |
| [`kk:fix`](../../skills/kk-fix/SKILL.md) | Repair a concrete bug or failure through reproduction, diagnosis, implementation, and verification. |
| [`kk:feature-lens`](../../skills/kk-feature-lens/SKILL.md) | Compare, copy, improve, or idiomatically port a feature from another repository. |
| [`kk:problem-solving`](../../skills/kk-problem-solving/SKILL.md) | Reframe a problem after complexity, assumptions, or repeated failed attempts block progress. |
| [`kk:sequential-thinking`](../../skills/kk-sequential-thinking/SKILL.md) | Work through a complex problem with explicit steps, hypothesis checks, and revisions. |
| [`kk:docs-seeker`](../../skills/kk-docs-seeker/SKILL.md) | Retrieve current library, framework, API, or repository documentation. |
| [`kk:repomix`](../../skills/kk-repomix/SKILL.md) | Pack a repository into an AI-friendly analysis artifact. |
| [`kk:context-engineering`](../../skills/kk-context-engineering/SKILL.md) | Diagnose or design context budgets, memory systems, token use, or multi-agent context flow. |

## Testing, Review, and Delivery

| Skill | Use it when |
| --- | --- |
| [`kk:test`](../../skills/kk-test/SKILL.md) | Run or design unit, integration, end-to-end, UI, coverage, build, or QA checks. |
| [`kk:web-testing`](../../skills/kk-web-testing/SKILL.md) | Build Playwright, Vitest, k6, accessibility, visual, load, or cross-browser tests. |
| [`kk:code-review`](../../skills/kk-code-review/SKILL.md) | Review pending changes, a commit, a PR, or a codebase for bugs and regressions. |
| [`kk:review-pr`](../../skills/kk-review-pr/SKILL.md) | Review a GitHub PR and optionally fix findings, post the review, or merge after CI. |
| [`kk:git`](../../skills/kk-git/SKILL.md) | Create focused conventional commits, push branches, open PRs, merge, or manage stacked PRs. |

## Browser, UI, and Media

| Skill | Use it when |
| --- | --- |
| [`kk:agent-browser`](../../skills/kk-agent-browser/SKILL.md) | Automate a clean browser, cloud browser, Electron app, form, screenshot, scrape, or exploratory QA flow. |
| [`kk:chrome-profile`](../../skills/kk-chrome-profile/SKILL.md) | Automate a real Chrome profile with the user's cookies, accounts, or authenticated tabs. |
| [`kk:frontend-design`](../../skills/kk-frontend-design/SKILL.md) | Build or reproduce a polished frontend with strong visual fidelity. |
| [`kk:ui-ux-pro-max`](../../skills/kk-ui-ux-pro-max/SKILL.md) | Make UX, design-system, accessibility, typography, layout, interaction, or responsive design decisions. |
| [`kk:threejs`](../../skills/kk-threejs/SKILL.md) | Build a Three.js, WebGL, WebGPU, GLTF, physics, animation, VR, or XR experience. |
| [`kk:ai-multimodal`](../../skills/kk-ai-multimodal/SKILL.md) | Analyze or generate images, audio, video, and documents, including OCR and transcription. |
| [`kk:media-processing`](../../skills/kk-media-processing/SKILL.md) | Transform media with FFmpeg, ImageMagick, background removal, batch jobs, or streaming formats. |

## Documentation and Visualization

| Skill | Use it when |
| --- | --- |
| [`kk:docs`](../../skills/kk-docs/SKILL.md) | Create, update, summarize, or audit project documentation or agent context. |
| [`kk:folder-context`](../../skills/kk-folder-context/SKILL.md) | Add compact agent instructions for one subfolder. |
| [`kk:project-organization`](../../skills/kk-project-organization/SKILL.md) | Choose paths, organize assets, or standardize a project layout. |
| [`kk:preview`](../../skills/kk-preview/SKILL.md) | Preview files or generate visual explanations, slides, diagrams, recaps, or plan reviews. |
| [`kk:mermaidjs-v11`](../../skills/kk-mermaidjs-v11/SKILL.md) | Create editable inline diagrams with Mermaid v11. |
| [`kk:tech-graph`](../../skills/kk-tech-graph/SKILL.md) | Create publish-grade SVG and PNG architecture or flow diagrams. |

## Quick Selection Rules

- Use `kk:plan` when the work needs phases or architecture.
- Use `kk:validate-plan` when an existing plan needs evidence-backed validation before implementation.
- Use `kk:implement` when scope is accepted and code must change.
- Use `kk:fix` when a concrete failure must be repaired.
- Use `kk:debug` when diagnosis is the requested deliverable.
- Use `kk:test` when verification is the requested deliverable.
- Use `kk:code-review` for a local diff and `kk:review-pr` for an existing GitHub PR.
- Use `kk:agent-browser` for clean automation and `kk:chrome-profile` for real logged-in browser state.
- Use `kk:mermaidjs-v11` for editable inline diagrams and `kk:tech-graph` for polished image output.

Next: [Use recipes and avoid pitfalls](./08-recipes-and-pitfalls.md).
