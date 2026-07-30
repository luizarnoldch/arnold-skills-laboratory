package runner

import (
	"fmt"
	"strings"
)

// New returns a TaskRunner by provider name.
func New(name string) (TaskRunner, error) {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "opencode":
		return &OpenCode{}, nil
	case "claude":
		return &Claude{}, nil
	case "codex":
		return &Codex{}, nil
	case "agent", "cursor_agent", "cursor-agent":
		return &CursorAgent{}, nil
	default:
		return nil, fmt.Errorf("unknown provider %q (want: opencode, claude, codex, agent)", name)
	}
}
