#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/lib/helpers.sh"
source "$SCRIPT_DIR/checks/list.sh"
source "$SCRIPT_DIR/checks/form-create.sh"
source "$SCRIPT_DIR/checks/form-update.sh"

show_help() {
  cat <<EOF
Next.js Frontend Reviewer Validation CLI

Usage:
  $(basename "$0") <target-folder> <entity-name> [options]

Arguments:
  target-folder    Path to the Next.js project root
  entity-name      Entity name in PascalCase (Product) or kebab-case (product)

Options:
  -t, --transport <type>   Transport type: trpc (default) or api
  -h, --help               Show this help message

Component Validation:
  Validates List, FormCreate, and FormUpdate components for:
  - File existence and structure
  - Proper hook imports and usage
  - Schema field coverage in tables and forms
  - Template placeholder removal

Output:
  One "COMPONENT:<component>:PASS" or "COMPONENT:<component>:FAIL:<reasons>" line per component,
  then a final "VERDICT:PASS|FAIL".
EOF
}

main() {
  local target="" entity="" transport="trpc"

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
    echo "Error: Missing required arguments" >&2
    show_help
    exit 1
  fi

  # Normalize entity name to PascalCase
  local entity_pascal="$entity"
  # If entity starts with lowercase, convert first letter to uppercase
  if [[ "$entity" =~ ^[a-z] ]]; then
    entity_pascal="$(echo "${entity:0:1}" | tr '[:lower:]' '[:upper:]')${entity:1}"
  fi

  local entity_camel=$(to_camel_case "$entity_pascal")
  local entity_kebab=$(to_kebab_case "$entity_pascal")

  # Validate target directory
  if [[ ! -d "$target" ]]; then
    echo "VERDICT:FAIL:target directory not found: $target"
    exit 1
  fi

  # Check if frontend feature directory exists
  local feature_dir="${target}/src/features/${entity_camel}"
  if [[ ! -d "$feature_dir" ]]; then
    echo "VERDICT:FAIL:frontend directory missing: $feature_dir"
    exit 1
  fi

  # Extract schema fields
  local schema_file="${feature_dir}/schemas/${entity_kebab}.schema.ts"
  local schema_fields=""
  
  if [[ -f "$schema_file" ]]; then
    schema_fields=$(extract_schema_fields "$schema_file")
  else
    echo "VERDICT:FAIL:schema file not found: $schema_file"
    exit 1
  fi

  # Run component checks
  local has_failures=false

  check_list_component "$target" "$entity_pascal" "$entity_camel" "$entity_kebab" "$schema_fields" "$transport"
  if [[ $? -ne 0 ]] || grep -q "COMPONENT:list:FAIL" <(check_list_component "$target" "$entity_pascal" "$entity_camel" "$entity_kebab" "$schema_fields" "$transport" 2>&1); then
    has_failures=true
  fi

  check_form_create "$target" "$entity_pascal" "$entity_camel" "$schema_fields" "$transport"
  if [[ $? -ne 0 ]] || grep -q "COMPONENT:form-create:FAIL" <(check_form_create "$target" "$entity_pascal" "$entity_camel" "$schema_fields" "$transport" 2>&1); then
    has_failures=true
  fi

  check_form_update "$target" "$entity_pascal" "$entity_camel" "$schema_fields" "$transport"
  if [[ $? -ne 0 ]] || grep -q "COMPONENT:form-update:FAIL" <(check_form_update "$target" "$entity_pascal" "$entity_camel" "$schema_fields" "$transport" 2>&1); then
    has_failures=true
  fi

  # Final verdict
  if [[ "$has_failures" == true ]]; then
    echo "VERDICT:FAIL"
    exit 1
  else
    echo "VERDICT:PASS"
    exit 0
  fi
}

main "$@"
