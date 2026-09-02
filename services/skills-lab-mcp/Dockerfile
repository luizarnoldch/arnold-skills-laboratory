# syntax=docker/dockerfile:1

FROM golang:1.26-bookworm AS build
WORKDIR /src
COPY services/skills-lab-mcp/go.mod services/skills-lab-mcp/go.sum ./
RUN go mod download
COPY services/skills-lab-mcp/ ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/skills_lab_mcp ./cmd/skills_lab_mcp

FROM debian:bookworm-slim AS runtime
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates coreutils && \
    rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=build /out/skills_lab_mcp ./
ENV LAB_API_URL=http://lab-api:18180 \
    ORCH_API_URL=http://orch:18182 \
    WORKSPACE_DEFAULT=/workspace \
    HTTP_TIMEOUT=30s
CMD ["sleep", "infinity"]
