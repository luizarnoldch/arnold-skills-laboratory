package benchmark_test

import (
	"os"
	"path/filepath"
	"testing"

	"skills-laboratory/lab-go-quality/internal/benchmark"
	"skills-laboratory/lab-go-quality/internal/grade"
	"skills-laboratory/lab-go-quality/internal/timing"
)

func TestComputeDelta(t *testing.T) {
	dir := t.TempDir()
	mustAppend := func(e timing.Entry) {
		t.Helper()
		if _, err := timing.Append(dir, e); err != nil {
			t.Fatal(err)
		}
	}
	mustAppend(timing.Entry{
		EvalID: 1, EvalSlug: "eval-a", Config: "with_skill", Run: 1,
		Path: "eval-a/with_skill/run_001", DurationMS: 2000, TotalTokens: 100,
	})
	mustAppend(timing.Entry{
		EvalID: 1, EvalSlug: "eval-a", Config: "without_skill", Run: 1,
		Path: "eval-a/without_skill/run_001", DurationMS: 1000, TotalTokens: 50,
	})

	writeGrading := func(rel string, rate float64) {
		t.Helper()
		p := filepath.Join(dir, rel)
		_ = os.MkdirAll(p, 0o755)
		rep := grade.Report{
			AssertionResults: []grade.AssertionResult{
				{Text: "x", Passed: rate >= 1, Evidence: "e"},
			},
		}
		if rate == 0.5 {
			rep.AssertionResults = []grade.AssertionResult{
				{Text: "x", Passed: true, Evidence: "e"},
				{Text: "y", Passed: false, Evidence: "e"},
			}
		}
		if err := grade.Write(filepath.Join(p, "grading.json"), rep); err != nil {
			t.Fatal(err)
		}
	}
	writeGrading("eval-a/with_skill/run_001", 1)
	writeGrading("eval-a/without_skill/run_001", 0.5)

	rep, err := benchmark.Compute(dir)
	if err != nil {
		t.Fatal(err)
	}
	if rep.RunSummary.WithSkill.PassRate.Mean != 1 {
		t.Fatalf("with pass=%v", rep.RunSummary.WithSkill.PassRate)
	}
	if rep.RunSummary.WithoutSkill.PassRate.Mean != 0.5 {
		t.Fatalf("without pass=%v", rep.RunSummary.WithoutSkill.PassRate)
	}
	if rep.RunSummary.Delta.PassRate != 0.5 {
		t.Fatalf("delta=%v", rep.RunSummary.Delta)
	}
	if rep.RunSummary.WithSkill.TimeSeconds.Mean != 2 {
		t.Fatalf("time=%v", rep.RunSummary.WithSkill.TimeSeconds)
	}
	if err := benchmark.Write(dir, rep); err != nil {
		t.Fatal(err)
	}
}
