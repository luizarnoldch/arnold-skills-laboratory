package provider

import (
	"context"
	"fmt"
)

// Claude is a stub for the Anthropic Claude Code CLI.
type Claude struct {
	base
}

func (c *Claude) Name() string { return "claude" }

func (c *Claude) Run(_ context.Context, query string) (RunResult, error) {
	return RunResult{}, fmt.Errorf("%w: Claude provider is a stub. Implement Run() against the Claude CLI, then detect skill %q tool-use. Query was: %q",
		ErrNotImplemented, c.cfg.SkillName, query)
}
