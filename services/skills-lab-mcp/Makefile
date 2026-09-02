LAB_API_URL ?= http://127.0.0.1:18180
ORCH_API_URL ?= http://127.0.0.1:18182
WORKSPACE_DEFAULT ?= .
HTTP_TIMEOUT ?= 30s

export LAB_API_URL ORCH_API_URL WORKSPACE_DEFAULT HTTP_TIMEOUT

.PHONY: build run test clean

build:
	mkdir -p bin
	go build -o bin/skills_lab_mcp ./cmd/skills_lab_mcp

run:
	go run ./cmd/skills_lab_mcp

test:
	go test ./...

clean:
	rm -rf bin
