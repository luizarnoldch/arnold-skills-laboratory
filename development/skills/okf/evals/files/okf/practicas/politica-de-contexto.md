---
type: Convention
title: Context window policy
description: How we manage token budget in our agents.
tags: [contexto, agentes, convenciones]
generated: { by: human:team, at: 2026-07-30T10:00:00Z }
status: stable
---

# Rule

Keep the context window under 70% before compacting. Past that, quality drops
and we prefer summarizing the thread over accumulating more.

# Related concepts

Filled via [RAG when appropriate](/practicas/cuando-usar-rag.md), watched with
the [context rot runbook](/runbooks/detectar-context-rot.md), and limited by
what [MCP](/runbooks/conectar-mcp.md) injects.

We also follow the still-missing [compaction checklist](/runbooks/compactar-contexto.md)
when the window crosses the threshold.
