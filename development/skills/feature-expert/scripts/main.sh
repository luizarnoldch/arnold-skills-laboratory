#!/usr/bin/env bash
set -euo pipefail

ROOT_FEATURES="./FEATURES.yml"
ARCHIVE_DIR="./spect/archive"
ARCHIVE_FEATURES="${ARCHIVE_DIR}/FEATURES.yml"
FEATURES_DIR="./spect/features"
ADRS_DIR="./spect/adrs"

init_workspace() {
    mkdir -p "${ARCHIVE_DIR}"
    mkdir -p "${FEATURES_DIR}"
    mkdir -p "${ADRS_DIR}"

    if [ ! -f "${ROOT_FEATURES}" ]; then
        cat <<EOF > "${ROOT_FEATURES}"
sprint_metadata:
  current_sprint: 1
  sprint_start: "$(date +%Y-%m-%d)"
  sprint_end: "$(date -d "+14 days" +%Y-%m-%d 2>/dev/null || date -v+14d +%Y-%m-%d)"

features: []
EOF
        echo "Initialized ${ROOT_FEATURES}"
    else
        echo "${ROOT_FEATURES} already exists."
    fi

    if [ ! -f "${ARCHIVE_FEATURES}" ]; then
        cat <<EOF > "${ARCHIVE_FEATURES}"
archived_sprints: []
archived_features: []
EOF
        echo "Initialized ${ARCHIVE_FEATURES}"
    else
        echo "${ARCHIVE_FEATURES} already exists."
    fi
}

validate_paths() {
    local prd_path="$1"
    if [ ! -f "${prd_path}" ]; then
        echo "Warning: Target PRD index file explicitly referenced does not exist: ${prd_path}"
        return 1
    fi
    return 0
}

case "${1:-}" in
    init)
        init_workspace
        ;;
    validate-prd)
        if [ -z "${2:-}" ]; then
            echo "Usage: $0 validate-prd <path_to_prd_index>"
            exit 1
        fi
        validate_paths "$2"
        ;;
    *)
        echo "Usage: $0 {init|validate-prd <path>}"
        exit 1
        ;;
esac
