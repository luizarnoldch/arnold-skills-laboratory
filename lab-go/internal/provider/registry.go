package provider

import (
	"fmt"
	"strings"
	"time"
)

// New returns a provider by name.
func New(name string, cfg Config) (Provider, error) {
	if cfg.Timeout <= 0 {
		cfg.Timeout = 60 * time.Second
	}
	key := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(name), "-", "_"))
	switch key {
	case "opencode":
		return &OpenCode{base: base{cfg: cfg}}, nil
	case "codex":
		return &Codex{base: base{cfg: cfg}}, nil
	case "claude":
		return &Claude{base: base{cfg: cfg}}, nil
	case "cursor_agent", "agent":
		return &CursorAgent{base: base{cfg: cfg}}, nil
	default:
		return nil, fmt.Errorf("unknown provider %q (opencode|codex|claude|cursor_agent)", name)
	}
}
