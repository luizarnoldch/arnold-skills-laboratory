# OKF v0.2 — Provenance, trust, and lifecycle

All fields in this document are **optional**. Absence is meaningful (e.g.
unverified vs verified) but never a reason to reject a concept.

## Actor convention

Used in `generated.by` and `verified[].by`:

| Form                 | Meaning                          | Example                         |
|----------------------|----------------------------------|---------------------------------|
| `<producer>/<version>` | Agent or tool                  | `reference_agent/gemini-2.5-pro`|
| `human:<id>`         | Person                           | `human:ahormati`                |
| `process:<id>`       | Automated non-LLM process        | `process:finance-nightly`       |

Producers MUST use the `human:` prefix for hand-authored or human-confirmed
content so trust tiers work.

## Provenance: `sources`

Records materials a concept derives from (external or in-bundle).

```yaml
sources:
  - id: ga4-schema
    resource: https://developers.google.com/analytics/bigquery/export-schema
    title: GA4 BigQuery Export schema
    author: team:ga4-docs
    usage_count: 5000
    last_modified: 2026-05-30
usage_window: { from: 2026-06-01, to: 2026-06-30 }
```

Per entry:

- `resource` — **required** within an entry. URL, bundle path, path under
  `references/`, or a scope descriptor that is not a path.
- `id` — optional stable key for claim attribution; SHOULD be set when the
  body cites the source.
- `title` — optional label.
- Credibility signals (optional, objective facts — OKF does **not** store a
  credibility score):
  - `author` — who produced the source (actor-like string).
  - `usage_count` — how often `resource` was exercised over `usage_window`
    (liveness/trend, not a precise cross-kind ranking).
  - `last_modified` — when the **source** last changed (`YYYY-MM-DD`);
    distinct from `generated.at` (when the concept was written).
- `usage_window` — usually a sibling of `sources` (`{ from, to }`); a single
  entry MAY override with its own window.

Lineage to other OKF concepts: point `resource` at that concept and use
normal links; recurse into the target’s `sources` if needed. Explicit
external `derived_from` graphs are out of scope for v0.2.

### Per-claim attribution

Use markdown footnotes whose **label** is `sources[].id`:

```markdown
The `events_` table is sharded daily as `events_YYYYMMDD`.[^ga4-schema]

[^ga4-schema]: GA4 BigQuery Export schema
```

Consumers join on the label → `sources` entry, not by parsing footnote prose.
Stable `id`s survive reordering; positional indexes do not.

**v0.1:** body `# Citations` lists are superseded; migrate to `sources` when
editing.

## Trust: `generated` and `verified`

Kept distinct: who **wrote** need not be who **confirmed**.

```yaml
generated: { by: reference_agent/gemini-2.5-pro, at: 2026-06-20T22:53:05Z }
verified:
  - { by: human:ahormati, at: 2026-06-25T09:00:00Z }
  - { by: process:finance-nightly, at: 2026-06-26T02:00:00Z }
```

- `generated.by` — required when `generated` is present (actor).
- `generated.at` — ISO 8601 last meaningful content change.
- `verified` — list of `{ by, at }` events. A bare mapping MUST be treated
  as a one-element list:

```yaml
verified: { by: human:ahormati, at: 2026-06-25T09:00:00Z }
```

`verified` is independent of `generated.at` (content can change without
re-confirmation; facts can be re-confirmed without regeneration).

**v0.1:** `timestamp` is superseded by `generated.at`. Consumers MAY fall
back to `timestamp` when `generated` is absent.

### Trust tiers (derived, advisory)

| Condition                         | Tier               |
|-----------------------------------|--------------------|
| No `verified`                     | unverified         |
| Only non-`human:` verifiers       | machine-confirmed  |
| Any `human:<id>` verifier         | human-reviewed     |

Not access control. Consumers MUST NOT reject unverified concepts.

## Lifecycle: `status`

```yaml
status: stable   # draft | stable | deprecated
```

- `draft` — incomplete / not reviewed.
- `stable` — ready (also the default when absent).
- `deprecated` — kept for links/history; no longer current.

## Lifecycle: `stale_after`

```yaml
stale_after: 2026-09-23   # YYYY-MM-DD; stale when today >= this date
```

Absolute date only (no relative TTL). Staleness is a plain date comparison.
