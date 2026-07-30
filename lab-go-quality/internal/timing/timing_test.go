package timing_test

import (
	"os"
	"path/filepath"
	"testing"

	"skills-laboratory/lab-go-quality/internal/timing"
)

func TestAppendStacksWithIndex(t *testing.T) {
	dir := t.TempDir()
	e1, err := timing.Append(dir, timing.Entry{
		EvalID: 1, EvalSlug: "eval-a", Config: "with_skill", Run: 1,
		Path: "eval-a/with_skill/run_001", DurationMS: 100, TotalTokens: 10,
	})
	if err != nil {
		t.Fatal(err)
	}
	if e1.Index != 1 {
		t.Fatalf("index=%d", e1.Index)
	}
	e2, err := timing.Append(dir, timing.Entry{
		EvalID: 1, EvalSlug: "eval-a", Config: "without_skill", Run: 1,
		Path: "eval-a/without_skill/run_001", DurationMS: 80, TotalTokens: 5,
	})
	if err != nil {
		t.Fatal(err)
	}
	if e2.Index != 2 {
		t.Fatalf("index=%d", e2.Index)
	}
	stack, err := timing.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(stack.Runs) != 2 {
		t.Fatalf("len=%d", len(stack.Runs))
	}
	runDir := filepath.Join(dir, "r")
	if err := os.MkdirAll(runDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := timing.WriteRunTiming(runDir, e1); err != nil {
		t.Fatal(err)
	}
}
