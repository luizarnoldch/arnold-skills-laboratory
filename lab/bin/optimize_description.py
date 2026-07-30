#!/usr/bin/env python3
"""
Train loop: evaluate train prompts and optionally rewrite SKILL.md description.

Validation data must NEVER be passed here — use lab/bin/evaluate.py on validation.json.
"""

from __future__ import annotations

import argparse
import json
import re
import subprocess
import sys
from datetime import datetime, timezone
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
if str(ROOT) not in sys.path:
    sys.path.insert(0, str(ROOT))

FRONTMATTER_RE = re.compile(r"^---\n(.*?)\n---\n?", re.DOTALL)
DESCRIPTION_RE = re.compile(
    r"^(description:\s*)(>?\s*.*?)(?=^[a-zA-Z_][\w-]*:|\Z)",
    re.DOTALL | re.MULTILINE,
)


def read_skill(path: Path) -> tuple[str, str, str]:
    """Return (raw, frontmatter, body)."""
    raw = path.read_text(encoding="utf-8")
    m = FRONTMATTER_RE.match(raw)
    if not m:
        raise ValueError(f"{path}: missing YAML frontmatter")
    return raw, m.group(1), raw[m.end() :]


def get_description(frontmatter: str) -> str:
    m = DESCRIPTION_RE.search(frontmatter)
    if not m:
        raise ValueError("frontmatter has no description field")
    value = m.group(2).strip()
    if value.startswith(">"):
        lines = []
        for line in value.splitlines():
            line = line.strip()
            if line.startswith(">"):
                line = line[1:].strip()
            lines.append(line)
        return " ".join(lines).strip()
    return value.strip().strip('"').strip("'")


def set_description(frontmatter: str, new_description: str) -> str:
    folded = new_description.strip().replace("\n", " ")
    block = f"description: >\n  {folded}\n"
    m = DESCRIPTION_RE.search(frontmatter)
    if not m:
        return frontmatter.rstrip() + "\n" + block
    start, end = m.span()
    return frontmatter[:start] + block + frontmatter[end:].lstrip("\n")


def write_skill(path: Path, frontmatter: str, body: str) -> None:
    body_text = body if body.startswith("\n") or not body else "\n" + body
    content = f"---\n{frontmatter.rstrip()}\n---{body_text}"
    if not content.endswith("\n"):
        content += "\n"
    path.write_text(content, encoding="utf-8")


def next_iteration_dir(base: Path) -> Path:
    base.mkdir(parents=True, exist_ok=True)
    nums = []
    for child in base.iterdir():
        if child.is_dir() and child.name.isdigit():
            nums.append(int(child.name))
    n = (max(nums) + 1) if nums else 1
    path = base / f"{n:03d}"
    path.mkdir(parents=True, exist_ok=True)
    return path


def run_evaluate(args: argparse.Namespace, out_path: Path) -> dict:
    cmd = [
        sys.executable,
        str(ROOT / "lab" / "bin" / "evaluate.py"),
        "--skill-name",
        args.skill_name,
        "--prompts",
        str(args.prompts),
        "--provider",
        args.provider,
        "--runs",
        str(args.runs),
        "--model",
        args.model,
        "--workdir",
        str(args.workdir),
        "--out",
        str(out_path),
        "--timeout",
        str(args.timeout),
        "--majority-threshold",
        str(args.majority_threshold),
    ]
    if args.log_dir:
        cmd.extend(["--log-dir", str(args.log_dir)])
    if args.no_logs:
        cmd.append("--no-logs")
    print("→", " ".join(cmd))
    subprocess.check_call(cmd)
    return json.loads(out_path.read_text(encoding="utf-8"))


def failures_from_results(payload: dict) -> list[dict]:
    return [r for r in payload.get("results", []) if not r.get("correct")]


