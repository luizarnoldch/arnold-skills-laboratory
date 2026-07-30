# stamp-note (smoke skill)

Minimal skill for exercising both laboratories:

- **lab-go** — trigger detection (`should_trigger`)
- **lab-go-quality** — output quality (with_skill vs without_skill)

Canonical skill: this directory (`SKILL.md` + `evals/`).

## Recommended: scripts

From the repo root:

```bash
# Trigger only
./scripts/lab/lab-go-skill-run.sh \
  --skill stamp-note \
  --provider opencode \
  --model digitalocean/deepseek-v4-pro \
  --runs 1

# Quality only
./scripts/lab/lab-go-quality-skill-run.sh \
  --skill stamp-note \
  --provider opencode \
  --model digitalocean/deepseek-v4-pro \
  --runs 1

# Both (trigger then quality)
./scripts/lab/complete-run-skill.sh \
  --skill stamp-note \
  --provider opencode \
  --model digitalocean/deepseek-v4-pro \
  --runs 1
```

See `./scripts/lab/*.sh --help` for all flags (`--timeout`, `--force-split`, `--baseline`, etc.).

## Manual commands (optional)

### Trigger lab (`lab-go`)

Workspace: `workspace/skills/stamp-note/`

```bash
go -C lab-go run ./cmd/splitprompts \
  -i "$(pwd)/workspace/skills/stamp-note/prompts/prompts.json" \
  -seed 42

go -C lab-go run ./cmd/evaluate \
  -skill-name stamp-note \
  -prompts "$(pwd)/workspace/skills/stamp-note/prompts/train.json" \
  -provider opencode \
  -model digitalocean/deepseek-v4-pro \
  -runs 1 \
  -workdir "$(pwd)/workspace/skills/stamp-note/sandbox" \
  -out "$(pwd)/workspace/skills/stamp-note/results/train/smoke.json"
```

Use **absolute paths** with `go -C` (relative paths resolve under `lab-go/`).

### Quality lab (`lab-go-quality`)

Workspace: `workspace/quality/stamp-note/`

```bash
REPO="$(pwd)"
go -C lab-go-quality run ./cmd/runevals \
  -evals "$REPO/development/skills/stamp-note/evals/evals.json" \
  -skill-path "$REPO/development/skills/stamp-note" \
  -workspace "$REPO/workspace/quality/stamp-note" \
  -provider opencode \
  -model digitalocean/deepseek-v4-pro \
  -runs 1

go -C lab-go-quality run ./cmd/grade \
  -iteration "$REPO/workspace/quality/stamp-note/iteration-1" \
  -evals "$REPO/development/skills/stamp-note/evals/evals.json" \
  -provider opencode \
  -model digitalocean/deepseek-v4-pro

go -C lab-go-quality run ./cmd/benchmark \
  -iteration "$REPO/workspace/quality/stamp-note/iteration-1"
```
