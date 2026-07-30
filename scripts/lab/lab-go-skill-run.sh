#!/usr/bin/env bash
# Trigger evaluation for one skill (lab-go): split (if needed) → evaluate train + validation.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"

SKILL=""
PROVIDER="opencode"
MODEL="digitalocean/deepseek-v4-pro"
RUNS=1
TIMEOUT=120
SEED=42
TAG="$(date +%Y%m%d_%H%M%S)"
FORCE_SPLIT=0
SKIP_VALIDATION=0

usage() {
  cat <<EOF
Usage: $(basename "$0") --skill <name> [options]

Trigger lab (lab-go) for a skill: stratified split if needed, then evaluate
train.json and validation.json.

Options:
  -s, --skill NAME          Skill name (required)
      --provider NAME       CLI provider (default: opencode)
      --model ID            Model id (default: digitalocean/deepseek-v4-pro)
      --runs N              Runs per prompt (default: 1)
      --timeout SECS        Per-run timeout seconds (default: 120)
      --seed N              Split RNG seed (default: 42)
      --tag TAG             Results filename tag (default: timestamp)
      --force-split         Re-run splitprompts even if train/validation exist
      --skip-validation     Only evaluate train.json
  -h, --help                Show this help

Paths (fixed convention):
  workspace/skills/<name>/{prompts,sandbox,logs,results}
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    -s|--skill) SKILL="${2:?}"; shift 2 ;;
    --provider) PROVIDER="${2:?}"; shift 2 ;;
    --model) MODEL="${2:?}"; shift 2 ;;
    --runs) RUNS="${2:?}"; shift 2 ;;
    --timeout) TIMEOUT="${2:?}"; shift 2 ;;
    --seed) SEED="${2:?}"; shift 2 ;;
    --tag) TAG="${2:?}"; shift 2 ;;
    --force-split) FORCE_SPLIT=1; shift ;;
    --skip-validation) SKIP_VALIDATION=1; shift ;;
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
WS="${REPO_ROOT}/workspace/skills/${SKILL}"
PROMPTS_DIR="${WS}/prompts"
PROMPTS_JSON="${PROMPTS_DIR}/prompts.json"
TRAIN_JSON="${PROMPTS_DIR}/train.json"
VAL_JSON="${PROMPTS_DIR}/validation.json"
SANDBOX="${WS}/sandbox"
LOG_DIR="${WS}/logs"
TRAIN_OUT_DIR="${WS}/results/train"
VAL_OUT_DIR="${WS}/results/validation"

die() { echo "error: $*" >&2; exit 1; }

[[ -d "${SKILL_DIR}" ]] || die "skill directory not found: ${SKILL_DIR}"
[[ -d "${WS}" ]] || die "trigger workspace not found: ${WS}"
[[ -f "${PROMPTS_JSON}" ]] || die "prompts.json not found: ${PROMPTS_JSON}"
[[ -d "${SANDBOX}" ]] || die "sandbox not found: ${SANDBOX}"

mkdir -p "${LOG_DIR}" "${TRAIN_OUT_DIR}" "${VAL_OUT_DIR}"

log() { echo "==> $*"; }

run_go() {
  local pkg="$1"
  shift
  log "go -C lab-go run ${pkg} $*"
  (cd "${REPO_ROOT}" && go -C lab-go run "${pkg}" "$@")
}

print_summary() {
  local label="$1"
  local path="$2"
  if [[ ! -f "${path}" ]]; then
    echo "${label}: (missing ${path})"
    return
  fi
  python3 - "${path}" "${label}" <<'PY'
import json, sys
path, label = sys.argv[1], sys.argv[2]
with open(path) as f:
    data = json.load(f)
s = data.get("summary") or {}
acc = s.get("accuracy_pct", s.get("accuracy"))
correct = s.get("correct")
total = s.get("total")
print(f"{label}: accuracy={acc} correct={correct}/{total}  ({path})")
PY
}

# --- split ---
if [[ "${FORCE_SPLIT}" -eq 1 || ! -f "${TRAIN_JSON}" || ! -f "${VAL_JSON}" ]]; then
  log "splitting prompts (seed=${SEED})"
  run_go ./cmd/splitprompts -i "${PROMPTS_JSON}" -seed "${SEED}"
else
  log "using existing train.json / validation.json (pass --force-split to regenerate)"
fi

[[ -f "${TRAIN_JSON}" ]] || die "train.json missing after split: ${TRAIN_JSON}"
[[ "${SKIP_VALIDATION}" -eq 1 || -f "${VAL_JSON}" ]] || die "validation.json missing: ${VAL_JSON}"

TRAIN_OUT="${TRAIN_OUT_DIR}/${TAG}.json"
VAL_OUT="${VAL_OUT_DIR}/${TAG}.json"

# --- train ---
log "evaluate TRAIN provider=${PROVIDER} model=${MODEL} runs=${RUNS}"
run_go ./cmd/evaluate \
  -skill-name "${SKILL}" \
  -prompts "${TRAIN_JSON}" \
  -provider "${PROVIDER}" \
  -model "${MODEL}" \
  -runs "${RUNS}" \
  -timeout "${TIMEOUT}" \
  -workdir "${SANDBOX}" \
  -log-dir "${LOG_DIR}" \
  -out "${TRAIN_OUT}"

# --- validation ---
if [[ "${SKIP_VALIDATION}" -eq 0 ]]; then
  log "evaluate VALIDATION provider=${PROVIDER} model=${MODEL} runs=${RUNS}"
  run_go ./cmd/evaluate \
    -skill-name "${SKILL}" \
    -prompts "${VAL_JSON}" \
    -provider "${PROVIDER}" \
    -model "${MODEL}" \
    -runs "${RUNS}" \
    -timeout "${TIMEOUT}" \
    -workdir "${SANDBOX}" \
    -log-dir "${LOG_DIR}" \
    -out "${VAL_OUT}"
fi

echo
echo "======== lab-go summary (${SKILL}) ========"
print_summary "train" "${TRAIN_OUT}"
if [[ "${SKIP_VALIDATION}" -eq 0 ]]; then
  print_summary "validation" "${VAL_OUT}"
fi
echo "tag=${TAG}"
echo "done."
