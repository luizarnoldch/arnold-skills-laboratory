#!/usr/bin/env bash
set -euo pipefail

# check_service <target> <entity_kebab> <entity_pascal> <entity_camel>
#
# Only called for tRPC transport. REST transport (--transport api) has no
# service.ts -- server.sh only generates <entity>.api.ts in that case, per
# assets/api-server.md (a client wrapper, not a hosted route handler).
check_service() {
  local target="$1" entity_kebab="$2" entity_pascal="$3" entity_camel="$4"
  local file="$target/src/features/$entity_kebab/server/$entity_kebab.service.ts"
  local reasons=()

  require_file "$file" "service file missing" || { fail_layer service "${reasons[@]}"; return 1; }

  require_grep "$file" "\* as ${entity_camel}Repository from \"./$entity_kebab.repository\"" "must import * as ${entity_camel}Repository from ./$entity_kebab.repository" || true

  for fn in list get create update remove; do
    require_grep "$file" "export const ${fn} = " "must export ${fn}(...)" || true
  done

  require_grep "$file" '"[A-Za-z]+ not found"' "get/update/remove must throw an Error(\"... not found\") after a null repository lookup" || true
  require_grep "$file" '"[A-Za-z]+ already exists"' "create must throw an Error(\"... already exists\") on duplicate" || true

  if grep -Eq '\b(prisma|db)\.[a-zA-Z]+\.(findMany|findUnique|findFirst|create|update|delete)\b' "$file"; then
    reasons+=("service must delegate to the repository, not call the ORM client directly")
  fi

  if [[ ${#reasons[@]} -eq 0 ]]; then
    pass_layer service
  else
    fail_layer service "${reasons[@]}"
  fi
}