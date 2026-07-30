#!/usr/bin/env bash
# Output-quality evaluation for one skill (lab-go-quality): runevals → grade → benchmark.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"

SKILL=""
PROVIDER="opencode"
MODEL="digitalocean/deepseek-v4-pro"
RUNS=1
TIMEOUT=600
BASELINE="none"
GRADE_TIMEOUT=300

usage() {
  cat <<EOF
Usage: $(basename "$0") --skill <name> [options]

Quality lab (lab-go-quality) for a skill: runevals (with/without skill),
grade assertions, then aggregate benchmark.json.

Options:
  -s, --skill NAME          Skill name (required)
      --provider NAME       CLI provider (default: opencode)
      --model ID            Model id (default: digitalocean/deepseek-v4-pro)
      --runs N              Runs per eval × config (default: 1)
      --timeout SECS        Per-run timeout for runevals (default: 600)
      --baseline MODE       none | snapshot (default: none)
      --grade-timeout SECS  Judge timeout seconds (default: 300)
  -h, --help                Show this help

Paths (fixed convention):
  development/skills/<name>/evals/evals.json
  workspace/quality/<name>/iteration-N/
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    -s|--skill) SKILL="${2:?}"; shift 2 ;;
    --provider) PROVIDER="${2:?}"; shift 2 ;;
    --model) MODEL="${2:?}"; shift 2 ;;
    --runs) RUNS="${2:?}"; shift 2 ;;
    --timeout) TIMEOUT="${2:?}"; shift 2 ;;
    --baseline) BASELINE="${2:?}"; shift 2 ;;
    --grade-timeout) GRADE_TIMEOUT="${2:?}"; shift 2 ;;
    -h|--help) usage; exit 0 ;;
    *) echo "error: unknown argument: $1" >&2; usage >&2; exit 2 ;;
  esac
done

if [[ -z "${SKILL}" ]]; then
  echo "error: --skill is required" >&2
  usage >&2
  exit 2
fi

SKILL_DIR="${REPO_ROOT}/development/skills/${SKILL}"
EVALS_JSON="${SKILL_DIR}/evals/evals.json"
WS="${REPO_ROOT}/workspace/quality/${SKILL}"

die() { echo "error: $*" >&2; exit 1; }
log() { echo "==> $*"; }

[[ -d "${SKILL_DIR}" ]] || die "skill directory not found: ${SKILL_DIR}"
[[ -f "${EVALS_JSON}" ]] || die "evals.json not found: ${EVALS_JSON}"
mkdir -p "${WS}"

run_goq() {
  local pkg="$1"
  shift
  log "go -C lab-go-quality run ${pkg} $*"
  (cd "${REPO_ROOT}" && go -C lab-go-quality run "${pkg}" "$@")
}

latest_iteration() {
  local best=""
  local best_n=0
  local d n
  shopt -s nullglob
  for d in "${WS}"/iteration-*; do
    [[ -d "${d}" ]] || continue
    n="${d##*/iteration-}"
    if [[ "${n}" =~ ^[0-9]+$ ]] && (( n > best_n )); then
      best_n="${n}"
      best="${d}"
    fi
  done
  shopt -u nullglob
  echo "${best}"
}

BEFORE="$(latest_iteration)"

log "runevals provider=${PROVIDER} model=${MODEL} runs=${RUNS} baseline=${BASELINE}"
run_goq ./cmd/runevals \
  -evals "${EVALS_JSON}" \
  -skill-path "${SKILL_DIR}" \
  -workspace "${WS}" \
  -provider "${PROVIDER}" \
  -model "${MODEL}" \
  -baseline "${BASELINE}" \
  -runs "${RUNS}" \
  -timeout "${TIMEOUT}"

ITER="$(latest_iteration)"
[[ -n "${ITER}" ]] || die "no iteration-* directory under ${WS}"
if [[ -n "${BEFORE}" && "${ITER}" == "${BEFORE}" ]]; then
  echo "warning: latest iteration unchanged (${ITER}); grading that directory anyway" >&2
fi

log "grade iteration=${ITER}"
run_goq ./cmd/grade \
  -iteration "${ITER}" \
  -evals "${EVALS_JSON}" \
  -provider "${PROVIDER}" \
  -model "${MODEL}" \
  -timeout "${GRADE_TIMEOUT}"

log "benchmark iteration=${ITER}"
run_goq ./cmd/benchmark -iteration "${ITER}"

echo
echo "======== lab-go-quality summary (${SKILL}) ========"
python3 - "${ITER}/benchmark.json" "${ITER}" <<'PY'
import json, sys
path, it = sys.argv[1], sys.argv[2]
with open(path) as f:
    data = json.load(f)
rs = data.get("run_summary") or {}
for key in ("with_skill", "without_skill", "old_skill"):
    block = rs.get(key)
    if not block:
        continue
    pr = (block.get("pass_rate") or {}).get("mean")
    ts = (block.get("time_seconds") or {}).get("mean")
    tok = (block.get("tokens") or {}).get("mean")
    print(f"{key}: pass_rate={pr} time_s={ts} tokens={tok}")
delta = rs.get("delta") or {}
print(f"delta: pass_rate={delta.get('pass_rate')} time_s={delta.get('time_seconds')} tokens={delta.get('tokens')}")
print(f"iteration: {it}")
print(f"timing: {it}/timing.json")
print(f"benchmark: {path}")
PY
echo "done."
