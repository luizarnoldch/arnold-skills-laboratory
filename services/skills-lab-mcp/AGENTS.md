# AGENTS.md — skills-lab-mcp

Servidor MCP stdio (Go) — fachada para LLMs externos sobre lab-api y orchestrator.

## Purpose

Expone 9 tools MCP para listar/consultar skills, gestionar test sets, encolar evals y optimización. No persiste datos ni ejecuta LLMs.

**Boundaries:** HTTP client only → `arnold-laboratory-api` + `arnold-lab-orchestrator`. Sin DB, sin puerto HTTP propio.

## Layout

| Path | Role |
|------|------|
| `cmd/skills_lab_mcp` | Entrypoint stdio |
| `pkg/config` | Env (`LAB_API_URL`, `ORCH_API_URL`, …) |
| `pkg/client/labapi` | Cliente laboratory-api |
| `pkg/client/orch` | Cliente orchestrator |
| `pkg/resolve` | skill_id/name, description/content defaults |
| `pkg/tools` | Lógica por dominio |
| `pkg/mcp` | Registro de tools con go-sdk |
| `Dockerfile` | Imagen multi-stage; CMD `sleep infinity` en runtime |

## Environment

Defaults local: `LAB_API_URL=http://127.0.0.1:18180`, `ORCH_API_URL=http://127.0.0.1:18182`, `WORKSPACE_DEFAULT=.`, `HTTP_TIMEOUT=30s`.

En Docker (compose): URLs internas `http://lab-api:18180`, `http://orch:18182`, `WORKSPACE_DEFAULT=/workspace`. El MCP se ejecuta con `docker compose exec -i skills-lab-mcp /app/skills_lab_mcp` (ver [`deploy/docker/mcp-exec.sh`](../../deploy/docker/mcp-exec.sh)).

Ver `.env.example`. Logs solo en **stderr** (stdout reservado a JSON-RPC MCP).

## Commands

```bash
make build   # bin/skills_lab_mcp
make run
make test
```

Docker (desde raíz del monorepo):

```bash
docker compose build skills-lab-mcp
docker compose up -d skills-lab-mcp
```

## Do not

- No escribir a stdout excepto el transporte MCP
- No importar módulos internos de otros servicios Go del monorepo
- No duplicar lógica de eval/optimize — delegar al orchestrator
