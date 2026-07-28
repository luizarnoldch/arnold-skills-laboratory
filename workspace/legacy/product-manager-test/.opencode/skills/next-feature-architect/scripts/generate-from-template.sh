#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/lib/helpers.sh"
source "$SCRIPT_DIR/lib/replacer.sh"

generate_from_template() {
  local template_path="$1"
  local entity="$2"
  local output_path="$3"

  local content
  content=$(cat "$template_path")
  content=$(replace_placeholders "$content" "$entity")

  mkdir -p "$(dirname "$output_path")"
  printf '%s\n' "$content" > "$output_path"
  echo "✓ Created: $output_path"
}
