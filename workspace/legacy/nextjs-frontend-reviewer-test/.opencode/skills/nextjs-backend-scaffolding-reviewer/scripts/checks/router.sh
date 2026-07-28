#!/usr/bin/env bash
set -euo pipefail

# check_router <target> <entity_kebab> <entity_pascal> <entity_camel> <transport>
#
# tRPC  -> server/<entity>.router.ts, per assets/trpc-router.md
# REST  -> server/<entity>.api.ts, a CLIENT wrapper around an external API via
#          apiFetch (assets/api-server.md + scripts/server.sh:generate_api_file).
#          It is NOT a Next.js route handler (no GET/POST exports, no Response.json) --
#          the app calls an existing external REST API, it does not host one.
check_router() {
  local target="$1" entity_kebab="$2" entity_pascal="$3" entity_camel="$4" transport="$5"
  local server_dir="$target/src/features/$entity_kebab/server"
  local reasons=()

  if [[ "$transport" == "trpc" ]]; then
    local file="$server_dir/$entity_kebab.router.ts"
    require_file "$file" "router file missing" || { fail_layer router "${reasons[@]}"; return 1; }

    require_grep "$file" 'createTRPCRouter.*from "@/trpc/init"' "must import createTRPCRouter from @/trpc/init" || true
    require_grep "$file" 'protectedProcedure' "must import/use protectedProcedure" || true
    require_grep "$file" 'baseProcedure' "must import/use baseProcedure" || true
    require_grep "$file" 'TRPCError.*from "@trpc/server"' "must import TRPCError from @trpc/server" || true
    require_grep "$file" "\* as ${entity_camel}Service from \"./$entity_kebab.service\"" "must import * as ${entity_camel}Service from ./$entity_kebab.service" || true

    for route in list get create update delete; do
      require_grep "$file" "^[[:space:]]*${route}:" "must define the '${route}' route" || true
    done

    require_grep "$file" "${entity_camel}Service\.list\(\)" "list route must call ${entity_camel}Service.list()" || true
    require_grep "$file" "${entity_camel}Service\.create\(" "create route must call ${entity_camel}Service.create(...)" || true
    require_grep "$file" "${entity_camel}Service\.update\(" "update route must call ${entity_camel}Service.update(...)" || true
    require_grep "$file" "${entity_camel}Service\.remove\(" "delete route must call ${entity_camel}Service.remove(...)" || true

    require_grep "$file" 'code: "NOT_FOUND"' "must throw TRPCError with code NOT_FOUND for not-found errors" || true
    require_grep "$file" 'code: "CONFLICT"' "must throw TRPCError with code CONFLICT for duplicate errors" || true
    require_grep "$file" 'code: "INTERNAL_SERVER_ERROR"' "must throw TRPCError with code INTERNAL_SERVER_ERROR as the fallback" || true

  elif [[ "$transport" == "api" ]]; then
    local file="$server_dir/$entity_kebab.api.ts"
    require_file "$file" "api client file missing" || { fail_layer router "${reasons[@]}"; return 1; }

    require_grep "$file" 'apiFetch.*from "@/lib/api"' "must import apiFetch from @/lib/api" || true
    require_grep "$file" "list${entity_pascal}s" "must export list${entity_pascal}s()" || true
    require_grep "$file" "get${entity_pascal}" "must export get${entity_pascal}(id)" || true
    require_grep "$file" "create${entity_pascal}" "must export create${entity_pascal}(data)" || true
    require_grep "$file" "update${entity_pascal}" "must export update${entity_pascal}(...)" || true
    require_grep "$file" "delete${entity_pascal}" "must export delete${entity_pascal}(id)" || true

    if grep -Eq '^\s*export\s+(async\s+)?function\s+(GET|POST|PUT|PATCH|DELETE)\b' "$file"; then
      reasons+=("looks like a Next.js route handler (GET/POST exports) -- this project's REST transport is a client wrapper calling apiFetch, not a hosted route handler")
    fi

  else
    reasons+=("unknown transport '$transport' (expected trpc or api)")
  fi

  if [[ ${#reasons[@]} -eq 0 ]]; then
    pass_layer router
  else
    fail_layer router "${reasons[@]}"
  fi
}