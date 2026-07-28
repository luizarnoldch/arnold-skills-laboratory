#!/usr/bin/env bash
set -euo pipefail

to_pascal_case() {
  local input="$1"
  if [[ "$input" =~ ^[A-Z] ]]; then
    echo "$input"
  else
    echo "$input" | sed -E 's/(^|[-_])([a-z])/\U\2/g; s/[-_]//g'
  fi
}

to_kebab_case() {
  local input="$1"
  if [[ "$input" =~ ^[a-z] ]]; then
    echo "$input"
  else
    echo "$input" | sed -E 's/([A-Z])/-\L\1/g; s/^-//; s/([a-z])([A-Z])/\1-\L\2/g'
  fi
}

to_camel_case() {
  local input="$1"
  local pascal
  pascal=$(to_pascal_case "$input")
  echo "$pascal" | sed -E 's/^([A-Z])/\L\1/'
}

# Emits a result line and returns the matching exit code.
#   pass_layer schema
#   fail_layer schema "reason one" "reason two"
pass_layer() {
  local layer="$1"
  echo "LAYER:${layer}:PASS"
  return 0
}

fail_layer() {
  local layer="$1"
  shift
  local reasons
  reasons=$(IFS=';'; echo "$*")
  echo "LAYER:${layer}:FAIL:${reasons}"
  return 1
}

# require_grep <file> <pattern> <label> -- appends to the caller's `reasons` array on miss.
# Caller must declare: local reasons=()
require_grep() {
  local file="$1"
  local pattern="$2"
  local label="$3"
  if [[ ! -f "$file" ]]; then
    reasons+=("$label (file missing: $file)")
    return 1
  fi
  if ! grep -Eq -- "$pattern" "$file"; then
    reasons+=("$label")
    return 1
  fi
  return 0
}

require_file() {
  local file="$1"
  local label="$2"
  if [[ ! -f "$file" ]]; then
    reasons+=("$label (file missing: $file)")
    return 1
  fi
  return 0
}

# extract_schema_fields <schema_file> <entity_camel>
# Extracts field names from the base [entity]Schema z.object({...}) block
extract_schema_fields() {
  local schema_file="$1"
  local entity_camel="$2"
  
  if [[ ! -f "$schema_file" ]]; then
    return 1
  fi
  
  # Find the base schema definition: export const [entity]Schema = z.object({
  # Extract field names from the z.object({...}) block
  sed -n "/export const ${entity_camel}Schema[[:space:]]*=[[:space:]]*z\.object(/,/^})/p" "$schema_file" 2>/dev/null \
    | sed -n '/^[[:space:]]*[a-zA-Z_][a-zA-Z0-9_]*[[:space:]]*:/p' \
    | sed 's/^[[:space:]]*//; s/[[:space:]]*:.*//' \
    | grep -v '^$'
}

# get_schema_create_fields <schema_file> <entity_camel>
# Returns field names that should appear in useCreate defaultValues
# (excludes id, createdAt, updatedAt, deletedAt)
get_schema_create_fields() {
  local schema_file="$1"
  local entity_camel="$2"
  extract_schema_fields "$schema_file" "$entity_camel" | grep -Ev '^(id|createdAt|updatedAt|deletedAt)$'
}

# get_schema_update_fields <schema_file> <entity_camel>
# Returns field names that should appear in useUpdate defaultValues
# (includes id, excludes createdAt, updatedAt, deletedAt)
get_schema_update_fields() {
  local schema_file="$1"
  local entity_camel="$2"
  extract_schema_fields "$schema_file" "$entity_camel" | grep -Ev '^(createdAt|updatedAt|deletedAt)$'
}
