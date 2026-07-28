#!/usr/bin/env bash
set -euo pipefail

# check_schema <target> <entity_kebab> <entity_pascal> <entity_camel> <database> [prisma_file]
check_schema() {
  local target="$1" entity_kebab="$2" entity_pascal="$3" entity_camel="$4" database="$5"
  local prisma_file="${6:-}"
  local file="$target/src/features/$entity_kebab/schemas/$entity_kebab.schema.ts"
  local reasons=()

  require_file "$file" "schema file missing" || { fail_layer schema "${reasons[@]}"; return 1; }

  require_grep "$file" 'from "zod"' "must import z from zod" || true
  require_grep "$file" "delete${entity_pascal}Schema[[:space:]]*=" "must export delete${entity_pascal}Schema" || true
  require_grep "$file" "create${entity_pascal}Schema[[:space:]]*=" "must export create${entity_pascal}Schema" || true
  require_grep "$file" "update${entity_pascal}Schema[[:space:]]*=" "must export update${entity_pascal}Schema" || true

  if [[ "$database" == "drizzle" ]]; then
    require_grep "$file" "createSelectSchema\(" "drizzle schema must derive from createSelectSchema" || true
    require_grep "$file" "createInsertSchema\(" "drizzle schema must derive from createInsertSchema" || true
    require_grep "$file" "createUpdateSchema\(" "drizzle schema must derive from createUpdateSchema" || true
  else
    require_grep "$file" "${entity_camel}Schema[[:space:]]*=[[:space:]]*z\.object" "must export ${entity_camel}Schema as z.object(...)" || true
    require_grep "$file" '\.omit\(' "create${entity_pascal}Schema must use .omit(...)" || true
    require_grep "$file" '\.partial\(\)' "update${entity_pascal}Schema must use .partial()" || true
  fi

  require_grep "$file" "type ${entity_pascal}[[:space:]]*=[[:space:]]*z\.infer" "must export type ${entity_pascal} via z.infer" || true
  require_grep "$file" "type Create${entity_pascal}Input" "must export type Create${entity_pascal}Input" || true
  require_grep "$file" "type Update${entity_pascal}Input" "must export type Update${entity_pascal}Input" || true

  if grep -Eq '\binterface\b' "$file"; then
    reasons+=("must not use the interface keyword")
  fi

  if grep -qE '^[[:space:]]*\.\.\.\s*,?\s*$' "$file"; then
    reasons+=("schema contains the '...' template placeholder -- custom model fields were not filled in (replace '...' with actual entity fields derived from the Prisma model)")
  fi

  # --- Prisma model cross-check ---
  if [[ -n "$prisma_file" && -f "$prisma_file" && "$database" == "prisma" ]]; then
    local prisma_fields
    prisma_fields=$(get_prisma_all_fields "$prisma_file" "$entity_pascal" 2>/dev/null) || true
    if [[ -n "$prisma_fields" ]]; then
      # Extract Zod field names from the z.object({...}) block
      local zod_fields
      zod_fields=$(sed -n '/z\.object({/,/})/p' "$file" 2>/dev/null \
        | sed -n '/^[[:space:]]*[a-zA-Z_][a-zA-Z0-9_]*[[:space:]]*:/p' \
        | sed 's/^[[:space:]]*//; s/[:].*//' || true)

      # Check missing: Prisma fields not in Zod schema
      local missing_fields=()
      while IFS= read -r pf; do
        [[ -z "$pf" ]] && continue
        if ! echo "$zod_fields" | grep -qFx "$pf"; then
          missing_fields+=("$pf")
        fi
      done <<< "$prisma_fields"

      if [[ ${#missing_fields[@]} -gt 0 ]]; then
        local missing_str
        missing_str=$(IFS=,; echo "${missing_fields[*]}")
        reasons+=("Zod schema missing Prisma model fields: $missing_str")
      fi

      # Check extra: Zod fields not in Prisma model
      local extra_fields=()
      while IFS= read -r zf; do
        [[ -z "$zf" ]] && continue
        if ! echo "$prisma_fields" | grep -qFx "$zf"; then
          extra_fields+=("$zf")
        fi
      done <<< "$zod_fields"

      if [[ ${#extra_fields[@]} -gt 0 ]]; then
        local extra_str
        extra_str=$(IFS=,; echo "${extra_fields[*]}")
        reasons+=("Zod schema has fields not in Prisma model: $extra_str")
      fi
    fi
  fi
  # --- end Prisma cross-check ---

  local field_names base_fields
  base_fields='^(id|createdAt|updatedAt|deletedAt)$'
  field_names=$(grep -oE '^[[:space:]]+[a-zA-Z_]+\s*:' "$file" 2>/dev/null | sed 's/^[[:space:]]*//;s/://;s/[[:space:]]*$//' || true)
  local has_custom=false
  while IFS= read -r fname; do
    if [[ -n "$fname" && ! "$fname" =~ $base_fields ]]; then
      has_custom=true
      break
    fi
  done <<< "$field_names"

  if [[ "$has_custom" == false ]]; then
    reasons+=("schema contains only base fields (id, createdAt, updatedAt) -- entity-specific custom fields are missing (the Prisma model fields were not translated into Zod schema fields)")
  fi

  if [[ ${#reasons[@]} -eq 0 ]]; then
    pass_layer schema
  else
    fail_layer schema "${reasons[@]}"
  fi
}
