#!/usr/bin/env python3
"""Split a prompt set into stratified train (60%) and validation (40%)."""

from __future__ import annotations

import argparse
import json
import random
import sys
from pathlib import Path


def load_prompts(path: Path) -> list[dict]:
    data = json.loads(path.read_text(encoding="utf-8"))
    if not isinstance(data, list):
        raise ValueError(f"{path}: expected a JSON array")
    for i, item in enumerate(data):
        for key in ("id", "query", "should_trigger"):
            if key not in item:
                raise ValueError(f"{path}: item[{i}] missing {key!r}")
    return data


def stratified_split(
    prompts: list[dict],
    *,
    train_ratio: float = 0.6,
    seed: int = 42,
) -> tuple[list[dict], list[dict]]:
    rng = random.Random(seed)
    positives = [p for p in prompts if p["should_trigger"]]
    negatives = [p for p in prompts if not p["should_trigger"]]
    rng.shuffle(positives)
    rng.shuffle(negatives)

    def take(items: list[dict]) -> tuple[list[dict], list[dict]]:
        n_train = int(round(len(items) * train_ratio))
        # Prefer exact 60/40 for even counts (e.g. 10 -> 6/4)
        return items[:n_train], items[n_train:]

    pos_train, pos_val = take(positives)
    neg_train, neg_val = take(negatives)

    train = sorted(pos_train + neg_train, key=lambda p: p["id"])
    validation = sorted(pos_val + neg_val, key=lambda p: p["id"])
    return train, validation


def write_json(path: Path, data: list[dict]) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(data, indent=2, ensure_ascii=False) + "\n", encoding="utf-8")


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("-i", "--input", required=True, type=Path, help="Full prompts.json")
    parser.add_argument("--train-out", type=Path, help="Output train.json (default: sibling train.json)")
    parser.add_argument(
        "--validation-out",
        type=Path,
        help="Output validation.json (default: sibling validation.json)",
    )
    parser.add_argument("--train-ratio", type=float, default=0.6, help="Train fraction (default 0.6)")
    parser.add_argument("--seed", type=int, default=42, help="RNG seed for reproducibility")
    args = parser.parse_args()

    prompts = load_prompts(args.input)
    train, validation = stratified_split(prompts, train_ratio=args.train_ratio, seed=args.seed)

    train_out = args.train_out or args.input.with_name("train.json")
    val_out = args.validation_out or args.input.with_name("validation.json")

    write_json(train_out, train)
    write_json(val_out, validation)

    pos_t = sum(1 for p in train if p["should_trigger"])
    neg_t = len(train) - pos_t
    pos_v = sum(1 for p in validation if p["should_trigger"])
    neg_v = len(validation) - pos_v

    print(f"Wrote {train_out} ({len(train)} = {pos_t} pos / {neg_t} neg)")
    print(f"Wrote {val_out} ({len(validation)} = {pos_v} pos / {neg_v} neg)")
    return 0


if __name__ == "__main__":
    sys.exit(main())
