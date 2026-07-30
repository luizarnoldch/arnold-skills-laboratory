#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/lib/helpers.sh"

show_help() {
  cat <<EOF
Next.js Feature Architect Scaffold CLI

Scaffolds frontend layers for a feature-based Next.js app:
pages, views, and components using hooks from next-backend-architect.

Usage:
  $(basename "$0") <target-folder> <entity-name> [options]

Arguments:
  target-folder    Path to the Next.js project root
  entity-name      Entity name in PascalCase (Product) or kebab-case (product)

Options:
  -p, --page <type>           Scaffold page.tsx (list|detail)
  -v, --view [type]           Scaffold view only (list|detail, prompts if omitted)
  -V, --view-full [type]      Scaffold view + components (list|detail, prompts if omitted)
  -a, --all                   Scaffold full frontend (pages + views + components)
  -t, --transport <type>      Transport type: trpc or api (auto-detected if omitted)
  -T, --tree                  Show available scaffold combinations
  -h, --help                  Show this help message

Examples:
  $(basename "$0") /path/to/project Product --page list
  $(basename "$0") . Product --page detail
  $(basename "$0") . Order --view-full list
  $(basename "$0") . Product --view detail
  $(basename "$0") . User --all
  $(basename "$0") . Product --all --transport api

EOF
}

main() {
  local target=""
  local entity=""
  local transport=""
  local page_type=""
  local scaffold_page=false
  local scaffold_view=false
  local scaffold_view_full=false
  local scaffold_all=false

  while [[ $# -gt 0 ]]; do
    case "$1" in
      -h|--help)
        show_help
        exit 0
        ;;
      -T|--tree)
        show_tree
        exit 0
        ;;
      -p|--page)
        if [[ -z "${2:-}" || "$2" =~ ^- ]]; then
          echo "Error: --page requires a value (list or detail)" >&2
          exit 1
        fi
        page_type="$2"
        if [[ ! "$page_type" =~ ^(list|detail)$ ]]; then
          echo "Error: Page type must be 'list' or 'detail'" >&2
          exit 1
        fi
        scaffold_page=true
        shift 2
        ;;
      -v|--view)
        scaffold_view=true
        if [[ -n "${2:-}" && "$2" =~ ^(list|detail)$ ]]; then
          page_type="$2"
          shift
        fi
        shift
        ;;
      -V|--view-full)
        scaffold_view_full=true
        if [[ -n "${2:-}" && "$2" =~ ^(list|detail)$ ]]; then
          page_type="$2"
          shift
        fi
        shift
        ;;
      -a|--all)
        scaffold_all=true
        shift
        ;;
      -t|--transport)
        if [[ -z "${2:-}" ]]; then
          echo "Error: --transport requires a value (trpc or api)" >&2
          exit 1
        fi
        transport="$2"
        if [[ ! "$transport" =~ ^(trpc|api)$ ]]; then
          echo "Error: Transport must be 'trpc' or 'api'" >&2
          exit 1
        fi
        shift 2
        ;;
      -*)
        echo "Error: Unknown option: $1" >&2
        show_help
        exit 1
        ;;
      *)
        if [[ -z "$target" ]]; then
          target="$1"
        elif [[ -z "$entity" ]]; then
          entity="$1"
        else
          echo "Error: Unexpected argument: $1" >&2
          exit 1
        fi
        shift
        ;;
    esac
  done

  if [[ -z "$target" ]]; then
    echo "Error: Target folder is required" >&2
    show_help
    exit 1
  fi

  validate_path "$target" || exit 1

  if [[ -z "$entity" ]]; then
    entity=$(prompt_entity) || exit 1
  fi

  validate_entity "$entity" || exit 1

  local entity_kebab
  entity_kebab=$(to_kebab_case "$entity")

  if [[ -z "$transport" ]]; then
    transport=$(detect_transport "$target" "$entity_kebab")
    echo "→ Detected transport: $transport"
  fi

  if [[ "$scaffold_all" == true ]]; then
    source "$SCRIPT_DIR/views.sh"
    scaffold_all_frontend "$target" "$entity" "$transport"
  fi

  if [[ "$scaffold_page" == true && "$scaffold_all" == false ]]; then
    source "$SCRIPT_DIR/pages.sh"
    scaffold_pages "$target" "$entity" "$page_type" "$transport"
  fi

  if [[ "$scaffold_view" == true && "$scaffold_all" == false ]]; then
    source "$SCRIPT_DIR/views.sh"
    scaffold_views "$target" "$entity" "$page_type" "$transport"
  fi

  if [[ "$scaffold_view_full" == true && "$scaffold_all" == false ]]; then
    source "$SCRIPT_DIR/views.sh"
    scaffold_views_and_components "$target" "$entity" "$page_type" "$transport"
  fi

  if [[ "$scaffold_page" == false && "$scaffold_view" == false && "$scaffold_view_full" == false && "$scaffold_all" == false ]]; then
    echo "Error: No scaffold flag specified. Use --page, --view, --view-full, or --all." >&2
    show_tree
    exit 1
  fi
}

main "$@"
