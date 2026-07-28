#!/usr/bin/env bash

# Validate [Entity]FormCreate component

check_form_create() {
  local target="$1"
  local entity_pascal="$2"
  local entity_camel="$3"
  local schema_fields="$4"
  local transport="$5"
  
  local failures=()
  
  local form_path="${target}/src/features/${entity_camel}/components/${entity_pascal}FormCreate.tsx"
  
  # Check 1: File exists
  if [[ ! -f "$form_path" ]]; then
    echo "COMPONENT:form-create:FAIL:File not found at expected path"
    return 1
  fi
  
  # Check 2: Hook import present
  if ! check_import "$form_path" "useCreate${entity_pascal}"; then
    failures+=("Missing hook import: useCreate${entity_pascal}")
  fi
  
  # Check 3: Form usage
  if ! contains_pattern "$form_path" "form\.handleSubmit|handleSubmit"; then
    failures+=("No form.handleSubmit detected")
  fi
  
  # Check 4: Template comments removed
  if has_template_comment "$form_path"; then
    failures+=("Template comments still present - form fields not filled in")
  fi
  
  # Check 5: Field coverage
  if [[ -n "$schema_fields" ]]; then
    local form_fields=$(extract_form_fields "$form_path")
    
    # If no fields detected, that's a problem
    if [[ -z "$form_fields" ]]; then
      failures+=("No form fields detected - only submit buttons present")
    else
      # Check for missing fields
      local missing_fields=$(find_missing_fields "$schema_fields" "$form_fields")
      if [[ -n "$missing_fields" ]]; then
        failures+=("Missing form fields: ${missing_fields//,/, }")
      fi
    fi
  fi
  
  # Report results
  if [[ ${#failures[@]} -eq 0 ]]; then
    echo "COMPONENT:form-create:PASS"
  else
    local failure_msg=$(IFS=';'; echo "${failures[*]}")
    echo "COMPONENT:form-create:FAIL:${failure_msg}"
  fi
}