def propose_description_with_opencode(
    *,
    workdir: Path,
    model: str,
    skill_name: str,
    current_description: str,
    failures: list[dict],
    timeout: float,
) -> str:
    fail_lines = []
    for f in failures:
        kind = "false_negative" if f["should_trigger"] else "false_positive"
        fail_lines.append(
            f"- id={f['id']} ({kind}) trigger_rate={f['trigger_rate']}: {f['query']!r}"
        )
    prompt = f"""You are optimizing the YAML frontmatter `description` of an agent skill named "{skill_name}".

Current description:
\"\"\"{current_description}\"\"\"

These train prompts were classified incorrectly (skill trigger vs expected):
{chr(10).join(fail_lines) or '(none)'}

Rewrite ONLY the description text so the model more reliably triggers on should_trigger=true
prompts and does NOT trigger on should_trigger=false prompts.

Rules:
- Return ONLY the new description plain text (no YAML, no markdown fences, no quotes wrapper).
- Keep it to 1-3 sentences / one folded paragraph.
- Include clear positive triggers AND explicit Do NOT / Does NOT negatives when helpful.
"""
    cmd = ["opencode", "run", "--model", model, prompt]
    try:
        completed = subprocess.run(
            cmd,
            cwd=str(workdir),
            capture_output=True,
            text=True,
            timeout=timeout,
            check=False,
        )
    except subprocess.TimeoutExpired as exc:
        raise RuntimeError("description optimization timed out") from exc

    text = (completed.stdout or "") + "\n" + (completed.stderr or "")
    # Prefer last non-empty substantial paragraph
    lines = [ln.strip() for ln in text.splitlines() if ln.strip()]
    if not lines:
        raise RuntimeError("optimizer returned empty output")
    # Drop obvious tool/log noise lines
    cleaned = [
        ln
        for ln in lines
        if not ln.startswith("Skill ")
        and '"name":' not in ln
        and not ln.startswith("===")
    ]
    candidate = " ".join(cleaned[-3:] if cleaned else lines[-3:]).strip()
    candidate = candidate.strip("`").strip('"').strip("'")
    if len(candidate) < 40:
        raise RuntimeError(f"optimizer output too short: {candidate!r}")
    return candidate


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--skill-name", required=True)
    parser.add_argument("--skill-md", required=True, type=Path, help="Path to SKILL.md to optimize")
    parser.add_argument("--prompts", required=True, type=Path, help="Train prompts JSON only")
    parser.add_argument("--workdir", required=True, type=Path)
    parser.add_argument("--iterations-dir", required=True, type=Path)
    parser.add_argument("--results-dir", required=True, type=Path, help="Directory for train result JSONs")
    parser.add_argument("--provider", default="opencode")
    parser.add_argument("--model", default="digitalocean/deepseek-4-flash")
    parser.add_argument("--runs", type=int, default=3)
    parser.add_argument("--timeout", type=float, default=60.0)
    parser.add_argument("--threshold", type=float, default=0.95, help="Stop when train accuracy >= this")
    parser.add_argument("--majority-threshold", type=float, default=0.5)
    parser.add_argument("--max-iters", type=int, default=5)
    parser.add_argument("--log-dir", type=Path)
    parser.add_argument("--no-logs", action="store_true")
    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="Evaluate once and report; do not rewrite description",
    )
    args = parser.parse_args()

    if "validation" in args.prompts.name.lower():
        print(
            "WARNING: prompts path looks like validation data. "
            "optimize_description must only use train.json.",
            file=sys.stderr,
        )

    skill_md = args.skill_md.resolve()
    args.results_dir.mkdir(parents=True, exist_ok=True)

    for iteration in range(1, args.max_iters + 1):
        print(f"\n######## TRAIN ITERATION {iteration}/{args.max_iters} ########")
        out_path = args.results_dir / f"iter_{iteration:03d}.json"
        payload = run_evaluate(args, out_path)
        accuracy = float(payload["summary"]["accuracy"])
        fails = failures_from_results(payload)

        raw, frontmatter, body = read_skill(skill_md)
        current_desc = get_description(frontmatter)
        iter_dir = next_iteration_dir(args.iterations_dir.resolve())
        (iter_dir / "SKILL.description.md").write_text(
            f"# description snapshot\n\n{current_desc}\n", encoding="utf-8"
        )
        decision = "stop_threshold_met" if accuracy >= args.threshold else "optimize"
        if args.dry_run:
            decision = "dry_run_stop"
        if not fails and accuracy >= args.threshold:
            decision = "stop_perfect"

        metrics = {
            "iteration": iteration,
            "timestamp": datetime.now(timezone.utc).isoformat(),
            "accuracy": accuracy,
            "accuracy_pct": payload["summary"]["accuracy_pct"],
            "threshold": args.threshold,
            "failures": fails,
            "decision": decision,
            "results_file": str(out_path),
        }
        (iter_dir / "metrics.json").write_text(
            json.dumps(metrics, indent=2, ensure_ascii=False) + "\n", encoding="utf-8"
        )
        print(
            f"Train accuracy: {payload['summary']['accuracy_pct']}% | "
            f"failures={len(fails)} | decision={decision}"
        )

        if decision.startswith("stop") or decision == "dry_run_stop":
            print("Stopping train loop.")
            return 0

        if args.provider != "opencode":
            print(
                f"Auto-optimize currently uses opencode for rewriting descriptions; "
                f"provider={args.provider} will still evaluate but rewrite requires opencode.",
                file=sys.stderr,
            )
            if args.provider != "opencode":
                # Still attempt rewrite via opencode helper for description text only
                pass

        try:
            new_desc = propose_description_with_opencode(
                workdir=args.workdir.resolve(),
                model=args.model,
                skill_name=args.skill_name,
                current_description=current_desc,
                failures=fails,
                timeout=max(args.timeout, 120.0),
            )
        except Exception as exc:  # noqa: BLE001
            print(f"Failed to propose new description: {exc}", file=sys.stderr)
            metrics["decision"] = "optimize_failed"
            (iter_dir / "metrics.json").write_text(
                json.dumps(metrics, indent=2, ensure_ascii=False) + "\n", encoding="utf-8"
            )
            return 1

        new_fm = set_description(frontmatter, new_desc)
        write_skill(skill_md, new_fm, body)
        (iter_dir / "proposed_description.md").write_text(new_desc + "\n", encoding="utf-8")
        print(f"Updated description in {skill_md}")
        print(f"Snapshot: {iter_dir}")

    print("Reached max_iters without meeting threshold.")
    return 0


if __name__ == "__main__":
    sys.exit(main())
