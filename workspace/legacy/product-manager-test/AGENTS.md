# OpenCode Agents & Skills — Usage Guide

This project uses a team of OpenCode [agents](#agents) and [skills](#skills) to scaffold, review, and manage Next.js full-stack features. You interact with them through natural-language chat or the `opencode` CLI.

---

## Quick-start Flow

Every feature follows this pipeline:

```
Idea ──► PRD ──► Backend ──► Frontend ──► Review
         │        │             │             │
         ▼        ▼             ▼             ▼
  product-manager  nextjs-architect     nextjs-architect    next-backend
   + prd-writer     -orchestrator      -orchestrator      -reviewer
                    (calls backend      (calls frontend
                     architect)          architect)
```

1. Start with an idea → get a PRD.
2. Confirm the PRD → the orchestrator scaffolds **backend first**.
3. Backend passes a gate check → the orchestrator scaffolds the **frontend**.
4. A reviewer validates everything → you see **PASS** or a fix list.

---

## What Lives in This Project

### Agents

All agent definitions are in `.opencode/agents/`.

| Agent | What it does | When it triggers |
|---|---|---|
| `nextjs-architect-orchestrator` | Routes your request to the right subagent. Enforces **backend-before-frontend**. Read-only + `write: deny` — never writes files itself, only delegates via Task. | Auto-triggered on: "create feature for X", "scaffold X", "build X feature", "CRUD for X", "full stack for X" |
| `next-backend-architect` | Thin Task-invocable wrapper around the `next-backend-architect` skill. Has `write: allow`. Scaffolds schema/router/service/repository/hooks via the project's CLI + templates. | Invoked by the orchestrator via Task — never call the skill directly, it has no permissions of its own. |
| `next-feature-architect` | Thin Task-invocable wrapper around the `next-feature-architect` skill. Has `write: allow`. Scaffolds pages/views/components wired to existing backend hooks. | Invoked by the orchestrator via Task, after the backend gate passes. |
| `next-backend-reviewer` | Runs the `next-backend-reviewer` skill's `scripts/validate.sh` (a deterministic bash script, not LLM judgment) against generated files, optionally runs `tsc --noEmit`, and reports **PASS/FAIL** with fix instructions. Read-only — `write: deny`. | Auto-triggered by the orchestrator, as its own Task call, after every scaffold run. |
| `product-manager` | Turn an idea into a structured PRD (`specs/*.md`) and a `FEATURES.yml` registry entry. | Auto-triggered on: "I want a new feature", "plan X", "idea for X" |

### Skills

All skill definitions are in `.opencode/skills/`.

| Skill | What it does | When it triggers |
|---|---|---|
| `prd-writer` | Provides the PRD template, EARS acceptance criteria format, and writing rules. | Loaded by `product-manager` automatically; never invoked directly. |
| `next-backend-architect` | Scaffolds backend layers (schema, router/service/repository, hooks) from project-specific templates in `assets/`. | Loaded by the `next-backend-architect` **agent** (same name) — that agent is what the orchestrator Tasks; the skill itself has no write permission and is never Task-invoked directly. |
| `next-feature-architect` | Scaffolds frontend layers (page, view, list component, create/update forms) with Suspense + SSR. | Loaded by the `next-feature-architect` **agent** (same name), after backend hooks exist. |
| `next-backend-reviewer` | Documents the target file layout and report format, and points to `scripts/validate.sh` — a bash script that does the actual per-layer PASS/FAIL checking (no LLM judgment involved). | Loaded by the `next-backend-reviewer` **agent** (same name), which runs the script and formats its output. |

> Every skill above (except `prd-writer`) has a same-named agent in `.opencode/agents/` that wraps it with the runtime permissions (`write`, `bash`) the skill's work requires, and is what actually gets Task-invoked. Skills themselves carry no permissions and are not directly Task-addressable.

---

## Workflow Gate Matrix

The orchestrator decides what to build based on what you ask.

| You say | Backend layers | Frontend layers | Review |
|---|---|---|---|
| "just the schema" | Schema only | — | ✅ |
| "backend only" / "server" / "hooks" | Schema + Router + Service + Repository + Hooks | — | ✅ |
| "full stack" / "everything" / "CRUD for X" | All | All | ✅ |
| "list page" / "table" / "form" | — (verifies hooks first) | Page + View + Components | ✅ |
| "create Product" (ambiguous) | Orchestrator asks: "Backend, frontend, or both?" | — | — |

> **Backend always comes before frontend.** The orchestrator will refuse to scaffold a page until `useCreateProduct.tsx` and the hook files exist.

---

## Step-by-Step Example: Building a "Product" Feature

This walkthrough shows both **natural-language chat** and **CLI** ways to drive the pipeline.

### Step 1 — Idea → PRD

**Chat:**
```
I need a product catalog with name, price, and description.
```

**CLI:**
```bash
opencode run -a product-manager \
  "Create a PRD for a product catalog with name, price, and description"
```

**What happens:**
- `product-manager` asks up to 5 clarifying questions (e.g., "Should products have categories?").
- It loads the `prd-writer` skill for the template format.
- It writes `specs/product-catalog.md` and appends an entry to `FEATURES.yml`.

**Output:**
```
specs/
└── product-catalog.md
FEATURES.yml          ← new registry entry
```

---

### Step 2 — PRD → Backend Scaffolding

**Chat:**
```
Scaffold the Product backend — Prisma, tRPC.
```

**CLI:**
```bash
opencode run -a nextjs-architect-orchestrator \
  "Full stack Product feature. Database: prisma. Transport: trpc. Target: src/features."
```

**What happens:**
1. Orchestrator parses the entity (`Product`), transport (`trpc`), ORM (`prisma`).
2. It invokes `next-backend-architect` via Task:
   ```
   Task: "Generate backend layers for Product. Layers: all. Transport: trpc. Database: prisma. Target: src/features."
   ```
3. `next-backend-architect` uses the project's CLI + template files in `assets/` to generate:

**Output:**
```
src/features/product/
├── schemas/
│   └── product.schema.ts          ← Zod schema, types
├── server/
│   ├── product.router.ts          ← tRPC router (list, get, create, update, delete)
│   ├── product.service.ts         ← business logic, error messages
│   └── product.repository.ts      ← Prisma delegation
└── hooks/
    ├── HydrateProducts.tsx
    ├── useListProducts.tsx
    ├── useSuspenseListProducts.tsx
    ├── useCreateProduct.tsx
    ├── useUpdateProduct.tsx
    └── useDeleteProduct.tsx
```

---

### Step 3 — Gate Check (automatic)

Before frontend scaffolding starts, the orchestrator verifies:

- `src/features/product/hooks/useSuspenseListProducts.tsx` exists?
- `src/features/product/hooks/useCreateProduct.tsx` exists?
- `src/features/product/hooks/useUpdateProduct.tsx` exists?
- `src/features/product/hooks/useDeleteProduct.tsx` exists?

If any are missing, the orchestrator **re-invokes** `next-backend-architect` with `layers: hooks` before continuing.

---

### Step 4 — Backend → Frontend Scaffolding

**Chat:**
```
Now build the Product list page.
```

**CLI:**
```bash
opencode run -a nextjs-architect-orchestrator \
  "Generate frontend layers for Product. Flag: --all. Hooks already at src/features/product/."
```

**What happens:**
1. Orchestrator invokes `next-feature-architect` via Task:
   ```
   Task: "Generate frontend layers for Product using flag --all. Schema and hooks are already at src/features/product/."
   ```
2. `next-feature-architect` reads the existing hooks and schema, then generates React components with Suspense, SSR, and TanStack Query integration.

**Output:**
```
src/app/products/page.tsx
src/features/product/views/ProductView.tsx
src/features/product/components/ProductList/index.tsx
src/features/product/components/ProductFormCreate.tsx
src/features/product/components/ProductFormUpdate.tsx
```

---

### Step 5 — Review

This step runs **automatically** after scaffolding finishes. You don't need to ask for it.

**What happens:**
1. Orchestrator invokes `next-backend-reviewer` via Task:
   ```
   Task: "Review backend layers for Product. Scope: --backend."
   ```
2. The reviewer loads its skill, runs `scripts/validate.sh` (a deterministic bash script — the check itself is not an LLM judgment call), and optionally `tsc --noEmit`.

**If PASS:**
```markdown
## Backend Review: Product

### Verdict: PASS

All backend layers conform to the `next-backend-architect` templates. Type check clean. Nothing to fix.
```

**If FAIL:**
```markdown
## Backend Review: Product

### Verdict: FAIL

### Required Fixes
- **File:** `src/features/product/server/product.router.ts`
  - **Issue:** `delete` route is missing.
  - **How to fix:** Add `delete: protectedProcedure.input(deleteProductSchema)...` with `try/catch` + `TRPCError` per the tRPC router template.
```

The orchestrator then **loops back** to `next-backend-architect` with the fixes and re-runs the reviewer until it passes.

---

## Invoking from the CLI

You can run any agent directly with `opencode run`.

### Run the orchestrator (full stack)
```bash
opencode run -a nextjs-architect-orchestrator \
  "Scaffold a full-stack Order feature with tRPC and Prisma"
```

### Run a single subagent directly
```bash
# Scaffold only the backend
opencode run -a nextjs-architect-orchestrator \
  "Generate backend layers for Order. Layers: all. Transport: trpc. Database: prisma."

# Review only
opencode run -a next-backend-reviewer \
  "Review backend layers for Product. Scope: --backend."
```

### Run a skill directly (not recommended)
Skills are designed to be **loaded by agents**, not invoked standalone. If you do invoke one, you must provide the full context the agent would normally supply:
```bash
# Not recommended — lacks orchestrator context
opencode run -s next-backend-architect "Generate hooks for Product"
```

> **Rule of thumb:** Use the **orchestrator** for full-stack work. Use **subagents directly** only when you know exactly which layer you need.

---

## Per-Agent Deep Dive

### `nextjs-architect-orchestrator`
- **Never writes files.** It only delegates via the `Task` tool.
- **Enforces order:** Backend → Gate → Frontend → Review.
- **Auto-detects:** ORM (Prisma vs Drizzle) by looking for `schema.prisma` or `drizzle.config.*`. Transport defaults to `trpc` unless you say "REST" or "API route".

### `next-backend-reviewer`
- **Read-only.** It cannot create, edit, or delete files (`write: deny`).
- **Script-driven, not LLM-judged.** Runs `scripts/validate.sh`, which greps generated files against the actual `next-backend-architect/assets/` templates. The PASS/FAIL verdict is deterministic — the agent formats the script's output, it doesn't re-derive compliance itself.
- **Runs type checks.** Passes `--typecheck` to the script when `tsconfig.json` exists, which runs `npx tsc --noEmit`.
- **REST-aware.** Knows `service.ts`/`repository.ts` don't exist under REST transport (REST's `<entity>.api.ts` is a client wrapper calling an external API, not a hosted route handler) — it doesn't check for files that were never meant to exist.
- **Reports structure:**
  - Layer results table (schema, router, service, repository, hooks)
  - Required fixes (file, issue, how to fix)
  - Type check result

