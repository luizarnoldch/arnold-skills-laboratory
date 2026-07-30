"""Provider contract for skill trigger evaluation CLIs."""

from __future__ import annotations

from abc import ABC, abstractmethod
from dataclasses import dataclass, field
from pathlib import Path
from typing import Optional


class ProviderNotImplementedError(NotImplementedError):
    """Raised when a provider stub is invoked before implementation."""


@dataclass
class RunResult:
    """Outcome of a single provider invocation for one prompt run."""

    stdout: str
    stderr: str = ""
    triggered: bool = False
    timed_out: bool = False
    returncode: Optional[int] = None
    metadata: dict = field(default_factory=dict)


class Provider(ABC):
    """
    Shared interface for LLM CLIs used by lab/bin/evaluate.py.

    Implementations must:
    - run `query` in `workdir`
    - stream or capture output
    - detect whether `skill_name` was invoked (provider-specific markers)
    - optionally terminate early once a trigger is observed
    """

    name: str = "base"

    def __init__(
        self,
        *,
        model: str,
        skill_name: str,
        workdir: Path,
        timeout_seconds: float = 60.0,
    ) -> None:
        self.model = model
        self.skill_name = skill_name
        self.workdir = Path(workdir).resolve()
        self.timeout_seconds = timeout_seconds

    @abstractmethod
    def run(self, query: str) -> RunResult:
        """Execute one prompt and return stdout plus trigger detection."""

    def detect_trigger(self, text: str) -> bool:
        """Default marker heuristics; providers may override."""
        targets = (
            f'Skill "{self.skill_name}"',
            f'"name":"{self.skill_name}"',
            f"'name': '{self.skill_name}'",
            f"skill:{self.skill_name}",
            f"/{self.skill_name}",
        )
        return any(t in text for t in targets)
