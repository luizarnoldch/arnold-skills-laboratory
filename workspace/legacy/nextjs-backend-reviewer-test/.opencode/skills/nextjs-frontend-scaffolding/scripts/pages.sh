#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/lib/helpers.sh"
source "$SCRIPT_DIR/lib/replacer.sh"

generate_page() {
  local target="$1"
  local entity="$2"
  local page_type="$3"
  local transport="${4:-trpc}"

  local entity_pascal
  local entity_camel
  local entity_kebab

  entity_pascal=$(to_pascal_case "$entity")
  entity_camel=$(to_camel_case "$entity")
  entity_kebab=$(to_kebab_case "$entity")

  if [[ "$page_type" == "list" ]]; then
    local page_dir="$target/src/app/${entity_kebab}s"
    mkdir -p "$page_dir"

    local page_file="$page_dir/page.tsx"

    cat > "$page_file" <<PAGETSX
import Hydrate${entity_pascal}s from "@/features/${entity_kebab}/hooks/Hydrate${entity_pascal}s"
import ${entity_pascal}View from "@/features/${entity_kebab}/views/${entity_pascal}View"

type ${entity_pascal}sPageProps = {
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>
}

const ${entity_pascal}sPage = async ({ searchParams }: ${entity_pascal}sPageProps) => {
  return (
    <Hydrate${entity_pascal}s>
      <${entity_pascal}View />
    </Hydrate${entity_pascal}s>
  )
}

export default ${entity_pascal}sPage
PAGETSX

    echo "✓ Created: $page_file"

  elif [[ "$page_type" == "detail" ]]; then
    local page_dir="$target/src/app/${entity_kebab}s/[id]"
    mkdir -p "$page_dir"

    local page_file="$page_dir/page.tsx"

    cat > "$page_file" <<PAGETSX
import Hydrate${entity_pascal} from "@/features/${entity_kebab}/hooks/Hydrate${entity_pascal}"
import ${entity_pascal}DetailView from "@/features/${entity_kebab}/views/${entity_pascal}DetailView"

type ${entity_pascal}DetailPageProps = {
  params: Promise<{ id: string }>
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>
}

const ${entity_pascal}DetailPage = async ({ params, searchParams }: ${entity_pascal}DetailPageProps) => {
  const { id } = await params

  return (
    <Hydrate${entity_pascal} ${entity_camel}Id={id}>
      <${entity_pascal}DetailView ${entity_camel}Id={id} />
    </Hydrate${entity_pascal}>
  )
}

export default ${entity_pascal}DetailPage
PAGETSX

    echo "✓ Created: $page_file"
  else
    echo "Error: Unknown page type: $page_type" >&2
    return 1
  fi
}

scaffold_pages() {
  local target="$1"
  local entity="$2"
  local page_type="$3"
  local transport="${4:-trpc}"

  if [[ -z "$page_type" ]]; then
    page_type=$(prompt_page_type) || exit 1
  fi

  generate_page "$target" "$entity" "$page_type" "$transport"
}