### `product-manager`
- **Never writes application code.** It only writes PRDs and `FEATURES.yml`.
- **Always asks first.** Up to 5 clarifying questions before drafting.
- **Requires confirmation.** After writing a PRD, it asks: "Should I scaffold this now?" before invoking the orchestrator.

---

## Per-Skill Deep Dive

### `next-backend-architect`
- **Uses a project-specific CLI.** Located at `scripts/main.sh`.
- **Uses template files.** All output is generated from `assets/` (e.g., `trpc-router.md`, `prisma-schema.md`, `trpc-service.md`).
- **Never improvise.** Even if the AI "knows" how a tRPC router should look, it must follow the templates.

### `next-feature-architect`
- **Works with existing hooks.** It reads `src/features/[entity]/hooks/` to understand mutation signatures.
- **Suspense-first.** All list views use `useSuspenseQuery`, hydrated via `HydrateClient` + `prefetch()` (tRPC) or `apiPrefetch()` (REST) — not `HydrationBoundary`/`dehydrate`.
- **Never reverses order.** It will error if hooks don't exist yet.

### `next-backend-reviewer`
- **Owns `scripts/validate.sh`.** All layer-compliance logic lives in `scripts/checks/*.sh`, not in prose — SKILL.md documents the file layout and report format, and points to the script as the single source of truth.
- **Never duplicate checks in prose.** If a check needs to change, edit the script. Keeping a parallel written checklist is what caused the reviewer to drift out of sync with the real templates before (e.g. checking for `prisma`/`@/lib/prisma` when every template actually imports `db` from `@/lib/db`).

