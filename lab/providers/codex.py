"""Codex CLI provider stub."""

from __future__ import annotations

from .base import Provider, ProviderNotImplementedError, RunResult


class CodexProvider(Provider):
    """
    Contract for OpenAI Codex CLI.

    Expected invocation (to implement):
      codex exec --model <model> <query>
      (or equivalent non-interactive Codex command)

    Trigger detection should look for skill tool-use markers for `skill_name`
    in stdout/stderr (JSON events or Skill "name" style logs).
    """

    name = "codex"

    def run(self, query: str) -> RunResult:
        raise ProviderNotImplementedError(
            "Codex provider is a stub. Implement run() against the Codex CLI, "
            f"then detect skill {self.skill_name!r} tool-use in the output. "
            f"Query was: {query!r}"
        )
