#!/usr/bin/env bash

# Helper functions for frontend validation

# Convert PascalCase to camelCase
# Note: Use printf to avoid newline between first char and rest
to_camel_case() {
  local str="$1"
  printf '%s%s' "$(echo "${str:0:1}" | tr '[:upper:]' '[:lower:]')" "${str:1}"
}

# Convert PascalCase to kebab-case
to_kebab_case() {
  local str="$1"
  echo "$str" | sed 's/\([A-Z]\)/-\1/g' | tr '[:upper:]' '[:lower:]' | sed 's/^-//'
}

# Extract schema fields from a schema file
# Returns comma-separated list of field names (excludes id, createdAt, updatedAt, deletedAt)
extract_schema_fields() {
  local schema_file="$1"
  
  if [[ ! -f "$schema_file" ]]; then
    echo ""
    return 1
  fi
  
  # Extract the z.object({...}) block and parse field names
  # This is a simple parser that looks for lines like "  fieldName: z.string()"
  local fields=""
  local in_object=false
  
  while IFS= read -r line; do
    # Detect start of z.object({
    if [[ "$line" =~ z\.object\(\{ ]]; then
      in_object=true
      continue
    fi
    
    # Detect end of object })
    if [[ "$in_object" == true && "$line" =~ ^\s*\}\) ]]; then
      break
    fi
    
    # Extract field names
    if [[ "$in_object" == true && "$line" =~ ^[[:space:]]*([a-zA-Z_][a-zA-Z0-9_]*)[[:space:]]*:[[:space:]]*z\. ]]; then
      local field_name="${BASH_REMATCH[1]}"
      
      # Skip auto-generated fields
      if [[ "$field_name" != "id" && "$field_name" != "createdAt" && "$field_name" != "updatedAt" && "$field_name" != "deletedAt" ]]; then
        if [[ -z "$fields" ]]; then
          fields="$field_name"
        else
          fields="$fields,$field_name"
        fi
      fi
    fi
  done < "$schema_file"
  
  echo "$fields"
}

# Check if a file contains a pattern
contains_pattern() {
  local file="$1"
  local pattern="$2"
  
  if [[ ! -f "$file" ]]; then
    return 1
  fi
  
  grep -qE "$pattern" "$file"
}

# Check if file contains import statement
check_import() {
  local file="$1"
  local import_pattern="$2"
  
  contains_pattern "$file" "import.*$import_pattern"
}

# Check if file contains template comment
has_template_comment() {
  local file="$1"
  
  if [[ ! -f "$file" ]]; then
    return 1
  fi
  
  grep -qE '\{/\*.*Add.*for each.*field.*\*/\}|{/\*.*TODO.*\*/}|{/\*.*PLACEHOLDER.*\*/}' "$file"
}

# Extract form field names from form.Field components or input elements
# Returns comma-separated list of field names
extract_form_fields() {
  local file="$1"
  
  if [[ ! -f "$file" ]]; then
    echo ""
    return 1
  fi
  
  local fields=""
  
  # Look for form.Field with name attribute: form.Field name="fieldName"
  # Use variable to avoid quote-escaping issues with regex containing [^>]
  local _re_form='form\.Field[^>]*name=["'\'']([a-zA-Z_][a-zA-Z0-9_]*)["'\'']'
  while IFS= read -r line; do
    if [[ "$line" =~ $_re_form ]]; then
      local field_name="${BASH_REMATCH[1]}"
      if [[ -z "$fields" ]]; then
        fields="$field_name"
      else
        fields="$fields,$field_name"
      fi
    fi
  done < "$file"
  
  # Also look for regular input elements with name attribute
  local _re_input='<input[^>]*name=["'\'']([a-zA-Z_][a-zA-Z0-9_]*)["'\'']'
  while IFS= read -r line; do
    if [[ "$line" =~ $_re_input ]]; then
      local field_name="${BASH_REMATCH[1]}"
      # Avoid duplicates
      if [[ ! ",$fields," =~ ,$field_name, ]]; then
        if [[ -z "$fields" ]]; then
          fields="$field_name"
        else
          fields="$fields,$field_name"
        fi
      fi
    fi
  done < "$file"
  
  echo "$fields"
}

# Check if table displays custom fields (beyond just ID)
# Returns 0 if custom fields found, 1 otherwise
check_table_fields() {
  local file="$1"
  local schema_fields="$2"
  
  if [[ ! -f "$file" ]]; then
    return 1
  fi
  
  if [[ -z "$schema_fields" ]]; then
    # No schema fields to check against, just verify not only showing ID
    # Look for more than one <th> element
    local th_count=$(grep -oE '<th[^>]*>' "$file" | wc -l)
    if [[ "$th_count" -gt 2 ]]; then
      # More than just ID and Actions columns
      return 0
    fi
    return 1
  fi
  
  # Check if any schema field is referenced in the component
  IFS=',' read -ra FIELDS <<< "$schema_fields"
  for field in "${FIELDS[@]}"; do
    # Look for references like [entity].$field or just $field as text
    if grep -qE "\[.*\]\.$field|>$field<|\"$field\"" "$file"; then
      return 0
    fi
  done
  
  return 1
}

# Compare two comma-separated lists and return missing items from first list
find_missing_fields() {
  local expected="$1"  # comma-separated
  local actual="$2"    # comma-separated
  
  if [[ -z "$expected" ]]; then
    echo ""
    return 0
  fi
  
  local missing=""
  IFS=',' read -ra EXPECTED_FIELDS <<< "$expected"
  
  for field in "${EXPECTED_FIELDS[@]}"; do
    if [[ ! ",$actual," =~ ,$field, ]]; then
      if [[ -z "$missing" ]]; then
        missing="$field"
      else
        missing="$missing,$field"
      fi
    fi
  done
  
  echo "$missing"
}
