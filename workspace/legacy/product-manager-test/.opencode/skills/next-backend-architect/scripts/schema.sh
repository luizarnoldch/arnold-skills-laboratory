#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SKILL_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
source "$SCRIPT_DIR/lib/helpers.sh"
source "$SCRIPT_DIR/lib/replacer.sh"

scaffold_schema() {
  local target="$1"
  local entity="$2"
  local database="$3"
  
  local entity_kebab
  entity_kebab=$(to_kebab_case "$entity")
  
  local schema_dir="$target/src/features/$entity_kebab/schemas"
  local schema_file="$schema_dir/$entity_kebab.schema.ts"
  
  if [[ -f "$schema_file" ]]; then
    echo "Warning: Schema already exists: $schema_file"
    read -p "Overwrite? [y/N]: " overwrite
    if [[ ! "$overwrite" =~ ^[Yy]$ ]]; then
      echo "Skipped."
      return 0
    fi
  fi
  
  mkdir -p "$schema_dir"
  
  local template_file="$SKILL_DIR/assets/${database}-schema.md"
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
  
  echo "$final_content" > "$schema_file"
  
  echo "✓ Created: $schema_file"
}
