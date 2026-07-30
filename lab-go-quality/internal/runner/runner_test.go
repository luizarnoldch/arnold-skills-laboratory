package runner_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"skills-laboratory/lab-go-quality/internal/runner"
)

func TestNewProviders(t *testing.T) {
	for _, name := range []string{"opencode", "claude", "codex", "agent", "cursor_agent"} {
		r, err := runner.New(name)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if r.Name() == "" {
			t.Fatalf("%s empty name", name)
		}
	}
	if _, err := runner.New("nope"); err == nil {
		t.Fatal("expected error")
	}
}

func TestBuildTaskPromptAndTokens(t *testing.T) {
	p := runner.BuildTaskPrompt(runner.TaskRequest{
		Prompt:     "do the thing",
		SkillPath:  "/skills/demo",
		InputFiles: []string{"/tmp/a.csv"},
		OutputDir:  "/out",
	})
	if !strings.Contains(p, "Skill path") || !strings.Contains(p, "do the thing") {
		t.Fatalf("bad prompt: %s", p)
	}
	cases := []struct {
		in   string
		want int64
	}{
		{`{"total_tokens": 1234}`, 1234},
		{`Token usage: 999`, 999},
		{`input_tokens=100 output_tokens=50`, 150},
		{"prompt_tokens: 10\ncompletion_tokens: 5", 15},
		{`no usage here`, 0},
	}
	for _, tc := range cases {
		if n := runner.ParseTokensBestEffort(tc.in); n != tc.want {
			t.Fatalf("%q => %d want %d", tc.in, n, tc.want)
		}
	}
}

func TestEvalInputDestRel(t *testing.T) {
	cases := []struct {
		src, name, want string
	}{
		{
			"/repo/development/skills/x/evals/files/create-feature/FEATURES.yml",
			"create-feature",
			"FEATURES.yml",
		},
		{
			"/repo/development/skills/x/evals/files/create-feature-link-prd/spect/features/prd_00010/index.md",
			"create-feature-link-prd",
			"spect/features/prd_00010/index.md",
		},
		{
			"/repo/development/skills/x/evals/files/FEATURES.yml",
			"create-feature",
			"FEATURES.yml",
		},
		{
			"/tmp/plain.csv",
			"create-feature",
			"plain.csv",
		},
	}
	for _, tc := range cases {
		got := runner.EvalInputDestRel(tc.src, tc.name)
		if got != tc.want {
			t.Fatalf("src=%s name=%s got=%q want=%q", tc.src, tc.name, got, tc.want)
		}
	}
}

func TestCopyEvalInputsPreservesTree(t *testing.T) {
	srcRoot := t.TempDir()
	filesRoot := filepath.Join(srcRoot, "evals", "files", "link-feature-prd")
	prd := filepath.Join(filesRoot, "spect", "features", "prd_00010", "index.md")
	feat := filepath.Join(filesRoot, "FEATURES.yml")
	if err := os.MkdirAll(filepath.Dir(prd), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(prd, []byte("# PRD\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(feat, []byte("features: []\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	dst := t.TempDir()
	dests, err := runner.CopyEvalInputs(dst, []string{feat, prd}, "link-feature-prd")
	if err != nil {
		t.Fatal(err)
	}
	if len(dests) != 2 {
		t.Fatalf("dests=%v", dests)
	}
	if _, err := os.Stat(filepath.Join(dst, "FEATURES.yml")); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dst, "spect", "features", "prd_00010", "index.md")); err != nil {
		t.Fatal(err)
	}
}
