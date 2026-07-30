#!/usr/bin/env bash
set -euo pipefail

HELPERS_LIB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

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

to_camel_case() {
  local input="$1"
  local pascal
  pascal=$(to_pascal_case "$input")
  echo "$pascal" | sed -E 's/^([A-Z])/\L\1/'
}

to_snake_case() {
  local input="$1"
  echo "$input" | sed -E 's/([A-Z])/_\L\1/g; s/^_//; s/([a-z])([A-Z])/\1_\L\2/g; s/-/_/g' | tr '[:upper:]' '[:lower:]'
}

prompt_entity() {
  read -p "Enter entity name (e.g., Product, user-profile): " entity
  echo "$entity"
}

prompt_page_type() {
  echo "Select page type:"
  echo "  1) list"
  echo "  2) detail"
  read -p "Choice [1-2]: " choice
  case "$choice" in
    1) echo "list" ;;
    2) echo "detail" ;;
    *) echo "Invalid choice" >&2; return 1 ;;
  esac
}

detect_transport() {
  local target="$1"
  local entity_kebab="$2"

  local hooks_dir="$target/src/features/$entity_kebab/hooks"

  if [[ ! -d "$hooks_dir" ]]; then
    echo "trpc"
    return
  fi

  local hydrate_file="$hooks_dir/Hydrate$(to_pascal_case "${entity_kebab//-/_}")s.tsx"
  if [[ -f "$hydrate_file" ]]; then
    if grep -q "apiPrefetch\|apiFetch\|list[Entity]s" "$hydrate_file" 2>/dev/null; then
      echo "api"
      return
    fi
  fi

  local list_file="$hooks_dir/useSuspenseList$(to_pascal_case "${entity_kebab//-/_}")s.tsx"
  if [[ -f "$list_file" ]]; then
    if grep -q "apiFetch\|from \"..\/server" "$list_file" 2>/dev/null; then
      echo "api"
      return
    fi
  fi

  if [[ -f "$list_file" ]]; then
    if grep -q "useTRPC\|trpc\." "$list_file" 2>/dev/null; then
      echo "trpc"
      return
    fi
  fi

  echo "trpc"
}

show_tree() {
  echo "Available scaffold combinations:"
  echo ""
  echo "  --page list                  → List page + view (src/app/[entity]s/page.tsx + views/)"
  echo "  --page detail                → Detail page + view (src/app/[entity]s/[id]/page.tsx + views/)"
  echo "  --view                       → View only (no page)"
  echo "  --all                        → Full feature frontend (page + view + components)"
  echo ""
  echo "  --transport trpc             → Use tRPC transport"
  echo "  --transport api              → Use REST API transport"
  echo ""
  echo "  Transport is auto-detected from existing hooks if not specified."
}

validate_hooks_exist() {
  local target="$1"
  local entity="$2"

  local entity_kebab
  entity_kebab=$(to_kebab_case "$entity")
  local hooks_dir="$target/src/features/$entity_kebab/hooks"

  if [[ ! -d "$hooks_dir" ]]; then
    echo "false"
    return
  fi

  local list_file="$hooks_dir/useSuspenseList$(to_pascal_case "$entity")s.tsx"
  if [[ -f "$list_file" ]]; then
    echo "true"
    return
  fi

  echo "false"
}
