package grade_test

import (
	"testing"

	"skills-laboratory/lab-go-quality/internal/grade"
)

func TestSummarizeAndParseJudgeJSON(t *testing.T) {
	s := grade.Summarize([]grade.AssertionResult{
		{Text: "a", Passed: true, Evidence: "ok"},
		{Text: "b", Passed: false, Evidence: "no"},
	})
	if s.Passed != 1 || s.Failed != 1 || s.Total != 2 || s.PassRate != 0.5 {
		t.Fatalf("%+v", s)
	}

	raw := "Here you go:\n```json\n{\"assertion_results\":[{\"text\":\"a\",\"passed\":true,\"evidence\":\"found\"}]}\n```\n"
	results, err := grade.ParseJudgeJSON(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || !results[0].Passed {
		t.Fatalf("%+v", results)
	}
}

func TestBuildTaskPromptViaJudgePrompt(t *testing.T) {
	p := grade.JudgePrompt("chart", []string{"has image"}, []string{"chart.png"}, "=== chart.png ===\nbinary", "ran ok")
	if p == "" {
		t.Fatal("empty prompt")
	}
}
