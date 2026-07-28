#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SKILL_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
source "$SCRIPT_DIR/lib/helpers.sh"
source "$SCRIPT_DIR/lib/replacer.sh"

generate_from_template() {
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

generate_api_file() {
  local entity="$1"
  local entity_pascal
  local entity_camel
  local entity_kebab

  entity_pascal=$(to_pascal_case "$entity")
  entity_camel=$(to_camel_case "$entity")
  entity_kebab=$(to_kebab_case "$entity")

  cat <<EOF
import { apiFetch } from "@/lib/api"
import type { Create${entity_pascal}Schema, Update${entity_pascal}Schema } from "../schemas/${entity_kebab}.schema"

export type ${entity_pascal}Response = {
  id: string
  name: string
  createdAt: string
  updatedAt: string
}

export async function list${entity_pascal}s(): Promise<${entity_pascal}Response[]> {
  return apiFetch<${entity_pascal}Response[]>(\`/${entity_kebab}s\`, { method: "GET" })
}

export async function get${entity_pascal}(id: string): Promise<${entity_pascal}Response> {
  return apiFetch<${entity_pascal}Response>(\`/${entity_kebab}s/\${id}\`, { method: "GET" })
}

export async function create${entity_pascal}(data: Create${entity_pascal}Schema): Promise<${entity_pascal}Response> {
  return apiFetch<${entity_pascal}Response>(\`/${entity_kebab}s\`, {
    method: "POST",
    body: JSON.stringify(data),
  })
}

export async function update${entity_pascal}(id: string, data: Update${entity_pascal}Schema): Promise<${entity_pascal}Response> {
  return apiFetch<${entity_pascal}Response>(\`/${entity_kebab}s/\${id}\`, {
    method: "PUT",
    body: JSON.stringify(data),
  })
}

export async function delete${entity_pascal}(id: string): Promise<void> {
  return apiFetch<void>(\`/${entity_kebab}s/\${id}\`, { method: "DELETE" })
}
EOF
}

scaffold_server() {
  local target="$1"
  local entity="$2"
  local transport="$3"
  local database="$4"

  local entity_kebab
  entity_kebab=$(to_kebab_case "$entity")

  local schema_exists
  schema_exists=$(validate_schema_exists "$target" "$entity")

  if [[ "$schema_exists" == "false" ]]; then
    echo "→ Schema not found, creating..."
    source "$SCRIPT_DIR/schema.sh"
    scaffold_schema "$target" "$entity" "$database"
  fi

  if [[ "$transport" == "trpc" ]]; then
    local server_dir="$target/src/features/$entity_kebab/server"
    mkdir -p "$server_dir"

    local router_file="$server_dir/${entity_kebab}.router.ts"
    local service_file="$server_dir/${entity_kebab}.service.ts"
    local repository_file="$server_dir/${entity_kebab}.repository.ts"

    generate_from_template "$SKILL_DIR/assets/trpc-router.md" "$router_file" "$entity"
    echo "✓ Created: $router_file"

    generate_from_template "$SKILL_DIR/assets/trpc-service.md" "$service_file" "$entity"
    echo "✓ Created: $service_file"

    if [[ "$database" == "prisma" ]]; then
      generate_from_template "$SKILL_DIR/assets/prisma-repository.md" "$repository_file" "$entity"
    else
      generate_from_template "$SKILL_DIR/assets/drizzle-repository.md" "$repository_file" "$entity"
    fi
    echo "✓ Created: $repository_file"

  elif [[ "$transport" == "api" ]]; then
    validate_api_client "$target" || true

    local server_dir="$target/src/features/$entity_kebab/server"
    mkdir -p "$server_dir"

    local api_file="$server_dir/${entity_kebab}.api.ts"
    generate_api_file "$entity" > "$api_file"
    echo "✓ Created: $api_file"
  fi
}
