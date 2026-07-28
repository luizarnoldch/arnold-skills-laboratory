"""Cursor Agent CLI provider stub."""

from __future__ import annotations

from .base import Provider, ProviderNotImplementedError, RunResult


class CursorAgentProvider(Provider):
    """
    Contract for Cursor Agent CLI (`agent`).

    Expected invocation (to implement):
      agent -p <query> --model <model>
      (or `cursor agent` equivalent non-interactive mode)

    Skills typically live under `.cursor/skills/<skill_name>/SKILL.md`.
    Detect trigger via skill invocation markers for `skill_name`.
    """

    name = "cursor_agent"

    def run(self, query: str) -> RunResult:
        raise ProviderNotImplementedError(
            "Cursor Agent provider is a stub. Implement run() against the agent CLI, "
            f"then detect skill {self.skill_name!r} tool-use in the output. "
            f"Query was: {query!r}"
        )
