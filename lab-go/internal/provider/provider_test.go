package provider_test

import (
	"context"
	"errors"
	"testing"

	"skills-laboratory/lab-go/internal/provider"
)

func TestDetectTrigger(t *testing.T) {
	cases := []struct {
		text string
		want bool
	}{
		{`Skill "feature-expert" loaded`, true},
		{`{"name":"feature-expert"}`, true},
		{`no skill here`, false},
		{`skill:feature-expert`, true},
	}
	for _, tc := range cases {
		if got := provider.DetectTrigger(tc.text, "feature-expert"); got != tc.want {
			t.Errorf("DetectTrigger(%q)=%v want %v", tc.text, got, tc.want)
		}
	}
}

func TestStubProviders(t *testing.T) {
	for _, name := range []string{"codex", "claude", "cursor_agent", "agent"} {
		p, err := provider.New(name, provider.Config{Model: "m", SkillName: "feature-expert", Workdir: "."})
		if err != nil {
			t.Fatalf("New(%s): %v", name, err)
		}
		_, err = p.Run(context.Background(), "hello")
		if !errors.Is(err, provider.ErrNotImplemented) {
			t.Fatalf("%s: want ErrNotImplemented, got %v", name, err)
		}
	}
}
