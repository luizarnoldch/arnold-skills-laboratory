# Skill trigger laboratory

Reusable tooling to create skills, measure whether LLM CLIs invoke them from prompts, optimize the skill `description` on train data, and report final accuracy on held-out validation data.

## Layout

```text
lab/
├── schema/prompts.schema.json
├── bin/
│   ├── split_prompts.py          # 60/40 stratified split
│   ├── evaluate.py               # run prompt set → triggers / trigger_rate / accuracy
│   └── optimize_description.py   # train loop (never use validation here)
└── providers/
    ├── opencode.py               # implemented
    ├── codex.py                  # stub
    ├── claude.py                 # stub
    └── cursor_agent.py           # stub (CLI `agent`)

workspace/skills/<skill-name>/
├── skill -> ../../../development/skills/<skill-name>
├── sandbox/                      # CLI cwd (.opencode/skills/...)
├── prompts/{prompts,train,validation}.json
├── results/{train,validation}/
├── logs/
└── iterations/
```

## Prompt contract

Each item in `prompts.json` (typically 20 = 10 positive / 10 negative):

```json
{
  "id": 1,
  "query": "Initialize the FEATURES.yml file for this project",
  "should_trigger": true
}
```

`id` is the stable index for logs and results. `should_trigger` is the ground truth for whether the skill must be called.

Result fields written by `evaluate.py`:

| Field | Meaning |
|-------|---------|
| `skill_name` | Skill under test |
| `triggers` | Times the skill fired across runs |
| `runs` | Runs for that prompt |
| `trigger_rate` | `triggers / runs` |
| `correct` | `(trigger_rate >= 0.5) == should_trigger` |

## Flow

### 1. Create the skill

Author `development/skills/<name>/SKILL.md` with YAML frontmatter (`name`, `description`). The train loop only rewrites `description`.

### 2. Create prompts and split

```bash
# edit workspace/skills/<name>/prompts/prompts.json (20 items)

python3 lab/bin/split_prompts.py \
  -i workspace/skills/<name>/prompts/prompts.json \
  --seed 42
```

Produces `train.json` (60%) and `validation.json` (40%), stratified by `should_trigger`.

### 3. Evaluate (any set)

```bash
python3 lab/bin/evaluate.py \
  --skill-name feature-expert \
  --prompts workspace/skills/feature-expert/prompts/train.json \
  --provider opencode \
  --runs 3 \
  --model digitalocean/deepseek-4-flash \
  --workdir workspace/skills/feature-expert/sandbox \
  --out workspace/skills/feature-expert/results/train/run_001.json
```

Providers: `opencode` (ready), `codex`, `claude`, `cursor_agent` / `agent` (stubs).

### 4. Train — optimize description from train metrics

```bash
python3 lab/bin/optimize_description.py \
  --skill-name feature-expert \
  --skill-md development/skills/feature-expert/SKILL.md \
  --prompts workspace/skills/feature-expert/prompts/train.json \
  --workdir workspace/skills/feature-expert/sandbox \
  --iterations-dir workspace/skills/feature-expert/iterations \
  --results-dir workspace/skills/feature-expert/results/train \
  --provider opencode \
  --runs 3 \
  --threshold 0.95 \
  --max-iters 5
```

Behavior:

1. Runs `evaluate.py` on **train** only.
2. If accuracy < `--threshold`, asks the LLM for a new `description`, snapshots under `iterations/NNN/`, updates `SKILL.md`.
3. Stops when accuracy ≥ threshold, no failures, `--dry-run`, or `--max-iters`.

**Never pass `validation.json` to this script.**

### 5. Validation — final accuracy only

```bash
python3 lab/bin/evaluate.py \
  --skill-name feature-expert \
  --prompts workspace/skills/feature-expert/prompts/validation.json \
  --provider opencode \
  --runs 3 \
  --model digitalocean/deepseek-4-flash \
  --workdir workspace/skills/feature-expert/sandbox \
  --out workspace/skills/feature-expert/results/validation/run_001.json
```

Target: accuracy near **100%** on the held-out set after train converges.

## Example skill

`feature-expert` is wired under [`workspace/skills/feature-expert/`](../workspace/skills/feature-expert/) with sandbox symlink into OpenCode skills.

Legacy scaffolding / old triggert iterations live in [`workspace/legacy/`](../workspace/legacy/).
