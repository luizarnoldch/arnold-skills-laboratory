# lab-go-quality — laboratorio de calidad de salida

Mide si el skill produce buenas salidas (with_skill vs baseline). La fiabilidad de trigger queda en [lab-go](./lab-go.md).

## Pipeline

```text
                    ┌──────────────────────────────────────────┐
                    │         lab-go-quality pipeline          │
                    └──────────────────────────────────────────┘
         │            │            │            │            │
         ▼            ▼            ▼            ▼            ▼
     scaffold      runevals       grade      benchmark     iterate
         │            │            │            │            │
    evals/ +          │            │            │      feedback.json
    workspace/        │            │            │      + grading/benchmark
    quality/<name>    │            │            │            │
                      ▼            │            │            ▼
              NextIterationDir     │            │     LLM propone
              skill-snapshot/      │            │     proposed-SKILL.md
                      │            │            │     (-apply → SKILL.md)
                      ▼            │            │
         for cada eval × config × run:
           config = with_skill
                  + without_skill  (baseline none)
                  | old_skill      (baseline snapshot)
                      │
                      ▼
              sandbox + copy inputs
              InstallSkill? (sí / no / snapshot)
              TaskRunner.Run (SIN early-kill)
              → outputs/, transcript, duration, tokens
              Append timing.json (index + path)
                      │
                      ▼
              grade(-index N):
                lee outputs + assertions
                LLM judge → grading.json
                      │
                      ▼
              benchmark:
                stack timing + grading
                mean pass_rate / time / tokens
                delta = with_skill - baseline
                → benchmark.json
```

## Una corrida de eval (`cmd/runevals`)

```text
  evals.json
      │
      ├─ eval-slug/
      │     ├─ with_skill/run_001/
      │     │     sandbox/  outputs/  transcript.log  timing.json
      │     └─ without_skill/run_001/   (o old_skill/)
      │
      └─ timing.json (stack global)
            [{index, eval_slug, config, path, duration_ms, tokens}, ...]
```

## Grading (`internal/grade`)

```text
  expected + assertions + file excerpts + transcript
                    │
                    ▼
              JudgePrompt → LLM
                    │
                    ▼
         assertion_results[{text, passed, evidence}]
                    │
                    ▼
              pass_rate = passed/total → grading.json
```

## Relación con lab-go

```text
  SKILL.md (description + body)
       │
       ├──────────────────► lab-go ─────────► ¿CUÁNDO dispara?
       │                    (triggers)         accuracy vs should_trigger
       │
       └──────────────────► lab-go-quality ─► ¿QUÉ TAN BIEN produce?
                            (outputs)          with_skill vs baseline
```
