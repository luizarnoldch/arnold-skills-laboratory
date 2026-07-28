package evalset_test

import (
	"os"
	"path/filepath"
	"testing"

	"skills-laboratory/lab-go-quality/internal/evalset"
)

func TestLoadValidateAndSlug(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "evals.json")
	set := evalset.Set{
		SkillName: "demo",
		Evals: []evalset.Case{
			{ID: 1, Name: "Top Months Chart", Prompt: "make a chart", ExpectedOutput: "chart"},
			{ID: 2, Prompt: "Clean missing emails please now", ExpectedOutput: "csv"},
		},
	}
	if err := evalset.Write(path, set); err != nil {
		t.Fatal(err)
	}
	got, err := evalset.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.SkillName != "demo" || len(got.Evals) != 2 {
		t.Fatalf("unexpected load: %+v", got)
	}
	if s := evalset.Slug(got.Evals[0]); s != "eval-top-months-chart" {
		t.Fatalf("slug name: %q", s)
	}
	if s := evalset.Slug(got.Evals[1]); s != "eval-2-clean-missing-emails-please" {
		t.Fatalf("slug prompt: %q", s)
	}
}

func TestValidateRejectsDuplicateID(t *testing.T) {
	s := evalset.Set{
		SkillName: "x",
		Evals: []evalset.Case{
			{ID: 1, Prompt: "a", ExpectedOutput: "b"},
			{ID: 1, Prompt: "c", ExpectedOutput: "d"},
		},
	}
	if err := s.Validate(); err == nil {
		t.Fatal("expected duplicate id error")
	}
}

func TestLoadRealFeatureExpertEvals(t *testing.T) {
	path := filepath.Join("..", "..", "..", "development", "skills", "feature-expert", "evals", "evals.json")
	if _, err := os.Stat(path); err != nil {
		t.Skip("feature-expert evals not present")
	}
	if _, err := evalset.Load(path); err != nil {
		t.Fatal(err)
	}
}
