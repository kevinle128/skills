# Visual Explanation Routing

Use this file when a workflow asks for a visual explanation, diagram, slide deck,
diff review, or recap. Load `../SKILL.md` first for command syntax, then use this
file to choose the mode.

## Mode Selection

| Need | Preview mode |
|---|---|
| View an existing Markdown file or directory | `/kevinle128-skills:preview <path>` |
| Explain a concept or code path | `/kevinle128-skills:preview --explain <topic>` |
| Generate a focused architecture/data-flow diagram | `/kevinle128-skills:preview --diagram <topic>` |
| Terminal-friendly diagram only | `/kevinle128-skills:preview --ascii <topic>` |
| Self-contained HTML explanation | `/kevinle128-skills:preview --html --explain <topic>` |
| Slide deck | `/kevinle128-skills:preview --html --slides <topic>` |
| Visual diff review for a branch, PR, or commit | `/kevinle128-skills:preview --html --diff [ref]` |
| Compare an implementation plan to code | `/kevinle128-skills:preview --html --plan-review <plan>` |
| Recap recent project context | `/kevinle128-skills:preview --html --recap [timeframe]` |

## Specialist Handoffs

- Mermaid syntax: load `/kevinle128-skills:mermaidjs-v11`.
- Publish-grade SVG/PNG architecture diagrams: use `/kevinle128-skills:tech-graph`.
- Generated images or multimodal analysis: use `/kevinle128-skills:ai-multimodal`.
- UI/UX style selection for slides or high-polish HTML: use
  `/kevinle128-skills:ui-ux-pro-max`.
- Documentation update after a durable visual: use `/kevinle128-skills:docs update` and
  `../../docs/references/documentation-management.md`.

## Output Rules

- Prefer the active plan's `visuals/` folder when a plan exists.
- If no plan exists, save under `plans/visuals/`.
- For HTML output, always include the theme toggle required by
  `html-css-patterns.md`.
- For diagrams, render and inspect the output; syntax validity alone is not
  enough.
