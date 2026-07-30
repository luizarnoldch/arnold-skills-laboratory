#!/usr/bin/env bash
set -euo pipefail

REPLACER_LIB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$REPLACER_LIB_DIR/helpers.sh"

replace_placeholders() {
  local content="$1"
  local entity="$2"

  local entity_pascal
  local entity_camel
  local entity_kebab
  local entity_snake

  entity_pascal=$(to_pascal_case "$entity")
  entity_camel=$(to_camel_case "$entity")
  entity_kebab=$(to_kebab_case "$entity")
  entity_snake=$(to_snake_case "$entity")

  content="${content//\[Entity\]/$entity_pascal}"
  content="${content//\[entity\]/$entity_camel}"
  content="${content//\[entityTable\]/$entity_snake}"
  content="${content//\[entity-kebab\]/$entity_kebab}"

  echo "$content"
}
