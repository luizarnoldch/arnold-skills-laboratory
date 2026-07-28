#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SKILL_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
source "$SCRIPT_DIR/lib/helpers.sh"
source "$SCRIPT_DIR/lib/replacer.sh"

HOOK_FILES=(
  "hydrate:hydrate-hook.md:Hydrate[Entity]s.tsx"
  "list:list-hook.md:useList[Entity]s.tsx"
  "list-suspense:list-suspense-hook.md:useSuspenseList[Entity]s.tsx"
  "create:create-hook.md:useCreate[Entity].tsx"
  "update:update-hook.md:useUpdate[Entity].tsx"
  "delete:delete-hook.md:useDelete[Entity].tsx"
)

HOOK_API_FILES=(
  "hydrate:hydrate-api-hook.md:Hydrate[Entity]s.tsx"
  "list:list-api-hook.md:useList[Entity]s.tsx"
  "list-suspense:list-suspense-api-hook.md:useSuspenseList[Entity]s.tsx"
  "create:create-api-hook.md:useCreate[Entity].tsx"
  "update:update-api-hook.md:useUpdate[Entity].tsx"
  "delete:delete-api-hook.md:useDelete[Entity].tsx"
)

generate_hook() {
  local template_file="$1"
  local output_file="$2"
  local entity="$3"
  
  if [[ ! -f "$template_file" ]]; then
    echo "Error: Template not found: $template_file" >&2
    return 1
  fi
  
  local template_content
  template_content=$(cat "$template_file")
  
  local code_content
  code_content=$(extract_code_block "$template_content")
  
  local final_content
  final_content=$(replace_placeholders "$code_content" "$entity")
  
  echo "$final_content" > "$output_file"
}

scaffold_hooks() {
  local target="$1"
  local entity="$2"
  local transport="$3"
  local database="${4:-prisma}"
  
  local entity_kebab
  entity_kebab=$(to_kebab_case "$entity")
  
  local schema_exists
  schema_exists=$(validate_schema_exists "$target" "$entity")
  
  if [[ "$schema_exists" == "false" ]]; then
    echo "→ Schema not found, creating..."
    source "$SCRIPT_DIR/schema.sh"
    scaffold_schema "$target" "$entity" "$database"
  fi
  
  local server_exists
  server_exists=$(validate_server_exists "$target" "$entity" "$transport")
  
  if [[ "$server_exists" == "false" ]]; then
    echo "→ Server layer not found, creating..."
    source "$SCRIPT_DIR/server.sh"
    scaffold_server "$target" "$entity" "$transport" "$database"
  fi
  
  local hooks_dir="$target/src/features/$entity_kebab/hooks"
  mkdir -p "$hooks_dir"
  
  local files
  if [[ "$transport" == "trpc" ]]; then
    files=("${HOOK_FILES[@]}")
  else
    files=("${HOOK_API_FILES[@]}")
  fi
  
  for entry in "${files[@]}"; do
    IFS=':' read -r name template output_name <<< "$entry"
    
    local output_file="$hooks_dir/${output_name//\[Entity\]/$(to_pascal_case "$entity")}"
    output_file="${output_file//\[entity\]/$(to_camel_case "$entity")}"
    
    local template_file="$SKILL_DIR/assets/$template"
    
    generate_hook "$template_file" "$output_file" "$entity"
    echo "✓ Created: $output_file"
  done
}
