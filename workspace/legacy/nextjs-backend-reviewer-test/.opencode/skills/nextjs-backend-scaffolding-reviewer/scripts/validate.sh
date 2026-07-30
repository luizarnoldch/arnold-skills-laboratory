#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/lib/helpers.sh"
source "$SCRIPT_DIR/lib/prisma.sh"
source "$SCRIPT_DIR/checks/schema.sh"
source "$SCRIPT_DIR/checks/router.sh"
source "$SCRIPT_DIR/checks/service.sh"
source "$SCRIPT_DIR/checks/repository.sh"
source "$SCRIPT_DIR/checks/hooks.sh"

show_help() {
  cat <<EOF
Next.js Backend Reviewer Validation CLI

Usage:
  $(basename "$0") <target-folder> <entity-name> [options]

Arguments:
  target-folder    Path to the Next.js project root
  entity-name      Entity name in PascalCase (Product) or kebab-case (product)

Options:
  -t, --transport <type>   Transport type: trpc (default) or api
  -d, --database <type>    Database type: prisma (default) or drizzle
  -P, --prisma-schema <path>  Path to Prisma schema file (default: <target>/prisma/schema.prisma)
  -T, --typecheck          Also run 'npx tsc --noEmit' in the target folder
  -h, --help               Show this help message

Model Field Validation:
  When --prisma-schema is provided (or auto-detected), the validator
  cross-references fields in the Zod schema, useCreate, and useUpdate
  defaultValues against the actual Prisma model definition to catch
  missing, extra, or placeholder fields.

Output:
  One "LAYER:<layer>:PASS" or "LAYER:<layer>:FAIL:<reasons>" line per layer,
  optionally "TYPECHECK:PASS|FAIL|SKIPPED", then a final "VERDICT:PASS|FAIL".

  REST transport (--transport api) only checks schema, router (the api.ts
  client wrapper), and hooks -- there is no service.ts/repository.ts under
  REST transport.
EOF
}

main() {
  local target="" entity="" transport="trpc" database="prisma" run_typecheck=false
  local prisma_schema=""

  while [[ $# -gt 0 ]]; do
    case "$1" in
      -h|--help)
        show_help
        exit 0
        ;;
      -t|--transport)
        transport="${2:-}"
        if [[ ! "$transport" =~ ^(trpc|api)$ ]]; then
          echo "Error: --transport must be 'trpc' or 'api'" >&2
          exit 1
        fi
        shift 2
        ;;
      -d|--database)
        database="${2:-}"
        if [[ ! "$database" =~ ^(prisma|drizzle)$ ]]; then
          echo "Error: --database must be 'prisma' or 'drizzle'" >&2
          exit 1
        fi
        shift 2
        ;;
      -P|--prisma-schema)
        prisma_schema="${2:-}"
        if [[ -z "$prisma_schema" ]]; then
          echo "Error: --prisma-schema requires a file path" >&2
          exit 1
        fi
        shift 2
        ;;
      -T|--typecheck)
        run_typecheck=true
        shift
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

  if [[ -z "$target" || -z "$entity" ]]; then
    echo "Error: target-folder and entity-name are required" >&2
    show_help
    exit 1
  fi

  if [[ ! -d "$target" ]]; then
    echo "VERDICT:FAIL:target directory not found: $target"
    exit 1
  fi

  local entity_pascal entity_kebab entity_camel
  entity_pascal=$(to_pascal_case "$entity")
  entity_kebab=$(to_kebab_case "$entity")
  entity_camel=$(to_camel_case "$entity")

  local feature_dir="$target/src/features/$entity_kebab"
  if [[ ! -d "$feature_dir" ]]; then
    echo "VERDICT:FAIL:backend directory missing -- layers were not generated ($feature_dir)"
    exit 1
  fi

  # Resolve Prisma schema path
  if [[ -z "$prisma_schema" ]]; then
    local default_schema="$target/prisma/schema.prisma"
    if [[ -f "$default_schema" ]]; then
      prisma_schema="$default_schema"
    fi
  fi

  if [[ -n "$prisma_schema" && ! -f "$prisma_schema" ]]; then
    echo "Warning: --prisma-schema path does not exist: $prisma_schema -- skipping model field cross-check" >&2
    prisma_schema=""
  fi

  local overall_pass=true

  check_schema "$target" "$entity_kebab" "$entity_pascal" "$entity_camel" "$database" "$prisma_schema" || overall_pass=false
  check_router "$target" "$entity_kebab" "$entity_pascal" "$entity_camel" "$transport" || overall_pass=false

  if [[ "$transport" == "trpc" ]]; then
    check_service "$target" "$entity_kebab" "$entity_pascal" "$entity_camel" || overall_pass=false
    check_repository "$target" "$entity_kebab" "$entity_pascal" "$entity_camel" "$database" || overall_pass=false
  fi

  check_hooks "$target" "$entity_kebab" "$entity_pascal" "$entity_camel" "$transport" "$prisma_schema" || overall_pass=false

  if [[ "$run_typecheck" == true ]]; then
    if [[ -f "$target/tsconfig.json" ]]; then
      local tsc_output
      if tsc_output=$(cd "$target" && npx tsc --noEmit 2>&1); then
        echo "TYPECHECK:PASS"
      else
        echo "TYPECHECK:FAIL:$(echo "$tsc_output" | head -n 10 | tr '\n' '|')"
        overall_pass=false
      fi
    else
      echo "TYPECHECK:SKIPPED:no tsconfig.json"
    fi
  fi

  if [[ "$overall_pass" == true ]]; then
    echo "VERDICT:PASS"
  else
    echo "VERDICT:FAIL"
    exit 1
  fi
}

main "$@"
