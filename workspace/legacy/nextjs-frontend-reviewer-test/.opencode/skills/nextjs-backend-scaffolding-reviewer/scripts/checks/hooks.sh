#!/usr/bin/env bash
set -euo pipefail

_check_hook_default_values() {
  local hook_file="$1"
  local expected_fields_str="$2"
  local hook_label="$3"

  if ! grep -q 'defaultValues:' "$hook_file" 2>/dev/null; then
    reasons+=("${hook_label}: missing defaultValues object")
    return 1
  fi

  if grep -Pzoq '\{[^}]*\[[a-zA-Z]+\]\s*:' "$hook_file" 2>/dev/null; then
    reasons+=("${hook_label}: defaultValues contains template placeholder keys (e.g. [requiredFields]) -- replace with actual entity fields")
    return 1
  fi

  if grep -q '\[defaultValue\]\|\[primitiveInitualValue\]\|\[primitiveInitialValue\]' "$hook_file" 2>/dev/null; then
    reasons+=("${hook_label}: defaultValues contains template placeholder values (e.g. [defaultValue]) -- replace with actual field init values")
    return 1
  fi

  if [[ -z "$expected_fields_str" ]]; then
    return 0
  fi

  local hook_fields
  hook_fields=$(sed -n '/defaultValues:[[:space:]]*{/,/}[[:space:]]*as[[:space:]].*Input/p' "$hook_file" 2>/dev/null \
    | sed -n 's/^[[:space:]]*\([a-zA-Z_][a-zA-Z0-9_]*\)[[:space:]]*:.*/\1/p' || true)

  [[ -z "$hook_fields" ]] && return 0

  local missing_fields=()
  while IFS= read -r expected; do
    [[ -z "$expected" ]] && continue
    if ! echo "$hook_fields" | grep -qFx "$expected"; then
      missing_fields+=("$expected")
    fi
  done <<< "$expected_fields_str"

  if [[ ${#missing_fields[@]} -gt 0 ]]; then
    local missing_str
    missing_str=$(IFS=,; echo "${missing_fields[*]}")
    reasons+=("${hook_label}: defaultValues missing expected fields: $missing_str")
  fi
}

# check_hooks <target> <entity_kebab> <entity_pascal> <entity_camel> <transport> [prisma_file]
check_hooks() {
  local target="$1" entity_kebab="$2" entity_pascal="$3" entity_camel="$4" transport="$5"
  local prisma_file="${6:-}"
  local dir="$target/src/features/$entity_kebab/hooks"
  local reasons=()

  local hydrate="$dir/Hydrate${entity_pascal}s.tsx"
  local use_list="$dir/useList${entity_pascal}s.tsx"
  local use_suspense="$dir/useSuspenseList${entity_pascal}s.tsx"
  local use_create="$dir/useCreate${entity_pascal}.tsx"
  local use_update="$dir/useUpdate${entity_pascal}.tsx"
  local use_delete="$dir/useDelete${entity_pascal}.tsx"

  for f in "$hydrate" "$use_list" "$use_suspense" "$use_create" "$use_update" "$use_delete"; do
    require_file "$f" "$(basename "$f") missing" || true
  done
  if [[ ${#reasons[@]} -gt 0 ]]; then
    fail_layer hooks "${reasons[@]}"
    return 1
  fi

  for f in "$use_suspense" "$use_create" "$use_update" "$use_delete"; do
    require_grep "$f" '^"use client"' "$(basename "$f") must start with \"use client\"" || true
  done

  if [[ "$transport" == "trpc" ]]; then
    require_grep "$hydrate" 'HydrateClient' "$(basename "$hydrate") must import HydrateClient from @/trpc/server" || true
    require_grep "$hydrate" "prefetch\(trpc\.${entity_camel}\.list" "$(basename "$hydrate") must call prefetch(trpc.${entity_camel}.list.queryOptions())" || true

    require_grep "$use_suspense" 'useTRPC.*from "@/trpc/client"' "$(basename "$use_suspense") must import useTRPC from @/trpc/client" || true
    require_grep "$use_suspense" "useSuspenseQuery\(trpc\.${entity_camel}\.list" "$(basename "$use_suspense") must call useSuspenseQuery(trpc.${entity_camel}.list.queryOptions())" || true

    require_grep "$use_list" "useQuery\(trpc\.${entity_camel}\.list" "$(basename "$use_list") must call useQuery(trpc.${entity_camel}.list.queryOptions())" || true

    require_grep "$use_create" 'useForm' "$(basename "$use_create") must use useForm from @tanstack/react-form" || true
    require_grep "$use_create" "create${entity_pascal}Schema" "$(basename "$use_create") must wire create${entity_pascal}Schema into the form validators" || true
    require_grep "$use_create" 'toast' "$(basename "$use_create") must import toast from sonner" || true

    require_grep "$use_update" "trpc\.${entity_camel}\.update\.mutationOptions" "$(basename "$use_update") must call trpc.${entity_camel}.update.mutationOptions()" || true

    require_grep "$use_delete" "trpc\.${entity_camel}\.delete\.mutationOptions" "$(basename "$use_delete") must call trpc.${entity_camel}.delete.mutationOptions()" || true

    for f in "$hydrate" "$use_list" "$use_suspense" "$use_create" "$use_update" "$use_delete"; do
      if grep -Eq 'apiFetch|apiPrefetch' "$f"; then
        reasons+=("$(basename "$f") mixes REST transport (apiFetch/apiPrefetch) into a tRPC feature")
      fi
    done

  elif [[ "$transport" == "api" ]]; then
    require_grep "$hydrate" 'apiPrefetch' "$(basename "$hydrate") must call apiPrefetch({ queryKey, queryFn })" || true
    require_grep "$hydrate" "list${entity_pascal}s.*from \"\.\./server/$entity_kebab\.api\"" "$(basename "$hydrate") must import list${entity_pascal}s from ../server/$entity_kebab.api" || true

    require_grep "$use_suspense" "list${entity_pascal}s" "$(basename "$use_suspense") must call list${entity_pascal}s() from the REST api client" || true
    require_grep "$use_list" "list${entity_pascal}s" "$(basename "$use_list") must call list${entity_pascal}s() from the REST api client" || true
    require_grep "$use_create" "create${entity_pascal}" "$(basename "$use_create") must call create${entity_pascal}(...) from the REST api client" || true
    require_grep "$use_update" "update${entity_pascal}" "$(basename "$use_update") must call update${entity_pascal}(...) from the REST api client" || true
    require_grep "$use_delete" "delete${entity_pascal}" "$(basename "$use_delete") must call delete${entity_pascal}(...) from the REST api client" || true

    for f in "$hydrate" "$use_list" "$use_suspense" "$use_create" "$use_update" "$use_delete"; do
      if grep -Eq 'useTRPC\(\)' "$f"; then
        reasons+=("$(basename "$f") mixes tRPC transport (useTRPC) into a REST feature")
      fi
    done

  else
    reasons+=("unknown transport '$transport' (expected trpc or api)")
  fi

  # --- Schema.ts field cross-check for useCreate / useUpdate (both transports) ---
  local schema_file="$target/src/features/$entity_kebab/schemas/$entity_kebab.schema.ts"
  if [[ -f "$schema_file" ]]; then
    local schema_create_fields schema_update_fields
    schema_create_fields=$(get_schema_create_fields "$schema_file" "$entity_camel" 2>/dev/null) || true
    schema_update_fields=$(get_schema_update_fields "$schema_file" "$entity_camel" 2>/dev/null) || true

    if [[ -n "$schema_create_fields" ]]; then
      _check_hook_default_values "$use_create" "$schema_create_fields" "useCreate${entity_pascal}"
    fi

    if [[ -n "$schema_update_fields" ]]; then
      _check_hook_default_values "$use_update" "$schema_update_fields" "useUpdate${entity_pascal}"
    fi
  fi

  # --- Additional Prisma model field cross-check when available (both transports) ---
  if [[ -n "$prisma_file" && -f "$prisma_file" ]]; then
    local prisma_create_fields prisma_update_fields
    prisma_create_fields=$(get_prisma_create_fields "$prisma_file" "$entity_pascal" 2>/dev/null) || true
    prisma_update_fields=$(get_prisma_update_fields "$prisma_file" "$entity_pascal" 2>/dev/null) || true

    if [[ -n "$prisma_create_fields" ]]; then
      _check_hook_default_values "$use_create" "$prisma_create_fields" "useCreate${entity_pascal} (Prisma cross-check)"
    fi

    if [[ -n "$prisma_update_fields" ]]; then
      _check_hook_default_values "$use_update" "$prisma_update_fields" "useUpdate${entity_pascal} (Prisma cross-check)"
    fi
  fi

  if [[ ${#reasons[@]} -eq 0 ]]; then
    pass_layer hooks
  else
    fail_layer hooks "${reasons[@]}"
  fi
}
