#!/usr/bin/env bash
set -euo pipefail

# check_repository <target> <entity_kebab> <entity_pascal> <entity_camel> <database>
#
# Only called for tRPC transport (see checks/service.sh for why REST has none).
#
# IMPORTANT: both assets/prisma-repository.md and assets/drizzle-repository.md
# import the client as `db` from "@/lib/db" -- there is no `prisma` import from
# "@/lib/prisma" anywhere in this project's templates. A prior hand-edit of
# SKILL.md briefly required `prisma`/"@/lib/prisma" here; that would fail every
# correctly-scaffolded repository (verified against the real generated fixture
# at workspaces/eval-1-exploring/outputs/.../product.repository.ts, which uses
# `db`). Keep this check pinned to `db`.
check_repository() {
  local target="$1" entity_kebab="$2" entity_pascal="$3" entity_camel="$4" database="$5"
  local file="$target/src/features/$entity_kebab/server/$entity_kebab.repository.ts"
  local reasons=()

  require_file "$file" "repository file missing" || { fail_layer repository "${reasons[@]}"; return 1; }

  require_grep "$file" 'import \{ db \} from "@/lib/db"' 'must import { db } from "@/lib/db" (not `prisma` -- see script header note)' || true

  for fn in findAll findById findByName create update remove; do
    require_grep "$file" "export const ${fn} = " "must export ${fn}(...)" || true
  done

  if [[ "$database" == "drizzle" ]]; then
    require_grep "$file" 'from "drizzle-orm"' "drizzle repository must import eq/desc from drizzle-orm" || true
    require_grep "$file" 'db\.select\(\)\.from\(' "findAll must use db.select().from(...) (drizzle query builder)" || true
    require_grep "$file" 'db\.insert\(' "create must use db.insert(...) (drizzle query builder)" || true
  else
    require_grep "$file" "db\.${entity_camel}\.findMany\(" "findAll must call db.${entity_camel}.findMany(...) (prisma delegate)" || true
    require_grep "$file" "db\.${entity_camel}\.create\(" "create must call db.${entity_camel}.create(...) (prisma delegate)" || true
  fi

  if grep -Eq '\bif\b|\bthrow\b' "$file"; then
    reasons+=("repository must be pure ORM delegation -- no business logic (if/throw found)")
  fi

  if [[ ${#reasons[@]} -eq 0 ]]; then
    pass_layer repository
  else
    fail_layer repository "${reasons[@]}"
  fi
}