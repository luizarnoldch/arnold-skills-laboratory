# OKF v0.2 — Attested Computation

An **Attested Computation** concept carries a sanctioned way to compute a
value so a consumer can confirm the agent ran the blessed computation
instead of improvising. Provenance answers “where did this claim come
from”; attestation answers “was this number produced the way we said.”

OKF records the contract; it does **not** execute anything.

## Standalone concept

Use `type: Attested Computation` as its own file. Metrics, tables, or
narratives **link** to it with normal markdown links.

Why standalone:

- `runtime` defines what `parameters` mean (SQL binds vs dbt vars vs Python).
- One computation can back many consumers.
- Trust / staleness / attester are per computation.

## Contract fields

In addition to provenance/trust/lifecycle families:

| Field         | Role |
|---------------|------|
| `runtime`     | **Required for this type.** How to run (e.g. `bigquery`, `postgres`, `dbt`, `python`). |
| `parameters`  | Named typed holes the agent may fill: `{ name, type, required }`. |
| `computation` | Optional path to a computation file; absent ⇒ body `# Computation` fence. |
| `executor`    | `resource` → run instructions/code; `receipt` → fields a run must return. |
| `attester`    | `resource` → deterministic (no-LLM) checker over a receipt. |

Packaging behind `resource` (Skill, script, container) is out of scope; OKF
fixes the interface only.

### Example

```markdown
---
type: Attested Computation
title: Revenue for fiscal year
description: Recognized revenue for a fiscal year, per Finance's definition.
status: stable
runtime: bigquery
parameters:
  - { name: year, type: integer, required: true }
executor:
  resource: references/skills/run-on-bq.md
  receipt: [job_id, executed_sql, result]
attester:
  resource: references/attesters/revenue.py
generated: { by: reference_agent/gemini-2.5-pro, at: 2026-06-20T22:53:05Z }
verified: { by: human:ahormati, at: 2026-06-25T09:00:00Z }
stale_after: 2026-09-23
sources:
  - id: rev-policy
    resource: https://wiki.acme/finance/revenue-recognition
    title: Revenue recognition policy
---

# Computation

    SELECT SUM(amount) AS revenue
    FROM finance.recognized_revenue
    WHERE fiscal_year = @year

Binds only declared parameters per policy.[^rev-policy]

[^rev-policy]: Revenue recognition policy
```

Or point at a file instead of an inline fence:

```yaml
runtime: bigquery
computation: references/computations/lib/revenue.sql
parameters:
  - { name: year, type: integer, required: true }
```

The agent MAY supply **values** for declared parameters only; it MUST NOT
author or edit the computation text. Binding and attestation comparison are
the consumer’s job (compare expanded artifact in the receipt to the
sanctioned computation).

## Concepts that use a computation

Keep narrative metrics readable; link out per figure:

```markdown
---
type: Metric
title: Revenue
description: Recognized revenue for a fiscal year.
tags: [finance, revenue]
status: stable
generated: { by: reference_agent/gemini-2.5-pro, at: 2026-06-20T22:53:05Z }
---

# Definition

Computed by [the revenue computation](/computations/revenue.md).
```

Co-locate computations under e.g. `computations/` with an `index.md` — a
directory choice, not a frontmatter one.

## Consumer flow (informative)

Runtime artifacts (receipts, verdicts) are **not** stored in the bundle.

1. Discover via `type: Attested Computation` or a link from a metric.
2. Load frontmatter contract + computation body/file.
3. Parameterize with declared params only.
4. Execute via `executor`; obtain receipt shaped by `executor.receipt`.
5. Run `attester` over the receipt; surface failures (do not silently drop).
6. Warn or refuse when `today >= stale_after`.

## Verification vs attestation

| Signal      | Confirms                         | Where        | Cadence   |
|-------------|----------------------------------|--------------|-----------|
| `verified`  | Definition still matches policy  | In bundle    | Slow/doc  |
| Attestation | This run used the sanctioned path| Runtime only | Per call  |

Both are needed: a stale definition can still attest; a fresh definition
still needs per-run attestation.

## Producer guidance

- Prefer Attested Computation when agents will compute numbers users must
  trust (finance metrics, SLAs).
- For ordinary practices/runbooks, skip this type — use plain concepts.
- Never invent `executor` / `attester` paths that do not exist; leave
  `# Open questions` or omit attestation fields until real.
