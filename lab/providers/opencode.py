"""OpenCode CLI provider — full implementation for initial lab runs."""

from __future__ import annotations

import subprocess
import threading
import time
from pathlib import Path

from .base import Provider, RunResult


class OpenCodeProvider(Provider):
    name = "opencode"

    def run(self, query: str) -> RunResult:
        cmd = ["opencode", "run", "--model", self.model, query]
        proc = subprocess.Popen(
            cmd,
            cwd=str(self.workdir),
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT,
            text=True,
            bufsize=1,
        )

        chunks: list[str] = []
        triggered = False
        timed_out = False
        lock = threading.Lock()

        def reader() -> None:
            nonlocal triggered
            assert proc.stdout is not None
            for line in proc.stdout:
                with lock:
                    chunks.append(line)
                    if not triggered and self.detect_trigger("".join(chunks)):
                        triggered = True
                        if proc.poll() is None:
                            proc.kill()
                        break

        t = threading.Thread(target=reader, daemon=True)
        t.start()

        deadline = time.monotonic() + self.timeout_seconds
        while proc.poll() is None and time.monotonic() < deadline:
            if triggered:
                break
            time.sleep(0.05)

        if proc.poll() is None:
            if not triggered:
                timed_out = True
            proc.kill()

        t.join(timeout=2.0)
        with lock:
            stdout = "".join(chunks)

        if not triggered:
            triggered = self.detect_trigger(stdout)

        return RunResult(
            stdout=stdout,
            triggered=triggered,
            timed_out=timed_out,
            returncode=proc.returncode,
            metadata={"cmd": cmd, "workdir": str(self.workdir)},
        )
