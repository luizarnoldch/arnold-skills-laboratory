package grade

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// AssertionResult is one graded assertion.
type AssertionResult struct {
	Text     string `json:"text"`
	Passed   bool   `json:"passed"`
	Evidence string `json:"evidence"`
}

// Summary aggregates assertion results.
type Summary struct {
	Passed   int     `json:"passed"`
	Failed   int     `json:"failed"`
	Total    int     `json:"total"`
	PassRate float64 `json:"pass_rate"`
}

// Report is grading.json.
type Report struct {
	Index             int               `json:"index,omitempty"`
	EvalID            int               `json:"eval_id,omitempty"`
	EvalSlug          string            `json:"eval_slug,omitempty"`
	Config            string            `json:"config,omitempty"`
	Run               int               `json:"run,omitempty"`
	Path              string            `json:"path,omitempty"`
	AssertionResults  []AssertionResult `json:"assertion_results"`
	Summary           Summary           `json:"summary"`
}

// Summarize computes summary from assertion results.
func Summarize(results []AssertionResult) Summary {
	s := Summary{Total: len(results)}
	for _, r := range results {
		if r.Passed {
			s.Passed++
		} else {
			s.Failed++
		}
	}
	if s.Total > 0 {
		s.PassRate = float64(s.Passed) / float64(s.Total)
	}
	return s
}

// Write saves grading.json.
func Write(path string, r Report) error {
	r.Summary = Summarize(r.AssertionResults)
	raw, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	raw = append(raw, '\n')
	return os.WriteFile(path, raw, 0o644)
}

// Load reads grading.json.
func Load(path string) (Report, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Report{}, err
	}
	var r Report
	if err := json.Unmarshal(raw, &r); err != nil {
		return Report{}, err
	}
	return r, nil
}

// ListOutputFiles returns relative paths under outputsDir (recursive).
func ListOutputFiles(outputsDir string) ([]string, string, error) {
	var files []string
	var excerpts strings.Builder
	err := filepath.Walk(outputsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(outputsDir, path)
		files = append(files, rel)
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		content := string(data)
		if len(content) > 4000 {
			content = content[:4000] + "\n...[truncated]...\n"
		}
		excerpts.WriteString(fmt.Sprintf("=== %s ===\n%s\n\n", rel, content))
		return nil
	})
	return files, excerpts.String(), err
}

// JudgePrompt builds the LLM grading prompt.
func JudgePrompt(expected string, assertions []string, fileList []string, excerpts, transcriptExcerpt string) string {
	var b strings.Builder
	b.WriteString("You are grading an agent skill eval. For EACH assertion, decide PASS or FAIL with concrete evidence quoting or referencing the outputs.\n")
	b.WriteString("Require concrete evidence for PASS. Do not give benefit of the doubt.\n")
	b.WriteString("Respond with ONLY a JSON object matching:\n")
	b.WriteString(`{"assertion_results":[{"text":"...","passed":true,"evidence":"..."}]` + "}\n\n")
	b.WriteString("Expected output description:\n")
	b.WriteString(expected)
	b.WriteString("\n\nAssertions:\n")
	for i, a := range assertions {
		b.WriteString(fmt.Sprintf("%d. %s\n", i+1, a))
	}
	b.WriteString("\nOutput files:\n")
	if len(fileList) == 0 {
		b.WriteString("(none)\n")
	} else {
		for _, f := range fileList {
			b.WriteString("- " + f + "\n")
		}
	}
	b.WriteString("\nOutput excerpts:\n")
	b.WriteString(excerpts)
	if transcriptExcerpt != "" {
		b.WriteString("\nTranscript excerpt:\n")
		b.WriteString(transcriptExcerpt)
	}
	return b.String()
}

// ParseJudgeJSON extracts Report assertion_results from LLM stdout.
func ParseJudgeJSON(stdout string) ([]AssertionResult, error) {
	stdout = strings.TrimSpace(stdout)
	// try whole string
	var wrap struct {
		AssertionResults []AssertionResult `json:"assertion_results"`
	}
	if err := json.Unmarshal([]byte(stdout), &wrap); err == nil && len(wrap.AssertionResults) > 0 {
		return wrap.AssertionResults, nil
	}
	// find first { ... }
	start := strings.Index(stdout, "{")
	end := strings.LastIndex(stdout, "}")
	if start >= 0 && end > start {
		if err := json.Unmarshal([]byte(stdout[start:end+1]), &wrap); err == nil {
			return wrap.AssertionResults, nil
		}
	}
	return nil, fmt.Errorf("could not parse assertion_results JSON from judge output")
}

// RunJudge invokes provider CLI as a judge (same binaries as runners).
func RunJudge(ctx context.Context, provider, model, prompt string, timeout time.Duration) (string, error) {
	if timeout <= 0 {
		timeout = 5 * time.Minute
	}
	runCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var bin string
	var args []string
	switch strings.ToLower(provider) {
	case "opencode":
		bin = "opencode"
		args = []string{"run", "--model", model, prompt}
	case "claude":
		bin = "claude"
		args = []string{"-p", prompt}
		if model != "" {
			args = append(args, "--model", model)
		}
	case "codex":
		bin = "codex"
		args = []string{"exec"}
		if model != "" {
			args = append(args, "--model", model)
		}
		args = append(args, prompt)
	case "agent", "cursor_agent", "cursor-agent":
		bin = "agent"
		args = []string{"-p", prompt}
		if model != "" {
			args = append(args, "--model", model)
		}
	default:
		return "", fmt.Errorf("unknown provider %q", provider)
	}

	cmd := exec.CommandContext(runCtx, bin, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	out := stdout.String()
	if out == "" {
		out = stderr.String()
	}
	if err != nil && out == "" {
		return "", fmt.Errorf("judge %s: %w: %s", bin, err, stderr.String())
	}
	return out, nil
}

// EmptyFailReport builds a FAIL-all report when no assertions or judge fails.
func EmptyFailReport(assertions []string, reason string) Report {
	var results []AssertionResult
	for _, a := range assertions {
		results = append(results, AssertionResult{
			Text:     a,
			Passed:   false,
			Evidence: reason,
		})
	}
	return Report{
		AssertionResults: results,
		Summary:          Summarize(results),
	}
}
