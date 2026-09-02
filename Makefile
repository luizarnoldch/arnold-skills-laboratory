.PHONY: docker-up docker-down docker-build docker-ps docker-logs docker-verify

docker-up:
	docker compose --env-file .env up -d --build

docker-down:
	docker compose down

docker-build:
	docker compose build

docker-ps:
	docker compose ps

docker-logs:
	docker compose logs -f nginx web-ui lab-api chavez orch skills-lab-mcp

docker-verify:
	@echo "=== lab-api /ready ==="
	@docker compose exec -T lab-api wget -qO- http://127.0.0.1:18180/ready
	@echo ""
	@echo "=== chavez /health ==="
	@docker compose exec -T chavez wget -qO- http://127.0.0.1:18181/health
	@echo ""
	@echo "=== orch /ready ==="
	@docker compose exec -T orch wget -qO- http://127.0.0.1:18182/ready
	@echo ""
	@echo "=== skills-lab-mcp binary ==="
	@docker compose exec -T skills-lab-mcp test -x /app/skills_lab_mcp && echo "OK: /app/skills_lab_mcp"
	@echo ""
	@echo "=== nginx UI ==="
	@curl -sf -o /dev/null -w "HTTP %{http_code}\n" http://127.0.0.1:18321/
	@echo "=== orch routes via nginx ==="
	@curl -sf -o /dev/null -w "baseline-eval-jobs HTTP %{http_code}\n" http://127.0.0.1:18321/orch/api/v1/baseline-eval-jobs
	@curl -sf -o /dev/null -w "trigger-eval-batches HTTP %{http_code}\n" http://127.0.0.1:18321/orch/api/v1/trigger-eval-batches
	@curl -sf -o /dev/null -w "optimize-jobs HTTP %{http_code}\n" http://127.0.0.1:18321/orch/api/v1/optimize-jobs
