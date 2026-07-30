package timing

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Entry is one stacked timing record with a stable index.
type Entry struct {
	Index       int    `json:"index"`
	EvalID      int    `json:"eval_id"`
	EvalSlug    string `json:"eval_slug"`
	Config      string `json:"config"`
	Run         int    `json:"run"`
	Path        string `json:"path"`
	TotalTokens int64  `json:"total_tokens"`
	DurationMS  int64  `json:"duration_ms"`
	Provider    string `json:"provider,omitempty"`
	ExitCode    int    `json:"exit_code,omitempty"`
	TimedOut    bool   `json:"timed_out,omitempty"`
}

// Stack is the iteration-level timing.json document.
type Stack struct {
	Runs []Entry `json:"runs"`
}

// Path returns timing.json under iterationDir.
func Path(iterationDir string) string {
	return filepath.Join(iterationDir, "timing.json")
}

// Load reads timing.json or returns an empty stack.
func Load(iterationDir string) (Stack, error) {
	p := Path(iterationDir)
	raw, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return Stack{Runs: []Entry{}}, nil
		}
		return Stack{}, err
	}
	var s Stack
	if err := json.Unmarshal(raw, &s); err != nil {
		return Stack{}, fmt.Errorf("%s: %w", p, err)
	}
	if s.Runs == nil {
		s.Runs = []Entry{}
	}
	return s, nil
}

// Save writes timing.json.
func Save(iterationDir string, s Stack) error {
	raw, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	raw = append(raw, '\n')
	return os.WriteFile(Path(iterationDir), raw, 0o644)
}

// Append adds an entry with the next index and persists the stack.
func Append(iterationDir string, e Entry) (Entry, error) {
	s, err := Load(iterationDir)
	if err != nil {
		return Entry{}, err
	}
	next := 1
	for _, r := range s.Runs {
		if r.Index >= next {
			next = r.Index + 1
		}
	}
	e.Index = next
	s.Runs = append(s.Runs, e)
	if err := Save(iterationDir, s); err != nil {
		return Entry{}, err
	}
	return e, nil
}

// WriteRunTiming writes a per-run timing snapshot (single entry) next to the run.
func WriteRunTiming(runDir string, e Entry) error {
	raw, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return err
	}
	raw = append(raw, '\n')
	return os.WriteFile(filepath.Join(runDir, "timing.json"), raw, 0o644)
}