### `prd-writer`
- **Template provider.** Loaded automatically by `product-manager`.
- **Defines:** LLM-PRD v1 format, EARS acceptance criteria syntax, naming conventions, Implementation Triggers.

---

## Troubleshooting

### "Backend reviewer says FAIL — missing delete route"
The orchestrator will automatically re-invoke `next-backend-architect` with the fix instructions and re-run the reviewer. You don't need to do anything unless you want to fix it manually.

### "Frontend won't generate because hooks are missing"
The orchestrator's gate caught a missing backend layer. It will re-invoke the backend architect to regenerate hooks before allowing frontend scaffolding.

### "Where is FEATURES.yml?"
In the project root. If it doesn't exist, `product-manager` creates it during its first run.

### "How do I switch from tRPC to REST?"
Tell the orchestrator explicitly:
```bash
opencode run -a nextjs-architect-orchestrator \
  "Scaffold Product backend with REST API transport"
```
The orchestrator forwards `Transport: api` to `next-backend-architect`, which then uses `assets/api-server.md` (a client wrapper calling an external API via `apiFetch` — not a hosted route handler) and the `assets/*-api-hook.md` templates. Note REST transport has no `service.ts`/`repository.ts` — only `<entity>.api.ts`.

### "Reviewer says 'mixed transport detected'"
Some hook files import `useTRPC` while the router is REST (or vice versa). The orchestrator will re-invoke the backend architect with the correct transport to regenerate all layers consistently.

