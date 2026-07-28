"""LLM CLI providers for skill trigger evaluation."""

from __future__ import annotations

from .base import Provider, ProviderNotImplementedError, RunResult
from .claude import ClaudeProvider
from .codex import CodexProvider
from .cursor_agent import CursorAgentProvider
from .opencode import OpenCodeProvider

PROVIDERS: dict[str, type[Provider]] = {
    "opencode": OpenCodeProvider,
    "codex": CodexProvider,
    "claude": ClaudeProvider,
    "cursor_agent": CursorAgentProvider,
    "agent": CursorAgentProvider,
}


def get_provider(name: str, **kwargs) -> Provider:
    key = name.strip().lower().replace("-", "_")
    if key not in PROVIDERS:
        known = ", ".join(sorted(PROVIDERS))
        raise ValueError(f"Unknown provider {name!r}. Known: {known}")
    return PROVIDERS[key](**kwargs)


__all__ = [
    "PROVIDERS",
    "ClaudeProvider",
    "CodexProvider",
    "CursorAgentProvider",
    "OpenCodeProvider",
    "Provider",
    "ProviderNotImplementedError",
    "RunResult",
    "get_provider",
]
