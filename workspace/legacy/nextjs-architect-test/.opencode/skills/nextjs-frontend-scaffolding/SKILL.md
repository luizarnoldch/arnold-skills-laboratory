---
name: nextjs-frontend-scaffolding
description: "Creates and maintains Next.js frontend features: pages, components, views, and hook integrations with server-side data prefetching, Suspense boundaries, and mutations. Use this skill when: building React components (.tsx), creating Next.js pages with app router, implementing TanStack Query hooks, adding Suspense and error boundaries, building forms and mutations, configuring server-side rendering (SSR)."
---
# Next.js Feature Architect

**Frontend feature scaffolding — connects `src/app/` pages to hooks via views and components.**

## Rules

1. **Pages are thin.** `page.tsx` renders a view. No business logic.
2. **Views compose.** `[Entity]View.tsx` wraps ErrorBoundary → Suspense → components. No `use client`.
3. **Hooks stay in entrypoints.** Only `[Entity]List/index.tsx` calls hooks; sub-components receive props.
4. **Types and style.** Use `type` for Props (never `interface`). Arrow functions. `export default`.
5. **Default: suspense.** Use `useSuspenseList[Entity]s`. Fall back to `useList[Entity]s` only if explicitly requested.
6. **Never read `@/components/ui/*` files.** shadcn APIs are stable — import `<Button>`, `<Input>`, `<Dialog>`, `<AlertDialog>`, `<Table>`, `<Badge>`, `<Spinner>` directly without inspection.
7. **Read only `schemas/[entity].schema.ts` and `hooks/*.tsx`** to understand an entity. Skip server files.
8. **Never conditionally call hooks.** For create vs update forms, use separate components (`[Entity]FormCreate`, `[Entity]FormUpdate`).
9. **Never spawn a sub-agent for exploration.** Read these exact files directly, in this order, then stop:
  - schemas/[entity].schema.ts
  - hooks/useSuspenseList[Entity]s.tsx
  - hooks/useCreate[Entity].tsx
  - hooks/useUpdate[Entity].tsx
  - hooks/useDelete[Entity].tsx
  - hooks/Hydrate[Entity]s.tsx
These six files contain everything needed. Do not read server files, ui files, or any other part of the codebase.

**When ambiguous, ask:** "List page or detail page?"

**Dependency:** hooks must exist first — run `next-backend-architect` if missing. **Prefer CLI** (`./scripts/main.sh`).

---

## Output files: list page (`--all`)

Generate exactly these 5 files for `<entity>`:

| # | Path |
|---|------|
| 1 | `src/app/[entity]s/page.tsx` |
| 2 | `src/features/[entity]/views/[Entity]View.tsx` |
| 3 | `src/features/[entity]/components/[Entity]List/index.tsx` |
| 4 | `src/features/[entity]/components/[Entity]FormCreate.tsx` |
| 5 | `src/features/[entity]/components/[Entity]FormUpdate.tsx` |

**Layer flags:**
- `--page list` → file 1 only
- `--view` → file 2 only
- `--view-full` → files 2–5
- `--all` → files 1–5

**Templates:** `assets/templates/` — edit these files to change the generated code pattern.

---

## Architecture

```
src/features/[entity]/
├── views/
│   └── [Entity]View.tsx              ← no 'use client', Tailwind v4 grid layout
└── components/
    ├── [Entity]List/
    │   └── index.tsx                 ← 'use client', data + toggle state + useDelete[Entity]
    ├── [Entity]FormCreate.tsx        ← 'use client', useCreate[Entity]
    └── [Entity]FormUpdate.tsx        ← 'use client', useUpdate[Entity], keyed by id
```

Data flow:
```
page.tsx → Hydrate[Entity]s → [Entity]View → ErrorBoundary → Suspense → [Entity]List
↓ toggles
[Entity]FormCreate
[Entity]FormUpdate
```

UI behavior:
- List shows short data per row (ID + schema fields) to validate data loads
- "Create [Entity]" button shows `[Entity]FormCreate` inline (no dialog)
- "Edit" button per row shows `[Entity]FormUpdate` inline (no dialog)
- "Delete" button per row calls the delete mutation directly

---

## CLI

```bash
./scripts/main.sh <target> <entity> --page list      # File 1 only
./scripts/main.sh <target> <entity> --view           # File 2 only
./scripts/main.sh <target> <entity> --view-full      # Files 2–5
./scripts/main.sh <target> <entity> --all            # Files 1–5
```

Transport is auto-detected from existing hooks. Override: `--transport api`

Detection cues:
- **tRPC:** `useTRPC()`, `trpc.[entity].list.queryOptions()`
- **REST:** `apiFetch`, `apiPrefetch`, raw fetch

---

## Verify

1. `npx tsc --noEmit` — check all imports resolve. If hooks are missing, run `next-backend` agent first.
2. Start dev server → navigate to the page → check browser console for errors.
3. Click "Create [Entity]" → form appears inline. Fill fields, submit, revalidate data and confirm data appears in list.
4. Click "Edit" → update form appears inline with current values. Submit and revalidate data.
5. Click "Delete" → revalidate data.

---

## Notes

- Multi-entity views (e.g. products + categories on one page): `reference/multi-hydration.md`
- Component logic patterns (where mutations live): `reference/component-logic.md`
- REST API requires `src/lib/api.ts` — run `next-backend` agent for setup
