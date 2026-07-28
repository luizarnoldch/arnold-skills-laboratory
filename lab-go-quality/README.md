# Skill output quality laboratory (Go)

Eval-driven iteration for skill **output quality** (with_skill vs without_skill), following [agentskills.io evaluating-skills](https://agentskills.io/skill-creation/evaluating-skills.md).

Trigger reliability stays in [`lab-go/`](../lab-go/). This module measures whether the skill produces good outputs.

## Layout

```text
lab-go-quality/
├── schema/                 # evals, timing, grading, benchmark, feedback
├── cmd/
│   ├── scaffold/           # create evals/ + workspace/quality/<name>/
│   ├── runevals/           # run with_skill + baseline → timing stack
│   ├── grade/              # LLM judge → grading.json per run (by index)
│   ├── benchmark/          # aggregate → benchmark.json
│   └── iterate/            # propose SKILL.md improvements
└── internal/
    ├── evalset/
    ├── workspace/
    ├── timing/             # stacked indexed timing.json
    ├── runner/             # opencode, claude, codex, agent
    ├── grade/
    ├── benchmark/
    └── skillmd/

development/skills/<name>/
├── SKILL.md
└── evals/
    ├── evals.json
    └── files/

workspace/quality/<name>/
└── iteration-N/
    ├── timing.json           # stacked runs with index + path
    ├── benchmark.json
    ├── feedback.json         # human notes per eval-slug
    ├── skill-snapshot/
    └── eval-<slug>/
        ├── with_skill/run_001/{outputs/,timing.json,grading.json,transcript.log,sandbox/}
        └── without_skill/run_001/...
```

## Timing stack

Each execution appends to `iteration-N/timing.json` with a monotonic `index` and a `path` pointing at the run directory — so you can jump straight to the transcript/outputs to fix issues:

```json
{
  "runs": [
    {
      "index": 1,
      "eval_id": 1,
      "eval_slug": "eval-init-features-yml",
      "config": "with_skill",
      "run": 1,
      "path": "eval-init-features-yml/with_skill/run_001",
      "total_tokens": 0,
      "duration_ms": 23332,
      "provider": "opencode"
    }
  ]
}
```

Grade a single execution with `-index N`.

## Providers

| Name | CLI |
|------|-----|
| `opencode` | `opencode run --model <model> <query>` |
| `claude` | `claude -p <query> --model <model>` |
| `codex` | `codex exec --model <model> <query>` |
| `agent` / `cursor_agent` | `agent -p <query> --model <model>` |

Skills are installed per sandbox under `.opencode/skills/`, `.claude/skills/`, `.codex/skills/`, or `.cursor/skills/` respectively. Runs do **not** early-kill on trigger detection.

## Flow

### 1. Scaffold (optional)

```bash
go -C lab-go-quality run ./cmd/scaffold \
  -skill-name feature-expert \
  -skill-dir development/skills/feature-expert \
  -workspace-root workspace/quality
```

`feature-expert` already has `evals/evals.json` and sample input files.

### 2. Run evals

```bash
go -C lab-go-quality run ./cmd/runevals \
  -evals development/skills/feature-expert/evals/evals.json \
  -skill-path development/skills/feature-expert \
  -workspace workspace/quality/feature-expert \
  -provider opencode \
  -model digitalocean/deepseek-4-flash \
  -baseline none \
  -runs 1 \
  -timeout 600
```

`-baseline snapshot` compares against `skill-snapshot/` and writes `old_skill/` instead of `without_skill/`.

### 3. Grade

```bash
go -C lab-go-quality run ./cmd/grade \
  -iteration workspace/quality/feature-expert/iteration-1 \
  -evals development/skills/feature-expert/evals/evals.json \
  -provider opencode \
  -model digitalocean/deepseek-4-flash
```

### 4. Benchmark

```bash
go -C lab-go-quality run ./cmd/benchmark \
  -iteration workspace/quality/feature-expert/iteration-1
```

### 5. Human feedback + iterate

Edit `feedback.json` (empty string = OK). Then:

```bash
go -C lab-go-quality run ./cmd/iterate \
  -iteration workspace/quality/feature-expert/iteration-1 \
  -skill-md development/skills/feature-expert/SKILL.md \
  -provider opencode
```

Writes `proposed-SKILL.md`. Add `-apply` to overwrite the canonical skill.

## Build / test

```bash
go -C lab-go-quality test ./...
go -C lab-go-quality build ./...
```
