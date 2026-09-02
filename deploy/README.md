# Docker — Skills Laboratory

Stack containerizado del laboratorio de skills: **lab-api**, **chavez**, **orchestrator**, **web-ui**, **skills-lab-mcp** y **nginx** (entrada única en `:18321`).

## Requisitos

- Docker Engine 24+ y Docker Compose v2
- Submodules inicializados: `git submodule update --init --recursive`
- API key LLM (p. ej. `DEEPSEEK_API_KEY`) para evals reales

## Arranque rápido

```bash
cp .env.docker.example .env
# Editar .env y definir DEEPSEEK_API_KEY

docker compose --env-file .env up -d --build
```

UI: `http://<host>:18321/`

## Puertos

| Servicio | Puerto interno | Host |
|----------|----------------|------|
| nginx (edge) | 18321 | **18321** |
| lab-api | 18180 | — |
| chavez | 18181 | — |
| orchestrator | 18182 | — |
| web-ui | 4321 | — |
| skills-lab-mcp | — (stdio MCP) | — |

## Volúmenes

| Volumen | Contenedor | Uso |
|---------|------------|-----|
| `lab-api-data` | lab-api `/data` | SQLite `lab-skills.db` |
| `orch-data` | orch `/data` | SQLite `lab-orchestrator.db` |
| `chavez-data` | chavez `/data` | SQLite `chavez.db` |
| bind `WORKSPACE_HOST_PATH` | chavez `/workspace` | Sandboxes de eval (default `./workspace`) |

Las migraciones goose se aplican automáticamente al arrancar cada API Go (`db_tool up` en entrypoint).

## Verificación

```bash
# Estado de servicios
docker compose ps

# Readiness interno
docker compose exec lab-api wget -qO- http://127.0.0.1:18180/ready
docker compose exec orch wget -qO- http://127.0.0.1:18182/ready
docker compose exec chavez wget -qO- http://127.0.0.1:18181/health
docker compose exec skills-lab-mcp test -x /app/skills_lab_mcp

# UI vía nginx
curl -s -o /dev/null -w "%{http_code}\n" http://127.0.0.1:18321/

# Rutas orchestrator (deben devolver 200, no 404)
curl -s -o /dev/null -w "%{http_code}\n" http://127.0.0.1:18321/orch/api/v1/baseline-eval-jobs
curl -s -o /dev/null -w "%{http_code}\n" http://127.0.0.1:18321/orch/api/v1/trigger-eval-batches
curl -s -o /dev/null -w "%{http_code}\n" http://127.0.0.1:18321/orch/api/v1/optimize-jobs
```

### WebSocket

Las páginas de progreso (`/trigger-batches/[id]`, `/optimize-jobs/[id]`, `/baseline-evals/[id]`) usan WebSocket en `/orch/ws/*`. Nginx reenvía upgrades; en dev local Vite hace lo mismo. Confirmar en DevTools → Network → WS que la conexión a `ws://<host>:18321/orch/ws/...` queda en estado **101 Switching Protocols**.

### Smoke E2E (manual)

1. Crear skill en `/skills/new`
2. Añadir trigger queries y ejecutar eval global en `/evaluations`
3. Verificar métricas y progreso WS en detalle de batch/baseline

Requiere `DEEPSEEK_API_KEY` válida en `.env`.

## MCP (Cursor)

El servicio **skills-lab-mcp** usa transporte stdio (no HTTP). El contenedor permanece activo con `sleep infinity`; Cursor invoca el MCP vía `docker compose exec`.

Requisito: stack levantado (`docker compose up -d`).

### Cursor (`~/.cursor/mcp.json`)

```json
{
  "mcpServers": {
    "skills-lab": {
      "command": "/ruta/absoluta/al/repo/deploy/docker/mcp-exec.sh"
    }
  }
}
```

Alternativa manual:

```bash
docker compose exec -T -i skills-lab-mcp /app/skills_lab_mcp
```

En Docker, el MCP habla con `lab-api` y `orch` por la red interna (`http://lab-api:18180`, `http://orch:18182`).

## Comandos útiles

```bash
docker compose logs -f nginx web-ui orch
docker compose down          # conserva volúmenes SQLite
docker compose down -v       # borra datos SQLite
docker compose build lab-api # rebuild servicio individual
```

## Arquitectura

```
Browser → nginx:18321
            ├─ /           → web-ui:4321
            ├─ /lab-api/*  → lab-api:18180
            ├─ /orch/*     → orch:18182  (HTTP + WS + SSE)
            └─ /chavez/*   → chavez:18181

Cursor → docker compose exec → skills-lab-mcp (stdio)
                              ├→ lab-api:18180
                              └→ orch:18182
```

web-ui hace fetch SSR interno a `lab-api` vía `LAB_API_PROXY_TARGET`; el browser solo habla con nginx.
