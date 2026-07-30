package provider

import (
	"context"
	"fmt"
)

// CursorAgent is a stub for the Cursor Agent CLI (`agent`).
type CursorAgent struct {
	base
}

func (c *CursorAgent) Name() string { return "cursor_agent" }

func (c *CursorAgent) Run(_ context.Context, query string) (RunResult, error) {
	return RunResult{}, fmt.Errorf("%w: Cursor Agent provider is a stub. Implement Run() against the agent CLI, then detect skill %q tool-use. Query was: %q",
		ErrNotImplemented, c.cfg.SkillName, query)
}
