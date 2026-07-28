#!/usr/bin/env python3
"""Evaluate whether an LLM CLI triggers a skill for a prompt set."""

from __future__ import annotations

import argparse
import json
import sys
import time
from datetime import datetime, timezone
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
if str(ROOT) not in sys.path:
    sys.path.insert(0, str(ROOT))

from lab.providers import get_provider  # noqa: E402


def load_prompts(path: Path) -> list[dict]:
    data = json.loads(path.read_text(encoding="utf-8"))
    if not isinstance(data, list):
        raise ValueError(f"{path}: expected a JSON array")
    return data


def next_log_dir(base: Path) -> Path:
    base.mkdir(parents=True, exist_ok=True)
    existing = []
    for child in base.iterdir():
        if child.is_dir() and child.name.startswith("run_"):
            try:
                existing.append(int(child.name.split("_", 1)[1]))
            except ValueError:
                continue
    n = (max(existing) + 1) if existing else 1
    path = base / f"run_{n}"
    path.mkdir(parents=True, exist_ok=True)
    return path


def evaluate_prompt(
    provider,
    item: dict,
    *,
    runs: int,
    log_dir: Path | None,
    majority_threshold: float,
) -> dict:
    prompt_id = item["id"]
    query = item["query"]
    should_trigger = bool(item["should_trigger"])
    run_count = int(item.get("runs") or runs)
    trigger_count = 0

    for i in range(1, run_count + 1):
        result = provider.run(query)
        if result.triggered:
            trigger_count += 1
            status = "TRIGGERED"
        elif result.timed_out:
            status = "TIMED OUT"
        else:
            status = "NOT TRIGGERED"

        print(f"  └─ Run {i}/{run_count}: {status}")

        if log_dir is not None:
            name = f"id_{prompt_id}.log" if run_count == 1 else f"id_{prompt_id}_run_{i}.log"
            header = (
                "=== RUN METADATA ===\n"
                f"ID: {prompt_id}\n"
                f"Timestamp: {datetime.now(timezone.utc).isoformat()}\n"
                f"Query: {query}\n"
                f"Model: {provider.model}\n"
                f"Provider: {provider.name}\n"
                f"Run: {i}/{run_count}\n"
                f"Expected Trigger: {should_trigger}\n"
                f"Triggered: {result.triggered}\n"
                f"Timed Out: {result.timed_out}\n"
                "====================\n\n"
            )
            (log_dir / name).write_text(header + result.stdout, encoding="utf-8")

    trigger_rate = round(trigger_count / run_count, 2) if run_count else 0.0
    predicted = trigger_rate >= majority_threshold
    return {
        "id": prompt_id,
        "query": query,
        "should_trigger": should_trigger,
        "skill_name": provider.skill_name,
        "triggers": trigger_count,
        "runs": run_count,
        "trigger_rate": trigger_rate,
        "correct": predicted == should_trigger,
    }


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--skill-name", required=True, help="Skill name to detect")
    parser.add_argument("--prompts", required=True, type=Path, help="Prompt JSON (train or validation)")
    parser.add_argument(
        "--provider",
        default="opencode",
        help="CLI provider: opencode | codex | claude | cursor_agent",
    )
    parser.add_argument("--runs", type=int, default=1, help="Runs per prompt (default 1)")
    parser.add_argument("--model", default="digitalocean/deepseek-4-flash", help="Model id for the CLI")
    parser.add_argument("--workdir", type=Path, required=True, help="Sandbox cwd for the CLI")
    parser.add_argument("--out", type=Path, required=True, help="Results JSON output path")
    parser.add_argument("--timeout", type=float, default=60.0, help="Per-run timeout seconds")
    parser.add_argument(
        "--majority-threshold",
        type=float,
        default=0.5,
        help="trigger_rate >= threshold counts as triggered prediction",
    )
    parser.add_argument(
        "--log-dir",
        type=Path,
        help="Base log directory (creates run_N). Default: <workdir>/../logs",
    )
    parser.add_argument("--no-logs", action="store_true", help="Disable per-run log files")
    args = parser.parse_args()

    prompts = load_prompts(args.prompts)
    workdir = args.workdir.resolve()
    if not workdir.is_dir():
        raise SystemExit(f"workdir does not exist: {workdir}")

    provider = get_provider(
        args.provider,
        model=args.model,
        skill_name=args.skill_name,
        workdir=workdir,
        timeout_seconds=args.timeout,
    )

    log_dir = None
    if not args.no_logs:
        base = args.log_dir or (workdir.parent / "logs")
        log_dir = next_log_dir(base.resolve())

    print("=" * 50)
    print(f"Evaluating with provider={provider.name} model={args.model}")
    print(f"Skill: {args.skill_name} | Default runs: {args.runs} | Timeout: {args.timeout}s")
    print(f"Prompts: {args.prompts} ({len(prompts)} items)")
    if log_dir:
        print(f"Log directory: {log_dir}")
    print("=" * 50)

    results: list[dict] = []
    started = time.time()
    for item in prompts:
        runs = int(item.get("runs") or args.runs)
        print(
            f"\n[ID: {item['id']}] Query: {item['query']!r} | "
            f"Expected: {item['should_trigger']} | Runs: {runs}"
        )
        results.append(
            evaluate_prompt(
                provider,
                item,
                runs=args.runs,
                log_dir=log_dir,
                majority_threshold=args.majority_threshold,
            )
        )

    correct = sum(1 for r in results if r["correct"])
    accuracy = round(correct / len(results), 4) if results else 0.0
    payload = {
        "skill_name": args.skill_name,
        "provider": provider.name,
        "model": args.model,
        "prompts_file": str(args.prompts.resolve()),
        "workdir": str(workdir),
        "runs_default": args.runs,
        "majority_threshold": args.majority_threshold,
        "summary": {
            "total": len(results),
            "correct": correct,
            "accuracy": accuracy,
            "accuracy_pct": round(accuracy * 100, 2),
        },
        "elapsed_seconds": round(time.time() - started, 2),
        "results": results,
    }

    args.out.parent.mkdir(parents=True, exist_ok=True)
    args.out.write_text(json.dumps(payload, indent=2, ensure_ascii=False) + "\n", encoding="utf-8")

    print(f"\nAccuracy: {correct}/{len(results)} ({payload['summary']['accuracy_pct']}%)")
    print(f"Results saved to {args.out}")
    if log_dir:
        print(f"Logs saved in {log_dir}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
