---
type: Decision
title: When we use RAG (and when we do not)
description: Team criteria for RAG versus curated prompt or OKF context.
tags: [rag, decision, contexto]
generated: { by: human:team, at: 2026-07-30T10:00:00Z }
status: stable
---

# Decision

Use RAG only when the source does not fit entirely and changes often. For
stable documentation we prefer a curated OKF bundle.

# Why

Retrieved chunks cost budget under the
[context window policy](/practicas/politica-de-contexto.md): every injected
chunk is context we cannot use for reasoning.
