# Use Specialized Workflows

The primary workflows call specialized skills when a task crosses into browsers, visual design, media, documentation, diagrams, or large-context analysis.

Invoke these skills directly when the specialized operation is the complete task.

## Codebase Discovery and External Research

| Skill | Choose it when |
| --- | --- |
| [`kevinle128-skills:scout`](../../skills/scout/SKILL.md) | You need fast file, symbol, owner, caller, or test discovery. |
| [`kevinle128-skills:docs-seeker`](../../skills/docs-seeker/SKILL.md) | You need current framework, package, or API documentation. |
| [`kevinle128-skills:repomix`](../../skills/repomix/SKILL.md) | You need a repository packed into a compact artifact for analysis. |
| [`kevinle128-skills:feature-lens`](../../skills/feature-lens/SKILL.md) | You need to compare, copy, improve, or port one feature from another repository. |
| [`kevinle128-skills:context-engineering`](../../skills/context-engineering/SKILL.md) | Context limits, memory design, token use, or multi-agent context quality are the problem. |
| [`kevinle128-skills:sequential-thinking`](../../skills/sequential-thinking/SKILL.md) | A complex analysis needs explicit steps, revisions, and hypothesis checks. |
| [`kevinle128-skills:problem-solving`](../../skills/problem-solving/SKILL.md) | The current approach is stuck and needs systematic reframing. |

```text
/kevinle128-skills:scout Find the HTTP entry point, state owner, worker callback, and tests for live model updates.
```

```text
/kevinle128-skills:docs-seeker React Server Components cache invalidation
```

Use `kevinle128-skills:feature-lens` when the external repository is not only evidence but the source of a feature you want to adapt.

```text
/kevinle128-skills:feature-lens vercel/ai streaming --port
```

## Browser and Application Automation

Use [`kevinle128-skills:agent-browser`](../../skills/agent-browser/SKILL.md) for clean browser sessions, compact snapshots, screenshots, forms, scraping, exploratory QA, cloud browsers, Electron apps, or long autonomous browser runs.

```text
/kevinle128-skills:agent-browser Open the local app, create an account, change the model setting, and capture the final state.
```

Use [`kevinle128-skills:chrome-profile`](../../skills/chrome-profile/SKILL.md) when the flow requires the user's real Chrome profile, cookies, account, or an already authenticated tab.

```text
/kevinle128-skills:chrome-profile Use my work profile to verify the authenticated settings flow.
```

Use [`kevinle128-skills:web-testing`](../../skills/web-testing/SKILL.md) when the deliverable is a repeatable browser, visual, load, accessibility, or cross-browser test rather than a one-time automated inspection.

## Frontend and UX

Use [`kevinle128-skills:frontend-design`](../../skills/frontend-design/SKILL.md) to build or reproduce a polished interface from requirements, screenshots, or videos.

Use [`kevinle128-skills:ui-ux-pro-max`](../../skills/ui-ux-pro-max/SKILL.md) for design-system choices, typography, color, layout, accessibility, interaction states, responsive behavior, forms, charts, or UX review.

```text
/kevinle128-skills:frontend-design Build the account settings screen from the attached reference and preserve the existing design system.
```

Use [`kevinle128-skills:threejs`](../../skills/threejs/SKILL.md) for WebGL, WebGPU, GLTF, physics, animation, or XR experiences.

Visual work should finish with real browser inspection and the relevant viewport or interaction tests.

## Images, Audio, Video, and Documents

Use [`kevinle128-skills:ai-multimodal`](../../skills/ai-multimodal/SKILL.md) for vision analysis, OCR, transcription, design extraction, and multimodal generation.

Use [`kevinle128-skills:media-processing`](../../skills/media-processing/SKILL.md) for deterministic FFmpeg, ImageMagick, background removal, encoding, conversion, filters, thumbnails, batch processing, or streaming output.

Choose AI multimodal work when semantic interpretation or generation is required.

Choose media processing when a reproducible transformation command can produce the result.

## Documentation and Local Instructions

Use [`kevinle128-skills:docs`](../../skills/docs/SKILL.md) to create, update, summarize, or audit project documentation.

```text
/kevinle128-skills:docs update
```

Use [`kevinle128-skills:folder-context`](../../skills/folder-context/SKILL.md) when one subfolder needs durable local agent instructions beyond the repository root.

Use [`kevinle128-skills:project-organization`](../../skills/project-organization/SKILL.md) to select stable paths and organize files before creating new artifacts.

Documentation should point to executable owners instead of duplicating mutable implementation details.

## Diagrams and Explanations

| Skill | Output |
| --- | --- |
| [`kevinle128-skills:mermaidjs-v11`](../../skills/mermaidjs-v11/SKILL.md) | Inline Mermaid flowcharts, sequences, state diagrams, ER diagrams, timelines, and journeys. |
| [`kevinle128-skills:tech-graph`](../../skills/tech-graph/SKILL.md) | Publish-grade SVG and PNG technical diagrams. |
| [`kevinle128-skills:preview`](../../skills/preview/SKILL.md) | Visual file previews, HTML explanations, slides, plan reviews, recaps, and diagram inspection. |

```text
/kevinle128-skills:mermaidjs-v11 Draw the request-to-worker sequence for a live model update.
```

```text
/kevinle128-skills:tech-graph Create a publish-grade architecture diagram for the session runtime.
```

```text
/kevinle128-skills:preview --html --explain the session runtime architecture
```

Use Mermaid for documentation that should remain editable as text.

Use `kevinle128-skills:tech-graph` for a polished image artifact.

Use `kevinle128-skills:preview` to explain or visually review an existing artifact.

Next: [Browse the complete skill catalog](./07-skill-catalog.md).
