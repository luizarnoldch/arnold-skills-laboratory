#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/lib/helpers.sh"
source "$SCRIPT_DIR/lib/replacer.sh"
source "$SCRIPT_DIR/generate-from-template.sh"

TEMPLATES_DIR="$(dirname "$SCRIPT_DIR")/assets/templates"

generate_view() {
  local target="$1"
  local entity="$2"
  local view_type="$3"
  local transport="${4:-trpc}"

  local entity_pascal
  local entity_kebab

  entity_pascal=$(to_pascal_case "$entity")
  entity_kebab=$(to_kebab_case "$entity")

  local views_dir="$target/src/features/$entity_kebab/views"
  mkdir -p "$views_dir"

  if [[ "$view_type" == "list" ]]; then
    generate_from_template \
      "$TEMPLATES_DIR/view-list.tsx" \
      "$entity" \
      "$views_dir/${entity_pascal}View.tsx"

  elif [[ "$view_type" == "detail" ]]; then
    local view_file="$views_dir/${entity_pascal}DetailView.tsx"
    local entity_camel
    entity_camel=$(to_camel_case "$entity")

    cat > "$view_file" <<VIEWTSX
import { Suspense } from "react"
import { ErrorBoundary } from "react-error-boundary"

import ${entity_pascal}Detail from "../components/${entity_pascal}Detail"
import ${entity_pascal}DetailViewLoader from "../components/loaders/${entity_pascal}DetailViewLoader"
import ${entity_pascal}DetailViewError from "../components/error/${entity_pascal}DetailViewError"

type ${entity_pascal}DetailViewProps = {
  ${entity_camel}Id: string
}

const ${entity_pascal}DetailView = ({ ${entity_camel}Id }: ${entity_pascal}DetailViewProps) => {
  return (
    <div className="container mx-auto p-6">
      <ErrorBoundary fallback={<${entity_pascal}DetailViewError />}>
        <Suspense fallback={<${entity_pascal}DetailViewLoader />}>
          <${entity_pascal}Detail ${entity_camel}Id={${entity_camel}Id} />
        </Suspense>
      </ErrorBoundary>
    </div>
  )
}

export default ${entity_pascal}DetailView
VIEWTSX

    echo "✓ Created: $view_file"
  else
    echo "Error: Unknown view type: $view_type" >&2
    return 1
  fi
}

generate_components() {
  local target="$1"
  local entity="$2"
  local page_type="$3"

  local entity_pascal
  local entity_kebab

  entity_pascal=$(to_pascal_case "$entity")
  entity_kebab=$(to_kebab_case "$entity")

  local comp_dir="$target/src/features/$entity_kebab/components"
  local list_dir="$comp_dir/${entity_pascal}List"

  if [[ "$page_type" == "list" ]]; then
    generate_from_template \
      "$TEMPLATES_DIR/list-component.tsx" \
      "$entity" \
      "$list_dir/index.tsx"

    generate_from_template \
      "$TEMPLATES_DIR/form-create.tsx" \
      "$entity" \
      "$comp_dir/${entity_pascal}FormCreate.tsx"

    generate_from_template \
      "$TEMPLATES_DIR/form-update.tsx" \
      "$entity" \
      "$comp_dir/${entity_pascal}FormUpdate.tsx"

  elif [[ "$page_type" == "detail" ]]; then
    local entity_camel
    entity_camel=$(to_camel_case "$entity")

    local loaders_dir="$comp_dir/loaders"
    local error_dir="$comp_dir/error"
    local detail_dir="$comp_dir/${entity_pascal}Detail"

    mkdir -p "$loaders_dir" "$error_dir" "$detail_dir"

    cat > "$loaders_dir/${entity_pascal}DetailViewLoader.tsx" <<LOADER
const ${entity_pascal}DetailViewLoader = () => {
  return <div>Loading...</div>
}

export default ${entity_pascal}DetailViewLoader
LOADER
    echo "✓ Created: $loaders_dir/${entity_pascal}DetailViewLoader.tsx"

    cat > "$error_dir/${entity_pascal}DetailViewError.tsx" <<ERROR
const ${entity_pascal}DetailViewError = () => {
  return <div>Failed to load</div>
}

export default ${entity_pascal}DetailViewError
ERROR
    echo "✓ Created: $error_dir/${entity_pascal}DetailViewError.tsx"

    cat > "$detail_dir/index.tsx" <<DETAIL
"use client"

import useSuspenseList${entity_pascal}s from "../../hooks/useSuspenseList${entity_pascal}s"

type ${entity_pascal}DetailProps = {
  ${entity_camel}Id: string
}

const ${entity_pascal}Detail = ({ ${entity_camel}Id }: ${entity_pascal}DetailProps) => {
  const { ${entity_camel}s } = useSuspenseList${entity_pascal}s()
  const item = ${entity_camel}s.find((i) => i.id === ${entity_camel}Id)

  if (!item) {
    return <div>Not found</div>
  }

  return (
    <div>
      <h1>{item.name}</h1>
    </div>
  )
}

export default ${entity_pascal}Detail
DETAIL
    echo "✓ Created: $detail_dir/index.tsx"
  fi
}

scaffold_views() {
  local target="$1"
  local entity="$2"
  local view_type="$3"
  local transport="${4:-trpc}"

  if [[ -z "$view_type" ]]; then
    view_type=$(prompt_page_type) || exit 1
  fi

  generate_view "$target" "$entity" "$view_type" "$transport"
}

scaffold_views_and_components() {
  local target="$1"
  local entity="$2"
  local page_type="$3"
  local transport="${4:-trpc}"

  if [[ -z "$page_type" ]]; then
    page_type=$(prompt_page_type) || exit 1
  fi

  generate_view "$target" "$entity" "$page_type" "$transport"
  generate_components "$target" "$entity" "$page_type"
}

scaffold_all_frontend() {
  local target="$1"
  local entity="$2"
  local transport="${3:-trpc}"

  local hooks_exist
  hooks_exist=$(validate_hooks_exist "$target" "$entity")

  if [[ "$hooks_exist" == "false" ]]; then
    echo "Warning: No hooks found for $entity. Run next-backend-architect first." >&2
    echo "  $(dirname "$SCRIPT_DIR")/../../next-backend-architect/scripts/main.sh $target $entity --hooks"
    return 1
  fi

  source "$SCRIPT_DIR/pages.sh"
  generate_page "$target" "$entity" "list" "$transport"

  generate_view "$target" "$entity" "list" "$transport"
  generate_components "$target" "$entity" "list"
}
