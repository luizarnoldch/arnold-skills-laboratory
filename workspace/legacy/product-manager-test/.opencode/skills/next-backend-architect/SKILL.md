---
name: next-backend-architect
description: "Scaffolds backend feature layers for a Next.js app using a CLI + template system. Use this skill whenever the user wants to create, generate, scaffold, or add backend code for any feature or entity — including schemas, tRPC routers, REST API endpoints, services, repositories, or TanStack Query/Form hooks. Triggers on: 'create feature', 'scaffold', 'add CRUD', 'new entity', 'generate backend', 'tRPC router for X', 'hooks for X', 'repository for X', 'service for X', or any request to build backend layers for a named model or table. NEVER write these files from scratch using your general knowledge of tRPC, Prisma, or Drizzle — the project has a specific CLI and template system that must be followed. Improvising produces code that breaks the project's architecture even when it looks correct."
---

# Next.js Backend Architect

Generates backend feature layers via CLI + templates. The CLI is the source of truth — never write these files from your own knowledge of tRPC or Prisma.

## Decision flow

When a request arrives, follow this sequence exactly:

1. **Identify the entity name** — extract it from the user's request (e.g. "Product", "Order").
2. **Identify the ORM** — look for a Prisma schema (`schema.prisma`) or Drizzle file in the project. Default: Prisma.
3. **Identify the transport** — tRPC unless the user explicitly says "REST" or "API route". Default: tRPC.
4. **Identify which layers** — schema, server, hooks, or all. If ambiguous, ask: *"Which layers? `schema`, `server`, `hooks`, or `all`?"* Do not assume `all`.
5. **Run the CLI** (see below). Do not create files manually before attempting the CLI.
6. **If CLI fails**, fall back to templates (see Fallback section).

## Layer rules

Each layer is independent. Only create what was requested:

- `schema` → `schemas/`; after generation, update fields to match the entity model using Zod v4 types.
- `server` → `server/` files + auto-creates schema if missing (CLI handles this)
- `hooks` → `hooks/` files + auto-creates server + schema if missing (CLI handles this)
- `all` → full stack

The dependency chain (`hooks → server → schema`) is handled by the CLI automatically. Do not create dependency layers manually. In fallback mode, create in order: schema → server → hooks.

## Target directory

All output goes under `src/features/[entity]/`:

```
src/features/[entity]/
├── schemas/
│   └── [entity].schema.ts
├── server/
│   ├── [entity].router.ts      ← tRPC transport
│   ├── [entity].api.ts         ← REST transport
│   ├── [entity].service.ts
│   └── [entity].repository.ts
└── hooks/
    ├── Hydrate[Entity]s.tsx
    ├── useList[Entity]s.tsx
    ├── useSuspenseList[Entity]s.tsx
    ├── useCreate[Entity].tsx
    ├── useUpdate[Entity].tsx
    └── useDelete[Entity].tsx
```

## Running the CLI

Use `./scripts/main.sh` from the project root. Compose flags based on what was requested:

```bash
# Schema only
./scripts/main.sh <target> <Entity> --schema --database <prisma|drizzle>

# Server only (tRPC)
./scripts/main.sh <target> <Entity> --server --transport trpc --database <prisma|drizzle>

# Server only (REST)
./scripts/main.sh <target> <Entity> --server --transport api --database <prisma|drizzle>

# Hooks only
./scripts/main.sh <target> <Entity> --hooks --transport <trpc|api>

# Full stack
./scripts/main.sh <target> <Entity> --schema --server --hooks --transport <trpc|api> --database <prisma|drizzle>
```

`<target>` is the destination path (e.g. `src/features`). `<Entity>` is PascalCase.

## Fallback — when the CLI is unavailable

The CLI is unavailable when `./scripts/main.sh` does not exist or exits with a non-zero code unrelated to arguments.

**Do not improvise. Do not write files from your own knowledge of tRPC or Prisma.**

Template files live in `assets/`. Match your layer and ORM/transport to the correct file:

- Schema + Prisma → `assets/prisma-schema.md`
- Schema + Drizzle → `assets/drizzle-schema.md`
- Server tRPC router → `assets/trpc-router.md`
- Server tRPC service → `assets/trpc-service.md`
- Server tRPC Prisma repository → `assets/prisma-repository.md`
- Server tRPC Drizzle repository → `assets/drizzle-repository.md`
- Server REST → `assets/api-server.md`
- Hooks tRPC hydrate → `hydrate-hook.md`
- Hooks tRPC list → `assets/list-hook.md`
- Hooks tRPC list (suspense) → `assets/list-suspense-hook.md`
- Hooks tRPC create → `assets/create-hook.md`
- Hooks tRPC update → `assets/update-hook.md`
- Hooks tRPC delete → `assets/delete-hook.md`
- Hooks REST hydrate → `assets/hydrate-api-hook.md`
- Hooks REST list → `assets/list-api-hook.md`
- Hooks REST list (suspense) → `assets/list-suspense-api-hook.md`
- Hooks REST create → `assets/create-api-hook.md`
- Hooks REST update → `assets/update-api-hook.md`
- Hooks REST delete → `assets/delete-api-hook.md`

To apply a template: read the file, copy its content exactly, replace `[Entity]` with PascalCase, `[entity]` with camelCase, `[entityTable]` with snake_case, `[entity-kebab]` with kebab-case, then write to the correct path under `src/features/[entity]/`. Do not add logic or imports beyond what the template contains.

## REST transport requirement

If generating REST transport (`--transport api`), the project requires `src/lib/api.ts`. Check if it exists first. If missing, read `references/external-api.md` and create it before generating the server layer.