---

## File Paths Reference

| File | Description |
|---|---|
| `.opencode/agents/nextjs-architect-orchestrator.md` | Orchestrator agent definition (`write: deny` — delegates everything via Task) |
| `.opencode/agents/next-backend-architect.md` | Backend scaffold subagent definition (`write: allow`, wraps the skill below) |
| `.opencode/agents/next-feature-architect.md` | Frontend scaffold subagent definition (`write: allow`, wraps the skill below) |
| `.opencode/agents/next-backend-reviewer.md` | Reviewer subagent definition (`write: deny`, runs the skill's `scripts/validate.sh`) |
| `.opencode/agents/product-manager.md` | Product manager agent definition |
| `.opencode/skills/next-backend-architect/SKILL.md` | Backend scaffold skill (decision flow, CLI usage, template fallback) |
| `.opencode/skills/next-backend-architect/scripts/main.sh` | Backend architect CLI entry point |
| `.opencode/skills/next-backend-architect/assets/` | Template files used by the backend architect |
| `.opencode/skills/next-feature-architect/SKILL.md` | Frontend scaffold skill |
| `.opencode/skills/next-feature-architect/scripts/main.sh` | Frontend architect CLI entry point |
| `.opencode/skills/next-backend-reviewer/SKILL.md` | Backend review skill (file layout + report format; points to the script below) |
| `.opencode/skills/next-backend-reviewer/scripts/validate.sh` | Deterministic per-layer validation script — the actual PASS/FAIL logic |
| `.opencode/skills/prd-writer/SKILL.md` | PRD writing skill |
| `FEATURES.yml` | Feature registry managed by product-manager (created on first PRD if missing) |
| `specs/` | PRD output directory |
