#!/usr/bin/env bash
set -euo pipefail

extract_prisma_fields() {
  local prisma_file="$1"
  local model_name="$2"

  if [[ ! -f "$prisma_file" ]]; then
    echo "PRISMA_PARSE:FAIL:prisma schema file not found: $prisma_file" >&2
    return 1
  fi

  local model_start model_end
  model_start=$(grep -n "^model ${model_name}\b" "$prisma_file" | head -1 | cut -d: -f1)
  if [[ -z "$model_start" ]]; then
    echo "PRISMA_PARSE:FAIL:model '${model_name}' not found in $prisma_file" >&2
    return 1
  fi

  model_end=$(tail -n +"$model_start" "$prisma_file" | grep -n '^}' | head -1 | cut -d: -f1)
  if [[ -z "$model_end" ]]; then
    echo "PRISMA_PARSE:FAIL:model '${model_name}' closing brace not found" >&2
    return 1
  fi
  model_end=$((model_start + model_end - 1))

  local block
  block=$(sed -n "${model_start},${model_end}p" "$prisma_file")

  while IFS= read -r line; do
    [[ -z "$line" ]] && continue
    [[ "$line" =~ ^[[:space:]]*\/\/ ]] && continue
    [[ "$line" =~ ^[[:space:]]*model ]] && continue
    [[ "$line" =~ ^[[:space:]]*\} ]] && continue

    local field_name
    field_name=$(echo "$line" | awk '{print $1}')
    if [[ -z "$field_name" || "$field_name" =~ ^(model|\}|\/\/) ]]; then
      continue
    fi

    local field_type
    field_type=$(echo "$line" | awk '{print $2}')

    local is_optional=false
    if [[ "$field_type" =~ \?$ ]]; then
      is_optional=true
      field_type="${field_type%\?}"
    fi

    local is_db_managed=false
    if echo "$line" | grep -qE '@default\((cuid|uuid|autoincrement)\(\)\)'; then
      is_db_managed=true
    fi
    if echo "$line" | grep -q '@updatedAt'; then
      is_db_managed=true
    fi
    if echo "$line" | grep -q '@default(now())' && [[ "$field_name" =~ ^(createdAt|updatedAt|deletedAt)$ ]]; then
      is_db_managed=true
    fi

    local has_default=false
    if echo "$line" | grep -q '@default('; then
      has_default=true
    fi

    local is_id=false
    if echo "$line" | grep -q '@id'; then
      is_id=true
    fi

    echo "${field_name}|${field_type}|${is_optional}|${is_db_managed}|${has_default}|${is_id}"
  done <<< "$block"

  echo "PRISMA_PARSE:PASS"
  return 0
}

get_prisma_all_fields() {
  local prisma_file="$1"
  local model_name="$2"
  extract_prisma_fields "$prisma_file" "$model_name" 2>/dev/null | grep -v '^PRISMA_PARSE:' | cut -d'|' -f1
}

get_prisma_create_fields() {
  local prisma_file="$1"
  local model_name="$2"
  extract_prisma_fields "$prisma_file" "$model_name" 2>/dev/null | grep -v '^PRISMA_PARSE:' | while IFS='|' read -r name type optional db_managed has_default is_id; do
    if [[ "$db_managed" == "false" ]]; then
      echo "$name"
    fi
  done
}

get_prisma_update_fields() {
  local prisma_file="$1"
  local model_name="$2"
  extract_prisma_fields "$prisma_file" "$model_name" 2>/dev/null | grep -v '^PRISMA_PARSE:' | while IFS='|' read -r name type optional db_managed has_default is_id; do
    if [[ "$db_managed" == "false" || "$name" == "id" ]]; then
      echo "$name"
    fi
  done
}
