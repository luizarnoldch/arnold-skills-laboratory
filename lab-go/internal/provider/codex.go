package provider

import (
	"context"
	"fmt"
)

// Codex is a stub for the OpenAI Codex CLI.
type Codex struct {
	base
}

func (c *Codex) Name() string { return "codex" }

func (c *Codex) Run(_ context.Context, query string) (RunResult, error) {
	return RunResult{}, fmt.Errorf("%w: Codex provider is a stub. Implement Run() against the Codex CLI, then detect skill %q tool-use. Query was: %q",
		ErrNotImplemented, c.cfg.SkillName, query)
}
