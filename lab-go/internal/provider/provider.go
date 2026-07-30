package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

// ErrNotImplemented is returned by stub providers.
var ErrNotImplemented = errors.New("provider not implemented")

// RunResult is the outcome of a single CLI invocation.
type RunResult struct {
	Stdout    string
	Triggered bool
	TimedOut  bool
	Code      int
}

// Config holds shared provider settings.
type Config struct {
	Model     string
	SkillName string
	Workdir   string
	Timeout   time.Duration
}

// Provider runs prompts against an LLM CLI and detects skill triggers.
type Provider interface {
	Name() string
	Model() string
	SkillName() string
	Run(ctx context.Context, query string) (RunResult, error)
}

// DetectTrigger reports whether skill_name appears via known markers.
func DetectTrigger(text, skillName string) bool {
	targets := []string{
		fmt.Sprintf(`Skill "%s"`, skillName),
		fmt.Sprintf(`"name":"%s"`, skillName),
		fmt.Sprintf(`'name': '%s'`, skillName),
		fmt.Sprintf("skill:%s", skillName),
		fmt.Sprintf("/%s", skillName),
	}
	for _, t := range targets {
		if strings.Contains(text, t) {
			return true
		}
	}
	return false
}

type base struct {
	cfg Config
}

func (b base) Model() string     { return b.cfg.Model }
func (b base) SkillName() string { return b.cfg.SkillName }
