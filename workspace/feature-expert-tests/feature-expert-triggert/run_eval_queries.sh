#!/bin/bash

# Default values
SKILL_NAME="feature-expert"
QUERIES_FILE="/home/arnold/Workspace/skills-laboratory/workspace/feature-expert-tests/.opencode/skills/feature-expert/evals/eval_queries.json"
MODEL="digitalocean/deepseek-4-flash"
RUNS=1

# Helper function to print usage error
usage_error() {
  echo "Error: $1" >&2
  echo "Usage: $0 -s|--skill-name <skill_name> -f|--file /path/to/eval_queries.json [-m|--model <model_id>] [-r|--runs <runs>]" >&2
  exit 1
}

# Parse long and short options
PARSED_ARGS=$(getopt -o s:f:m:r: --long skill-name:,file:,model:,runs: -n "$0" -- "$@")
if [ $? -ne 0 ]; then
  usage_error "Invalid arguments provided."
fi

eval set -- "$PARSED_ARGS"

while true; do
  case "$1" in
    -s|--skill-name)
      SKILL_NAME="$2"
      shift 2
      ;;
    -f|--file)
      QUERIES_FILE="$2"
      shift 2
      ;;
    -m|--model)
      MODEL="$2"
      shift 2
      ;;
    -r|--runs)
      RUNS="$2"
      shift 2
      ;;
    --)
      shift
      break
      ;;
    *)
      usage_error "Unexpected option: $1"
      ;;
  esac
done

# Validate required flags, numeric values, and file presence
if [[ -z "$SKILL_NAME" ]]; then
  usage_error "The '-s|--skill-name' argument is required."
fi

if [[ -z "$QUERIES_FILE" ]]; then
  usage_error "The '-f|--file' argument is required."
fi

if [[ ! -f "$QUERIES_FILE" ]]; then
  usage_error "File '$QUERIES_FILE' was not found."
fi

if [[ ! "$RUNS" =~ ^[1-9][0-9]*$ ]]; then
  usage_error "The '-r|--runs' argument must be a positive integer."
fi

# Checks if OpenCode invoked the specified Skill tool for a given query.
check_triggered() {
  local query="$1"

  # Stream output through jq using recursive search (..), ensuring both NDJSON
  # streams and full payload arrays are properly traversed.
  opencode run "$query" -m "$MODEL" --format json 2>/dev/null \
    | jq -e --arg skill "$SKILL_NAME" '
        [
          ..
          | objects
          | select(
              # Match generic tool call objects targeting "skill"
              (
                (.type? == "tool_use" or .type? == "tool_call" or .type? == "call" or .function? != null) and
                (.name? == "skill" or .name? == "Skill" or .tool? == "skill" or .function?.name? == "skill") and
                (
                  .input?.skill == $skill or
                  .args?.skill == $skill or
                  .parameters?.skill == $skill or
                  .function?.arguments?.skill == $skill
                )
              )
              or
              # Match direct skill invocation names (e.g., tool name is directly the skill name)
              (
                .name? == $skill or .tool? == $skill or .function?.name? == $skill
              )
            )
        ] | length > 0
      ' > /dev/null 2>&1
}

# Setup output directory and filename
RESULTS_DIR="results"
mkdir -p "$RESULTS_DIR"

TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
OUTPUT_FILE="${RESULTS_DIR}/eval_${SKILL_NAME}_${TIMESTAMP}.json"

count=$(jq length "$QUERIES_FILE")

eval_results=$(for i in $(seq 0 $((count - 1))); do
  query=$(jq -r ".[$i].query" "$QUERIES_FILE")
  should_trigger=$(jq -r ".[$i].should_trigger" "$QUERIES_FILE")
  triggers=0

  for run in $(seq 1 $RUNS); do
    check_triggered "$query" && triggers=$((triggers + 1))
  done

  jq -n \
    --arg query "$query" \
    --argjson should_trigger "$should_trigger" \
    --argjson triggers "$triggers" \
    --argjson runs "$RUNS" \
    '{
      query: $query,
      should_trigger: $should_trigger,
      triggers: $triggers,
      runs: $runs,
      trigger_rate: ($triggers / $runs)
    }'
done | jq -s '.')

# Save to results folder and display on stdout
echo "$eval_results" | tee "$OUTPUT_FILE"
echo -e "\n[+] Evaluation completed. Results saved to: $OUTPUT_FILE" >&2
