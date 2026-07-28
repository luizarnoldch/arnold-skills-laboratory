#!/usr/bin/env bash
# Full skill eval: lab-go (trigger) then lab-go-quality (output quality).
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"

SKILL=""
PROVIDER="opencode"
MODEL="digitalocean/deepseek-v4-pro"
RUNS=1
TIMEOUT=120
QUALITY_TIMEOUT=600
GRADE_TIMEOUT=300
SEED=42
TAG="$(date +%Y%m%d_%H%M%S)"
BASELINE="none"
FORCE_SPLIT=0
SKIP_VALIDATION=0
SKIP_TRIGGER=0
SKIP_QUALITY=0

usage() {
  cat <<EOF
Usage: $(basename "$0") --skill <name> [options]

Runs both laboratories for one skill:
  1) scripts/lab/lab-go-skill-run.sh         (trigger)
  2) scripts/lab/lab-go-quality-skill-run.sh (quality)

Common options:
  -s, --skill NAME          Skill name (required)
      --provider NAME       CLI provider (default: opencode)
      --model ID            Model id (default: digitalocean/deepseek-v4-pro)
      --runs N              Runs per prompt/eval (default: 1)

Trigger (lab-go) options:
      --timeout SECS        Trigger per-run timeout (default: 120)
      --seed N              Split seed (default: 42)
      --tag TAG             Results filename tag (default: timestamp)
      --force-split         Regenerate train/validation
      --skip-validation     Only evaluate train
      --skip-trigger        Skip lab-go entirely

Quality (lab-go-quality) options:
      --quality-timeout S   runevals timeout (default: 600)
      --grade-timeout S     grade/judge timeout (default: 300)
      --baseline MODE       none | snapshot (default: none)
      --skip-quality        Skip lab-go-quality entirely

  -h, --help                Show this help

Example:
  $(basename "$0") --skill stamp-note --provider opencode \\
    --model digitalocean/deepseek-v4-pro --runs 1
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    -s|--skill) SKILL="${2:?}"; shift 2 ;;
    --provider) PROVIDER="${2:?}"; shift 2 ;;
    --model) MODEL="${2:?}"; shift 2 ;;
    --runs) RUNS="${2:?}"; shift 2 ;;
    --timeout) TIMEOUT="${2:?}"; shift 2 ;;
    --quality-timeout) QUALITY_TIMEOUT="${2:?}"; shift 2 ;;
    --grade-timeout) GRADE_TIMEOUT="${2:?}"; shift 2 ;;
    --seed) SEED="${2:?}"; shift 2 ;;
    --tag) TAG="${2:?}"; shift 2 ;;
    --baseline) BASELINE="${2:?}"; shift 2 ;;
    --force-split) FORCE_SPLIT=1; shift ;;
    --skip-validation) SKIP_VALIDATION=1; shift ;;
    --skip-trigger) SKIP_TRIGGER=1; shift ;;
    --skip-quality) SKIP_QUALITY=1; shift ;;
    -h|--help) usage; exit 0 ;;
    *) echo "error: unknown argument: $1" >&2; usage >&2; exit 2 ;;
  esac
done

if [[ -z "${SKILL}" ]]; then
  echo "error: --skill is required" >&2
  usage >&2
  exit 2
fi

TRIGGER_SH="${SCRIPT_DIR}/lab-go-skill-run.sh"
QUALITY_SH="${SCRIPT_DIR}/lab-go-quality-skill-run.sh"

[[ -x "${TRIGGER_SH}" || -f "${TRIGGER_SH}" ]] || { echo "error: missing ${TRIGGER_SH}" >&2; exit 1; }
[[ -x "${QUALITY_SH}" || -f "${QUALITY_SH}" ]] || { echo "error: missing ${QUALITY_SH}" >&2; exit 1; }

echo "######## complete-run-skill: ${SKILL} ########"
echo "repo=${REPO_ROOT}"
echo "provider=${PROVIDER} model=${MODEL} runs=${RUNS} tag=${TAG}"
echo

if [[ "${SKIP_TRIGGER}" -eq 0 ]]; then
  TRIGGER_ARGS=(
    --skill "${SKILL}"
    --provider "${PROVIDER}"
    --model "${MODEL}"
    --runs "${RUNS}"
    --timeout "${TIMEOUT}"
    --seed "${SEED}"
    --tag "${TAG}"
  )
  [[ "${FORCE_SPLIT}" -eq 1 ]] && TRIGGER_ARGS+=(--force-split)
  [[ "${SKIP_VALIDATION}" -eq 1 ]] && TRIGGER_ARGS+=(--skip-validation)
  bash "${TRIGGER_SH}" "${TRIGGER_ARGS[@]}"
else
  echo "==> skipping trigger (lab-go)"
fi

echo

if [[ "${SKIP_QUALITY}" -eq 0 ]]; then
  QUALITY_ARGS=(
    --skill "${SKILL}"
    --provider "${PROVIDER}"
    --model "${MODEL}"
    --runs "${RUNS}"
    --timeout "${QUALITY_TIMEOUT}"
    --grade-timeout "${GRADE_TIMEOUT}"
    --baseline "${BASELINE}"
  )
  bash "${QUALITY_SH}" "${QUALITY_ARGS[@]}"
else
  echo "==> skipping quality (lab-go-quality)"
fi

echo
echo "======== complete-run-skill finished (${SKILL}) ========"
echo "trigger results: ${REPO_ROOT}/workspace/skills/${SKILL}/results/{train,validation}/${TAG}.json"
echo "quality workspace: ${REPO_ROOT}/workspace/quality/${SKILL}/"
echo "done."
