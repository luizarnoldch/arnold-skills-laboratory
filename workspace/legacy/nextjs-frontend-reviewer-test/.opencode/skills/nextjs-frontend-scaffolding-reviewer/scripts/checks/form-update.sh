#!/usr/bin/env bash

# Validate [Entity]FormUpdate component

check_form_update() {
  local target="$1"
  local entity_pascal="$2"
  local entity_camel="$3"
  local schema_fields="$4"
  local transport="$5"
  
  local failures=()
  
  local form_path="${target}/src/features/${entity_camel}/components/${entity_pascal}FormUpdate.tsx"
  
  # Check 1: File exists
  if [[ ! -f "$form_path" ]]; then
    echo "COMPONENT:form-update:FAIL:File not found at expected path"
    return 1
  fi
  
  # Check 2: Hook import present
  if ! check_import "$form_path" "useUpdate${entity_pascal}"; then
    failures+=("Missing hook import: useUpdate${entity_pascal}")
  fi
  
  # Check 3: Entity prop present
  if ! contains_pattern "$form_path" "${entity_camel}.*:.*${entity_pascal}|${entity_camel}:.*type.*${entity_pascal}"; then
    failures+=("Missing or improperly typed entity prop (${entity_camel})")
  fi
  
  # Check 4: Form usage
  if ! contains_pattern "$form_path" "form\.handleSubmit|handleSubmit"; then
    failures+=("No form.handleSubmit detected")
  fi
  
  # Check 5: Template comments removed
  if has_template_comment "$form_path"; then
    failures+=("Template comments still present - form fields not filled in")
  fi
  
  # Check 6: Field coverage (should include id for updates)
  if [[ -n "$schema_fields" ]]; then
    local form_fields=$(extract_form_fields "$form_path")
    
    # If no fields detected, that's a problem
    if [[ -z "$form_fields" ]]; then
      failures+=("No form fields detected - only submit buttons present")
    else
      # For update forms, we expect the same fields as create
      # (id is usually passed via prop, not as a form field)
      local missing_fields=$(find_missing_fields "$schema_fields" "$form_fields")
      if [[ -n "$missing_fields" ]]; then
        failures+=("Missing form fields: ${missing_fields//,/, }")
      fi
    fi
  fi
  
  # Report results
  if [[ ${#failures[@]} -eq 0 ]]; then
    echo "COMPONENT:form-update:PASS"
  else
    local failure_msg=$(IFS=';'; echo "${failures[*]}")
    echo "COMPONENT:form-update:FAIL:${failure_msg}"
  fi
}
