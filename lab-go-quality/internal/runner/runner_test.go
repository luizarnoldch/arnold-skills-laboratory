package runner_test

import (
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
	n := runner.ParseTokensBestEffort(`{"total_tokens": 1234}`)
	if n != 1234 {
		t.Fatalf("tokens=%d", n)
	}
}
