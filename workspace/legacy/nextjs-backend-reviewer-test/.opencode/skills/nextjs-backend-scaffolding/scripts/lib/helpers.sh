#!/usr/bin/env bash
set -euo pipefail

show_tree() {
  echo "Available scaffold combinations:"
  echo ""
  echo "  --schema --database prisma    → Zod schema from Prisma types"
  echo "  --schema --database drizzle   → Zod schema from Drizzle types"
  echo ""
  echo "  --server                      → tRPC router + service + repository (defaults: trpc, prisma)"
  echo "  --server --transport api      → REST API client wrapper"
  echo "  --server --transport trpc --database drizzle → tRPC with Drizzle ORM"
  echo ""
  echo "  --hooks                       → TanStack Query hooks (queries + mutations)"
  echo "  --hooks --transport api       → TanStack Query hooks (REST API)"
  echo ""
  echo "  --all --database prisma       → Full stack (schema + server + hooks)"
  echo "  --all --transport api --database drizzle → Full stack with REST API + Drizzle"
}

validate_path() {
  local target="$1"
  if [[ ! -d "$target" ]]; then
    echo "Error: Target path does not exist: $target" >&2
    return 1
  fi
  if [[ ! -d "$target/src" ]]; then
    echo "Error: Target does not appear to be a Next.js project (no src/ folder)" >&2
    return 1
  fi
  return 0
}

validate_entity() {
  local entity="$1"
  if [[ -z "$entity" ]]; then
    echo "Error: Entity name is required" >&2
    return 1
  fi
  if [[ ! "$entity" =~ ^[A-Z][a-zA-Z0-9]*$ ]] && [[ ! "$entity" =~ ^[a-z][a-z0-9-]*$ ]]; then
    echo "Error: Entity must be PascalCase (Product) or kebab-case (product)" >&2
    return 1
  fi
  return 0
}

prompt_database() {
  echo "Select database:"
  echo "  1) prisma"
  echo "  2) drizzle"
  read -p "Choice [1-2]: " choice
  case "$choice" in
    1) echo "prisma" ;;
    2) echo "drizzle" ;;
    *) echo "Invalid choice" >&2; return 1 ;;
  esac
}

prompt_entity() {
  read -p "Enter entity name (e.g., Product, user-profile): " entity
  echo "$entity"
}

to_pascal_case() {
  local input="$1"
  if [[ "$input" =~ ^[A-Z] ]]; then
    echo "$input"
  else
    echo "$input" | sed -E 's/(^|[-_])([a-z])/\U\2/g; s/[-_]//g'
  fi
}

to_kebab_case() {
  local input="$1"
  if [[ "$input" =~ ^[a-z] ]]; then
    echo "$input"
  else
    echo "$input" | sed -E 's/([A-Z])/-\L\1/g; s/^-//; s/([a-z])([A-Z])/\1-\L\2/g'
  fi
}

to_snake_case() {
  local input="$1"
  echo "$input" | sed -E 's/([A-Z])/_\L\1/g; s/^_//; s/([a-z])([A-Z])/\1_\L\2/g; s/-/_/g' | tr '[:upper:]' '[:lower:]'
}

to_camel_case() {
  local input="$1"
  local pascal
  pascal=$(to_pascal_case "$input")
  echo "$pascal" | sed -E 's/^([A-Z])/\L\1/'
}

prompt_transport() {
  echo "Select transport:"
  echo "  1) trpc"
  echo "  2) api"
  read -p "Choice [1-2]: " choice
  case "$choice" in
    1) echo "trpc" ;;
    2) echo "api" ;;
    *) echo "Invalid choice" >&2; return 1 ;;
  esac
}

validate_schema_exists() {
  local target="$1"
  local entity="$2"
  
  local entity_kebab
  entity_kebab=$(to_kebab_case "$entity")
  local schema_file="$target/src/features/$entity_kebab/schemas/$entity_kebab.schema.ts"
  
  if [[ ! -f "$schema_file" ]]; then
    echo "false"
  else
    echo "true"
  fi
}

validate_api_client() {
  local target="$1"
  local api_client="$target/src/lib/api.ts"
  
  if [[ ! -f "$api_client" ]]; then
    echo "Warning: API client not found at $api_client" >&2
    echo "         See references/external-api.md for setup instructions" >&2
    return 1
  fi
  return 0
}

validate_server_exists() {
  local target="$1"
  local entity="$2"
  local transport="$3"
  
  local entity_kebab
  entity_kebab=$(to_kebab_case "$entity")
  
  if [[ "$transport" == "trpc" ]]; then
    local router_file="$target/src/features/$entity_kebab/server/${entity_kebab}.router.ts"
    if [[ -f "$router_file" ]]; then
      echo "true"
    else
      echo "false"
    fi
  else
    local api_file="$target/src/features/$entity_kebab/server/${entity_kebab}.api.ts"
    if [[ -f "$api_file" ]]; then
      echo "true"
    else
      echo "false"
    fi
  fi
}
