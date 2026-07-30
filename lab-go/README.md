# Skill trigger laboratory (Go)

Go port of [`lab/`](../lab/). Same JSON contracts, same workspace layout, stdlib-only module `skills-laboratory/lab-go`.

## Layout

```text
lab-go/
├── schema/prompts.schema.json
├── cmd/
│   ├── splitprompts/   # 60/40 stratified split
│   ├── evaluate/       # prompt set → triggers / trigger_rate / accuracy
│   └── optimize/       # train rewrite + validation winner selection
└── internal/
    ├── prompt/
    ├── result/
    ├── eval/
    ├── logdir/
    ├── skillmd/
    └── provider/       # opencode + stubs (codex, claude, cursor_agent)
```

Per-skill data stays in `workspace/skills/<name>/` (shared with the Python lab).

## Prompt contract

```json
{
  "id": 1,
  "query": "Initialize the FEATURES.yml file for this project",
  "should_trigger": true
}
```

Result fields: `skill_name`, `triggers`, `runs`, `trigger_rate`, `correct` where
`correct = (trigger_rate >= 0.5) == should_trigger`.

## Flow

### 1. Create the skill

Author `development/skills/<name>/SKILL.md` (`name` + `description` frontmatter). Train only rewrites `description`.

### 2. Split prompts

```bash
go -C lab-go run ./cmd/splitprompts \
  -i ../workspace/skills/<name>/prompts/prompts.json \
  -seed 42
```

### 3. Evaluate

```bash
go -C lab-go run ./cmd/evaluate \
  -skill-name feature-expert \
  -prompts ../workspace/skills/feature-expert/prompts/train.json \
  -provider opencode \
  -runs 3 \
  -model digitalocean/deepseek-4-flash \
  -workdir ../workspace/skills/feature-expert/sandbox \
  -out ../workspace/skills/feature-expert/results/train/go_run_001.json
```

Providers: `opencode` (ready), `codex`, `claude`, `cursor_agent` / `agent` (stubs).

### 4. Train — optimize description

Train prompts guide rewrites. Validation prompts (sibling `validation.json` by default, or `-validation-prompts`) score each iteration; the loop restores the description with the best validation accuracy (Agent Skills anti-overfitting rule). Descriptions are clamped to 1024 characters.

```bash
go -C lab-go run ./cmd/optimize \
  -skill-name feature-expert \
  -skill-md ../development/skills/feature-expert/SKILL.md \
  -prompts ../workspace/skills/feature-expert/prompts/train.json \
  -validation-prompts ../workspace/skills/feature-expert/prompts/validation.json \
  -workdir ../workspace/skills/feature-expert/sandbox \
  -iterations-dir ../workspace/skills/feature-expert/iterations \
  -results-dir ../workspace/skills/feature-expert/results/train \
  -provider opencode \
  -runs 3 \
  -threshold 0.95 \
  -max-iters 5
```

Do not pass `validation.json` as `-prompts` (that would leak hold-out into the rewrite). Use `-validation-prompts` instead.

### 5. Validation — final accuracy (optional sanity)

After optimize applies the validation winner, you can still run a standalone evaluation:

```bash
go -C lab-go run ./cmd/evaluate \
  -skill-name feature-expert \
  -prompts ../workspace/skills/feature-expert/prompts/validation.json \
  -provider opencode \
  -runs 3 \
  -workdir ../workspace/skills/feature-expert/sandbox \
  -out ../workspace/skills/feature-expert/results/validation/go_run_001.json
```

Target: accuracy near 100% after train converges.

## Build / test

```bash
cd lab-go && go test ./... && go build ./...
```

From repo root:

```bash
go -C lab-go test ./...
go -C lab-go build -o /tmp/lab-evaluate ./cmd/evaluate
```

Paths in flags are relative to the current working directory (not `lab-go/`), so prefer absolute paths or run from the repo root with paths like `workspace/skills/...` after `cd` into the repo and invoking binaries built with `-o`.