# skills-lab-mcp

Servidor MCP (stdio) en **Go** que expone las funcionalidades principales del Skills Laboratory para LLMs externos (Cursor, Claude Desktop, etc.).

Fachada sobre:

- **arnold-laboratory-api** (`:18180`) — skills, versiones, test sets
- **arnold-lab-orchestrator** (`:18182`) — evaluaciones y optimización async

## Prerrequisitos

Stack de laboratorio corriendo:

```bash
make docker-up   # desde la raíz del monorepo
```

Chavez (`:18181`) debe estar activo para jobs de eval/optimize.

## Instalación

### Docker (recomendado con stack completo)

El servicio se construye y levanta con el stack principal:

```bash
make docker-up   # desde la raíz del monorepo
```

El contenedor queda activo; Cursor invoca el MCP vía [`deploy/docker/mcp-exec.sh`](../../deploy/docker/mcp-exec.sh).

### Local (desarrollo)

```bash
cd services/skills-lab-mcp
make build
```

Binario: `bin/skills_lab_mcp`

## Configuración

| Variable | Default (local) | Default (Docker) | Descripción |
|----------|-----------------|------------------|-------------|
| `LAB_API_URL` | `http://127.0.0.1:18180` | `http://lab-api:18180` | Base URL laboratory-api |
| `ORCH_API_URL` | `http://127.0.0.1:18182` | `http://orch:18182` | Base URL orchestrator |
| `WORKSPACE_DEFAULT` | `.` | `/workspace` | Workspace para evals |
| `HTTP_TIMEOUT` | `30s` | `30s` | Timeout HTTP clientes |

Vía nginx docker desde el **host** (solo si usas binario local, no el contenedor MCP):

```bash
LAB_API_URL=http://127.0.0.1:18321/lab-api
ORCH_API_URL=http://127.0.0.1:18321/orch
```

### Cursor con Docker (`~/.cursor/mcp.json`)

Requiere stack levantado (`make docker-up`):

```json
{
  "mcpServers": {
    "skills-lab": {
      "command": "/datos/Work/command-center/arnold-skills-laboratory/deploy/docker/mcp-exec.sh"
    }
  }
}
```

### Cursor con binario local (`~/.cursor/mcp.json`)

```json
{
  "mcpServers": {
    "skills-lab": {
      "command": "/ruta/absoluta/services/skills-lab-mcp/bin/skills_lab_mcp",
      "env": {
        "LAB_API_URL": "http://127.0.0.1:18180",
        "ORCH_API_URL": "http://127.0.0.1:18182"
      }
    }
  }
}
```

### Migración desde la versión TypeScript (eliminada)

Si tenías configurado el MCP anterior con Node, reemplaza por el script Docker o el binario Go (ver arriba). Elimina entradas con `node` + `dist/index.js` o `skills-lab-go`.

## Tools (9)

| Tool | Descripción |
|------|-------------|
| `skills_list` | Lista skills |
| `skill_get` | Detalle por id/nombre + versiones current/test |
| `test_set_upload` | Carga prompts + split auto train/validation |
| `test_set_list` | Lista prompts por split |
| `test_set_update` | Modifica query, should_trigger o split |
| `baseline_eval_start` | Eval train+validation → `job_id` |
| `trigger_eval_start` | Eval por split → `job_id` |
| `optimize_start` | Optimización → `job_id` |
| `job_get` | Consulta job async |

Jobs async devuelven `{ job_id, job_type, status }`. Usa `job_get` hasta `status: "completed"`.

## Flujos E2E

### Explorar skill

```
skills_list → {}
skill_get → { "skill_name": "mi-skill" }
```

### Cargar test set

```
test_set_upload → {
  "skill_name": "mi-skill",
  "prompts": [
    { "prompt_index": 1, "query": "¿Cómo formateo LaTeX?", "should_trigger": true },
    { "prompt_index": 2, "query": "Lista archivos", "should_trigger": false }
  ]
}
```

### Evaluar baseline

```
baseline_eval_start → { "skill_name": "mi-skill", "runs_default": 3 }
job_get → { "job_type": "baseline_eval", "job_id": "<uuid>" }
```

### Optimizar

```
optimize_start → { "skill_name": "mi-skill", "max_iters": 3, "threshold": 0.95 }
job_get → { "job_type": "optimize", "job_id": "<uuid>" }
```

## Desarrollo

```bash
make build
make test
make run    # stdio — para MCP Inspector
```

SDK: [`github.com/modelcontextprotocol/go-sdk`](https://github.com/modelcontextprotocol/go-sdk)

## Arquitectura

```
LLM ←stdio→ skills-lab-mcp
              ├→ lab-api
              └→ orchestrator → chavez
```

Ver [`../AGENTS.md`](../AGENTS.md).
