# OKF v0.2 — Core structure

Condensed producer/consumer rules. Prefer this over inventing conventions.

## Bundle

A **bundle** is a directory tree of UTF-8 markdown files. Distribute as git
repo (recommended), tarball/zip, or a subdirectory of a larger repo.

```
path/to/bundle/
  index.md                 # Optional. Progressive-disclosure listing.
  log.md                   # Optional. Chronological update history.
  <concept>.md
  <subdirectory>/
    index.md
    <concept>.md
    ...
```

Organization is domain-defined. Tags live in frontmatter (`tags`), not a
separate tag file format.

## Reserved filenames

| Filename   | Role                                      |
|------------|-------------------------------------------|
| `index.md` | Directory listing — **not** a concept     |
| `log.md`   | Update history — **not** a concept        |

All other `.md` files are concepts.

## Concept documents

1. YAML frontmatter between `---` lines at the start of the file.
2. Markdown body after the closing `---`.

### Frontmatter

**Required**

- `type` — short non-empty string (e.g. `BigQuery Table`, `Metric`,
  `Playbook`, `Convention`, `Attested Computation`). Not centrally
  registered. Consumers MUST tolerate unknown types.

**Recommended**

- `title` — display name (else derive from filename)
- `description` — one-line summary
- `resource` — canonical URI of the underlying asset (omit for pure ideas)
- `tags` — YAML list of short strings

**Optional families** (see `provenance-trust.md`, `attested-computation.md`)

- Provenance / trust / lifecycle: `sources`, `generated`, `verified`,
  `status`, `stale_after`, `usage_window`
- Attested Computation: `runtime`, `parameters`, `computation`, `executor`,
  `attester`

**Extensions:** any extra keys allowed. Consumers SHOULD preserve unknown
keys on round-trip and MUST NOT reject them.

### Body

No required sections. Prefer structured markdown (headings, tables, lists,
fenced code). Conventional headings when applicable:

| Heading         | Purpose                                      |
|-----------------|----------------------------------------------|
| `# Schema`      | Columns / fields of an asset                 |
| `# Examples`    | Concrete usage                               |
| `# Computation` | Sanctioned computation (Attested Computation)|

Per-claim attribution uses footnotes keyed to `sources[].id`, not a body
`# Citations` list (v0.1 legacy — migrate when touching the file).

## Cross-linking

- **Absolute (recommended):** `/tables/orders.md` — relative to bundle root;
  stable if the file moves within its folder.
- **Relative:** `./other.md` — also valid.

A link asserts an **untyped** relationship; the surrounding prose names the
kind (joins-with, depends-on, etc.).

Consumers MUST tolerate broken links.

Path-valued fields (`resource`, `sources[].resource`, `computation`,
`executor.resource`, `attester.resource`) accept absolute URLs,
bundle-absolute `/…` paths, or relative paths. A `sources[].resource` MAY
instead be a non-path scope descriptor (e.g. “all queries in project X”).

Convention: `references/` may mirror external material, executors, or
attesters as first-class concepts — naming only, not required.

## Index files (`index.md`)

Support progressive disclosure: show what exists before opening documents.

- No frontmatter, **except** bundle-root `index.md` MAY set
  `okf_version: "0.2"`.
- Body: sections with bullet links and short descriptions.

```markdown
---
okf_version: "0.2"
---

# Practices

* [Context window policy](/practicas/politica-de-contexto.md) - Token budget rules.
```

Producers MAY auto-generate indexes; consumers MAY synthesize one if missing.

## Log files (`log.md`)

Flat, date-grouped entries, **newest first**. Date headings MUST be
`YYYY-MM-DD`. Leading bold words (`**Update**`, `**Creation**`,
`**Deprecation**`) are conventional, not required.

```markdown
# Directory Update Log

## 2026-07-30
* **Creation**: Seeded context-policy practice graph.
```

## Conformance (v0.2)

A bundle is conformant if:

1. Every non-reserved `.md` has parseable YAML frontmatter.
2. Every frontmatter has non-empty `type`.
3. Present `index.md` / `log.md` follow the shapes above.

Consumers MUST NOT reject bundles for missing optional fields, unknown
types, unknown keys, broken links, or missing indexes.

When trust/lifecycle/provenance/computation families are present, follow
those reference docs. Treat a bare `verified` mapping as a one-element list.

## Versioning

- Declare target with `okf_version` on root `index.md`.
- Minor bumps add optional fields; major bumps may break.
- Unknown declared versions: best-effort consume, do not refuse.

### v0.1 → v0.2 consumer fallbacks

- Prefer `generated.at`; fall back to legacy `timestamp` if absent.
- Prefer `sources`; MAY still parse legacy body `# Citations`.
