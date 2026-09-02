#!/bin/sh
set -e
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
ENV_FILE="$ROOT/.env"
if [ ! -f "$ENV_FILE" ]; then
	ENV_FILE="$ROOT/.env.docker.example"
fi
exec docker compose -f "$ROOT/docker-compose.yml" --env-file "$ENV_FILE" \
	exec -T -i skills-lab-mcp /app/skills_lab_mcp
