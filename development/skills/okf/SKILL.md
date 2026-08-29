---
name: okf
description: >
  Creates, maintains, migrates, and navigates Open Knowledge Format (OKF) v0.2
  knowledge bundles — directories of markdown concepts with YAML frontmatter.
  Use when the user asks about OKF, Open Knowledge Format, knowledge bundles,
  okf/, LLM-wiki patterns, concept graphs for agents, Attested Computation,
  progressive disclosure via index.md, or migrating OKF v0.1 to v0.2. Does NOT
  apply to CLAUDE.md / AGENTS.md / DESIGN.md behavior files, SKILL.md authoring,
  inventing textbook definitions the model already knows, app source-code edits,
  git/CI/package work, or treating OKF as an SEO ranking signal.
---
<skill_content name='okf'>

# Open Knowledge Format (OKF)

Produce and consume **OKF v0.2** bundles: a directory of markdown concept files
with YAML frontmatter. No SDK, no runtime, no platform — if you can `cat` a
file, you can read OKF.

OKF describes **what the project/domain knows** (tables, metrics, practices,
runbooks, decisions). It does **not** replace `CLAUDE.md` / `AGENTS.md`
(agent behavior) or `SKILL.md` (procedural methods).

## Core rules (v0.2)

1. **One concept per file.** Concept ID = path within the bundle without `.md`
   (e.g. `practicas/cuando-usar-rag.md` → `practicas/cuando-usar-rag`).
2. **YAML frontmatter required** on every non-reserved `.md` file. Only
   **`type`** is mandatory (non-empty). Unknown types and extra keys are fine.
3. **Reserved filenames:** `index.md` (directory map) and `log.md` (changelog)
   are never concepts.
4. **Links create untyped relationships.** Prefer absolute bundle-relative
   links (`/tables/orders.md`). Broken links are allowed (not-yet-written
   knowledge).
5. **Document real team knowledge**, not textbook definitions the model already
   knows. If unsure, add `# Open questions` instead of inventing facts.

Recommended frontmatter: `title`, `description`, `resource`, `tags`, plus
`generated` (replaces v0.1 `timestamp`). Optional trust/lifecycle families:
`sources`, `verified`, `status`, `stale_after` — see references.

Root `index.md` MAY declare:

```yaml
---
okf_version: "0.2"
---
```

(only place `index.md` may have frontmatter).

## Workflows

### 1. Seed a bundle (init)

Create a small graph (typically under `okf/`) with **four related concepts**,
plus root `index.md` and `log.md`:

1. A central practice / convention.
2. A decision governed by that practice.
3. A runbook that prevents a failure the practice cares about.
4. A runbook for a tool or process that consumes the same budget/policy.

Rules:

- Every concept file: parseable YAML frontmatter with non-empty `type`.
- Prefer types like `Convention`, `Decision`, `Runbook`, `Playbook`,
  `Metric`, `BigQuery Table` — descriptive, not registered centrally.
- Set `generated: { by: <actor>, at: <ISO-8601> }` (use `human:<id>` or
  `agent/<version>` per actor convention).
- Cross-link with absolute `/…` markdown links.
- Root `index.md`: sections listing concepts with short descriptions (from
  frontmatter). Include `okf_version: "0.2"` when seeding fresh.
- Root `log.md`: ISO date headings (`YYYY-MM-DD`), newest first; initial
  **Creation** / **Initialization** entry.
- Do **not** dump the whole repo into concepts. Start small.

### 2. Expand the graph

Add **only** concepts that are already linked from existing documents and
still missing. Each new file must fill a hole another document asked for.
Update `index.md` and append to `log.md`.

### 3. Consume a bundle

1. Read root (or relevant directory) `index.md` first — progressive disclosure.
2. Open **only** the concept files needed for the question.
3. Honor advisory signals: `status` (`draft` | `stable` | `deprecated`;
   absent ⇒ stable), `stale_after` (stale when `today >= stale_after`), and
   trust tiers from `verified` (none ⇒ unverified; non-human only ⇒
   machine-confirmed; any `human:` ⇒ human-reviewed).
4. Do not load the entire tree into context.

### 4. Migrate v0.1 → v0.2

- `timestamp` → `generated: { by: …, at: <former timestamp> }` (consumers may
  still fall back to legacy `timestamp`).
- Body `# Citations` list → frontmatter `sources` with stable `id`s; attribute
  claims via footnotes `[^id]` keyed to `sources[].id`.
- Optionally add `status`, `stale_after`, `verified` where known.
- Set root `okf_version: "0.2"` when migrating the bundle.

Deep field rules: `references/provenance-trust.md`. Attested Computation
concepts: `references/attested-computation.md`. Full structural conformance:
`references/spec-core.md`.

## Out of scope

- Authoring or rewriting `CLAUDE.md`, `AGENTS.md`, `DESIGN.md`, or `SKILL.md`.
- Inventing textbook definitions, fake schemas, or unsourced metrics.
- Treating OKF as SEO / SERP ranking.
- Unrelated app code, git, CI/CD, or package management (unless the user
  explicitly asks to document those as OKF concepts).

## Quick concept skeleton

```markdown
---
type: Convention
title: Context window policy
description: How we manage token budget in our agents.
tags: [context, agents]
generated: { by: human:team, at: 2026-07-30T12:00:00Z }
status: stable
---

# Rule

Keep the window under 70% before compacting.

# Related concepts

See [when we use RAG](/practicas/cuando-usar-rag.md).

# Open questions

- Exact compaction trigger in CI agents?
```

<skill_resources>
  <file>references/spec-core.md</file>
  <file>references/provenance-trust.md</file>
  <file>references/attested-computation.md</file>
</skill_resources>
</skill_content>
