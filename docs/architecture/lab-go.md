# lab-go — laboratorio de triggers

Mide si el skill se dispara bien (`should_trigger`).

## Pipeline

```text
                    ┌─────────────────────────────────────┐
                    │           lab-go pipeline           │
                    └─────────────────────────────────────┘
                                      │
          ┌───────────────────────────┼───────────────────────────┐
          ▼                           ▼                           ▼
   splitprompts                  evaluate                      optimize
   (cmd/splitprompts)            (cmd/evaluate)                (cmd/optimize)
          │                           │                           │
 prompts.json                ┌────────┴────────┐                  │
      │                      │  eval.Run()     │◄─────────────────┤
      ▼                      └────────┬────────┘                  │
 StratifiedSplit                      │                           │
 60% train / 40% val                  ▼                           │
      │                    for cada prompt:                       │
      ├─► train.json           for run 1..N:                      │
      └─► validation.json         provider.Run(query)             │
                                      │                           │
                                      ▼                           │
                         stream stdout; DetectTrigger             │
                         (kill early si trigger)                  │
                                      │                           │
                                      ▼                           │
                         trigger_rate = triggers/N                │
                         predicted = rate >= 0.5                  │
                         correct = predicted == should_trigger    │
                                      │                           │
                                      ▼                           │
                              results JSON + accuracy             │
                                                                  │
                    ┌─────────────────────────────────────────────┘
                    │  loop max-iters:
                    │    1. eval train  → accuracy + failures
                    │    2. eval validation → elige "best"
                    │    3. si accuracy >= threshold → STOP
                    │    4. else: LLM reescribe description
                    │       (solo frontmatter; clamp 1024)
                    │    5. al salir: restore best description
                    │       (gana validation accuracy)
                    └─────────────────────────────────────────────
```

## Núcleo por prompt

`internal/eval` + `provider`:

```text
  query ──► CLI (opencode/...) ──► stdout stream
                                        │
                          ¿aparece skill en output?
                           /                    \
                         SÍ                     NO
                          │                      │
                     TRIGGERED              timeout/fin
                     (kill proc)                 │
                          \                     /
                           \                   /
                            trigger_count++?
                                   │
                    rate = triggers/runs
                    correct <=> (rate>=0.5)==should_trigger
```

## Relación con lab-go-quality

```text
  SKILL.md (description + body)
       │
       ├──────────────────► lab-go ─────────► ¿CUÁNDO dispara?
       │                    (triggers)         accuracy vs should_trigger
       │
       └──────────────────► lab-go-quality ─► ¿QUÉ TAN BIEN produce?
                            (outputs)          with_skill vs baseline
```

Ver también [lab-go-quality.md](./lab-go-quality.md).
