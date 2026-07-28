"""Claude Code CLI provider stub."""

from __future__ import annotations

from .base import Provider, ProviderNotImplementedError, RunResult


class ClaudeProvider(Provider):
    """
    Contract for Anthropic Claude Code CLI.

    Expected invocation (to implement):
      claude -p <query> --model <model>
      (non-interactive / print mode)

    Skills typically live under `.claude/skills/<skill_name>/SKILL.md`.
    Detect trigger via Skill tool-use events mentioning `skill_name`.
    """

    name = "claude"

    def run(self, query: str) -> RunResult:
        raise ProviderNotImplementedError(
            "Claude provider is a stub. Implement run() against the Claude CLI, "
            f"then detect skill {self.skill_name!r} tool-use in the output. "
            f"Query was: {query!r}"
        )
