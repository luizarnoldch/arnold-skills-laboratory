#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/lib/helpers.sh"

show_help() {
  cat <<EOF
Next.js Backend Architect Scaffold CLI

Usage:
  $(basename "$0") <target-folder> <entity-name> [options]

Arguments:
  target-folder    Path to the Next.js project root
  entity-name      Entity name in PascalCase (Product) or kebab-case (product)

Options:
   -s, --schema               Scaffold Zod schema
   -S, --server               Scaffold server layer (router + service + repository or api)
   -H, --hooks                Scaffold hooks (queries + mutations + hydration)
   -a, --all                  Scaffold all layers (schema + server + hooks)
   -d, --database <type>      Database type: prisma (default) or drizzle
   -t, --transport <type>     Transport type: trpc (default) or api
   -T, --tree                 Show available scaffold combinations
   -h, --help                 Show this help message

Examples:
  $(basename "$0") /path/to/project Product --schema --database prisma
  $(basename "$0") ./my-app user-profile --schema --database drizzle
  $(basename "$0") . Order --server
  $(basename "$0") . Product --server --transport api
  $(basename "$0") . User --server --transport trpc --database drizzle
  $(basename "$0") . Product --hooks
  $(basename "$0") . Order --hooks --transport api
  $(basename "$0") . User --all --database prisma
  $(basename "$0") . Product --all --transport api --database drizzle

EOF
}

main() {
  local target=""
  local entity=""
  local database=""
  local transport=""
  local scaffold_schema=false
  local scaffold_server=false
  local scaffold_hooks=false
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
      -s|--schema)
        scaffold_schema=true
        shift
        ;;
      -S|--server)
        scaffold_server=true
        shift
        ;;
      -H|--hooks)
        scaffold_hooks=true
        shift
        ;;
      -a|--all)
        scaffold_all=true
        shift
        ;;
      -d|--database)
        if [[ -z "${2:-}" ]]; then
          echo "Error: --database requires a value (prisma or drizzle)" >&2
          exit 1
        fi
        database="$2"
        if [[ ! "$database" =~ ^(prisma|drizzle)$ ]]; then
          echo "Error: Database must be 'prisma' or 'drizzle'" >&2
          exit 1
        fi
        shift 2
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
  
  if [[ "$scaffold_all" == true ]]; then
    if [[ -z "$transport" ]]; then
      transport="trpc"
      echo "→ Using default transport: trpc"
    fi
    
    if [[ -z "$database" ]]; then
      database="prisma"
      echo "→ Using default database: prisma"
    fi
    
    source "$SCRIPT_DIR/schema.sh"
    scaffold_schema "$target" "$entity" "$database"
    
    source "$SCRIPT_DIR/server.sh"
    scaffold_server "$target" "$entity" "$transport" "$database"
    
    source "$SCRIPT_DIR/hooks.sh"
    scaffold_hooks "$target" "$entity" "$transport" "$database"
  fi
  
  if [[ "$scaffold_schema" == true && "$scaffold_all" == false ]]; then
    if [[ -z "$database" ]]; then
      database=$(prompt_database) || exit 1
    fi
    source "$SCRIPT_DIR/schema.sh"
    scaffold_schema "$target" "$entity" "$database"
  fi
  
  if [[ "$scaffold_server" == true && "$scaffold_all" == false ]]; then
    if [[ -z "$transport" ]]; then
      transport="trpc"
      echo "→ Using default transport: trpc"
    fi
    
    if [[ "$transport" == "trpc" && -z "$database" ]]; then
      database="prisma"
      echo "→ Using default database: prisma"
    fi
    
    source "$SCRIPT_DIR/server.sh"
    scaffold_server "$target" "$entity" "$transport" "${database:-prisma}"
  fi
  
  if [[ "$scaffold_hooks" == true && "$scaffold_all" == false ]]; then
    if [[ -z "$transport" ]]; then
      transport="trpc"
      echo "→ Using default transport: trpc"
    fi
    
    if [[ -z "$database" ]]; then
      database="prisma"
      echo "→ Using default database: prisma"
    fi
    
    source "$SCRIPT_DIR/hooks.sh"
    scaffold_hooks "$target" "$entity" "$transport" "$database"
  fi
  
  if [[ "$scaffold_schema" == false && "$scaffold_server" == false && "$scaffold_hooks" == false && "$scaffold_all" == false ]]; then
    echo "Error: No scaffold flag specified. Use --schema, --server, --hooks, or --all." >&2
    show_tree
    exit 1
  fi
}

main "$@"
