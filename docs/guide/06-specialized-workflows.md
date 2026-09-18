# Use Specialized Workflows

The primary workflows call specialized skills when a task crosses into browsers, visual design, media, documentation, diagrams, or large-context analysis.

Invoke these skills directly when the specialized operation is the complete task.

## Codebase Discovery and External Research

| Skill | Choose it when |
| --- | --- |
| [`kk:scout`](../../skills/kk-scout/SKILL.md) | You need fast file, symbol, owner, caller, or test discovery. |
| [`kk:docs-seeker`](../../skills/kk-docs-seeker/SKILL.md) | You need current framework, package, or API documentation. |
| [`kk:repomix`](../../skills/kk-repomix/SKILL.md) | You need a repository packed into a compact artifact for analysis. |
| [`kk:feature-lens`](../../skills/kk-feature-lens/SKILL.md) | You need to compare, copy, improve, or port one feature from another repository. |
| [`kk:context-engineering`](../../skills/kk-context-engineering/SKILL.md) | Context limits, memory design, token use, or multi-agent context quality are the problem. |
| [`kk:sequential-thinking`](../../skills/kk-sequential-thinking/SKILL.md) | A complex analysis needs explicit steps, revisions, and hypothesis checks. |
| [`kk:problem-solving`](../../skills/kk-problem-solving/SKILL.md) | The current approach is stuck and needs systematic reframing. |

```text
/kk:scout Find the HTTP entry point, state owner, worker callback, and tests for live model updates.
```

```text
/kk:docs-seeker React Server Components cache invalidation
```

Use `kk:feature-lens` when the external repository is not only evidence but the source of a feature you want to adapt.

```text
/kk:feature-lens vercel/ai streaming --port
```

## Browser and Application Automation

Use [`kk:agent-browser`](../../skills/kk-agent-browser/SKILL.md) for clean browser sessions, compact snapshots, screenshots, forms, scraping, exploratory QA, cloud browsers, Electron apps, or long autonomous browser runs.

```text
/kk:agent-browser Open the local app, create an account, change the model setting, and capture the final state.
```

Use [`kk:chrome-profile`](../../skills/kk-chrome-profile/SKILL.md) when the flow requires the user's real Chrome profile, cookies, account, or an already authenticated tab.

```text
/kk:chrome-profile Use my work profile to verify the authenticated settings flow.
```

Use [`kk:web-testing`](../../skills/kk-web-testing/SKILL.md) when the deliverable is a repeatable browser, visual, load, accessibility, or cross-browser test rather than a one-time automated inspection.

## Frontend and UX

Use [`kk:frontend-design`](../../skills/kk-frontend-design/SKILL.md) to build or reproduce a polished interface from requirements, screenshots, or videos.

Use [`kk:ui-ux-pro-max`](../../skills/kk-ui-ux-pro-max/SKILL.md) for design-system choices, typography, color, layout, accessibility, interaction states, responsive behavior, forms, charts, or UX review.

```text
/kk:frontend-design Build the account settings screen from the attached reference and preserve the existing design system.
```

Use [`kk:threejs`](../../skills/kk-threejs/SKILL.md) for WebGL, WebGPU, GLTF, physics, animation, or XR experiences.

Visual work should finish with real browser inspection and the relevant viewport or interaction tests.

## Images, Audio, Video, and Documents

Use [`kk:ai-multimodal`](../../skills/kk-ai-multimodal/SKILL.md) for vision analysis, OCR, transcription, design extraction, and multimodal generation.

Use [`kk:media-processing`](../../skills/kk-media-processing/SKILL.md) for deterministic FFmpeg, ImageMagick, background removal, encoding, conversion, filters, thumbnails, batch processing, or streaming output.

Choose AI multimodal work when semantic interpretation or generation is required.

Choose media processing when a reproducible transformation command can produce the result.

## Documentation and Local Instructions

Use [`kk:docs`](../../skills/kk-docs/SKILL.md) to create, update, summarize, or audit project documentation.

```text
/kk:docs update
```

Use [`kk:folder-context`](../../skills/kk-folder-context/SKILL.md) when one subfolder needs durable local agent instructions beyond the repository root.

Use [`kk:project-organization`](../../skills/kk-project-organization/SKILL.md) to select stable paths and organize files before creating new artifacts.

Documentation should point to executable owners instead of duplicating mutable implementation details.

## Diagrams and Explanations

| Skill | Output |
| --- | --- |
| [`kk:mermaidjs-v11`](../../skills/kk-mermaidjs-v11/SKILL.md) | Inline Mermaid flowcharts, sequences, state diagrams, ER diagrams, timelines, and journeys. |
| [`kk:tech-graph`](../../skills/kk-tech-graph/SKILL.md) | Publish-grade SVG and PNG technical diagrams. |
| [`kk:preview`](../../skills/kk-preview/SKILL.md) | Visual file previews, HTML explanations, slides, plan reviews, recaps, and diagram inspection. |

```text
/kk:mermaidjs-v11 Draw the request-to-worker sequence for a live model update.
```

```text
/kk:tech-graph Create a publish-grade architecture diagram for the session runtime.
```

```text
/kk:preview --html --explain the session runtime architecture
```

Use Mermaid for documentation that should remain editable as text.

Use `kk:tech-graph` for a polished image artifact.

Use `kk:preview` to explain or visually review an existing artifact.

Next: [Browse the complete skill catalog](./07-skill-catalog.md).
