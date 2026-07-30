#!/usr/bin/env bash

# Validate [Entity]List component

check_list_component() {
  local target="$1"
  local entity_pascal="$2"
  local entity_camel="$3"
  local entity_kebab="$4"
  local schema_fields="$5"
  local transport="$6"
  
  local failures=()
  
  # Try to find the list component
  local list_path="${target}/src/features/${entity_camel}/components/${entity_pascal}List/index.tsx"
  
  if [[ ! -f "$list_path" ]]; then
    # Try alternative location
    list_path="${target}/src/features/${entity_camel}/components/${entity_pascal}List.tsx"
  fi
  
  # Check 1: File exists
  if [[ ! -f "$list_path" ]]; then
    echo "COMPONENT:list:FAIL:File not found at expected path"
    return 1
  fi
  
  # Check 2: Schema import present
  if ! check_import "$list_path" "from.*schemas/${entity_kebab}.schema"; then
    failures+=("Missing schema import")
  fi
  
  # Check 3: Hook usage (useSuspenseList or useList)
  if ! contains_pattern "$list_path" "useSuspenseList${entity_pascal}s|useList${entity_pascal}s"; then
    failures+=("Missing list hook usage (useSuspenseList${entity_pascal}s or useList${entity_pascal}s)")
  fi
  
  # Check 4: Delete hook present
  if ! contains_pattern "$list_path" "useDelete${entity_pascal}"; then
    failures+=("Missing delete hook (useDelete${entity_pascal})")
  fi
  
  # Check 5: Field display in table
  if ! check_table_fields "$list_path" "$schema_fields"; then
    failures+=("Table only shows ID field - no custom schema fields displayed")
  fi
  
  # Check 6: Template comments removed
  if has_template_comment "$list_path"; then
    failures+=("Template comments still present")
  fi
  
  # Check 7: Form integration
  if ! contains_pattern "$list_path" "${entity_pascal}FormCreate"; then
    failures+=("Missing ${entity_pascal}FormCreate integration")
  fi
  
  if ! contains_pattern "$list_path" "${entity_pascal}FormUpdate"; then
    failures+=("Missing ${entity_pascal}FormUpdate integration")
  fi
  
  # Report results
  if [[ ${#failures[@]} -eq 0 ]]; then
    echo "COMPONENT:list:PASS"
  else
    local failure_msg=$(IFS=';'; echo "${failures[*]}")
    echo "COMPONENT:list:FAIL:${failure_msg}"
  fi
}
