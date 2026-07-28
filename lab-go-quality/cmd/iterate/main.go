package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"skills-laboratory/lab-go-quality/internal/grade"
	"skills-laboratory/lab-go-quality/internal/skillmd"
	"skills-laboratory/lab-go-quality/internal/timing"
)

func main() {
	iteration := flag.String("iteration", "", "Path to iteration-N")
	skillMD := flag.String("skill-md", "", "Path to canonical SKILL.md")
	providerName := flag.String("provider", "opencode", "Provider for proposing edits")
	model := flag.String("model", "digitalocean/deepseek-4-flash", "Model id")
	timeout := flag.Float64("timeout", 300, "Timeout seconds")
	apply := flag.Bool("apply", false, "Overwrite SKILL.md with the proposal")
	flag.Parse()

	if *iteration == "" || *skillMD == "" {
		fmt.Fprintln(os.Stderr, "error: -iteration and -skill-md are required")
		flag.Usage()
		os.Exit(2)
	}

	current, err := skillmd.ReadRaw(*skillMD)
	must(err)

	stack, err := timing.Load(*iteration)
	must(err)

	var failed strings.Builder
	var transcripts strings.Builder
	for _, e := range stack.Runs {
		if e.Config != "with_skill" {
			continue
		}
		gPath := filepath.Join(*iteration, e.Path, "grading.json")
		r, err := grade.Load(gPath)
		if err != nil {
			continue
		}
		for _, a := range r.AssertionResults {
			if a.Passed {
				continue
			}
			failed.WriteString(fmt.Sprintf("- [index=%d %s] %s\n  evidence: %s\n",
				e.Index, e.Path, a.Text, a.Evidence))
		}
		tPath := filepath.Join(*iteration, e.Path, "transcript.log")
		if raw, err := os.ReadFile(tPath); err == nil {
			excerpt := string(raw)
			if len(excerpt) > 2500 {
				excerpt = excerpt[:2500] + "\n...[truncated]...\n"
			}
			transcripts.WriteString(fmt.Sprintf("=== index=%d %s ===\n%s\n\n", e.Index, e.Path, excerpt))
		}
	}

	feedback := ""
	if raw, err := os.ReadFile(filepath.Join(*iteration, "feedback.json")); err == nil {
		feedback = string(raw)
	}

	prompt := fmt.Sprintf(`You improve agent skills using eval signals.

Guidelines:
- Generalize from feedback; do not patch only the test cases.
- Keep the skill lean; fewer better instructions beat exhaustive rules.
- Prefer reasoning-based instructions ("Do X because Y") over rigid ALWAYS/NEVER.
- Bundle repeated helper work into scripts/ when transcripts show duplication.

Return ONLY the full updated SKILL.md contents (YAML frontmatter + body). No markdown fences.

Current SKILL.md:
%s

Failed assertions (with_skill):
%s

Human feedback.json:
%s

Transcript excerpts:
%s
`, current, empty(failed.String(), "(none)"), empty(feedback, "{}"), empty(transcripts.String(), "(none)"))

	to := time.Duration(*timeout * float64(time.Second))
	out, err := runPropose(context.Background(), *providerName, *model, prompt, to)
	must(err)

	proposed := stripFences(out)
	outPath := filepath.Join(*iteration, "proposed-SKILL.md")
	must(os.WriteFile(outPath, []byte(proposed), 0o644))
	fmt.Printf("wrote %s\n", outPath)

	if *apply {
		must(skillmd.WriteRaw(*skillMD, proposed))
		fmt.Printf("applied to %s\n", *skillMD)
	} else {
		fmt.Println("review the proposal; re-run with -apply to overwrite the skill")
	}
}

func runPropose(ctx context.Context, provider, model, prompt string, timeout time.Duration) (string, error) {
	runCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	var bin string
	var args []string
	switch strings.ToLower(provider) {
	case "opencode":
		bin, args = "opencode", []string{"run", "--model", model, prompt}
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
	if err != nil && strings.TrimSpace(out) == "" {
		return "", fmt.Errorf("%s: %w: %s", bin, err, stderr.String())
	}
	return out, nil
}

func stripFences(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "```") {
		s = strings.TrimPrefix(s, "```markdown")
		s = strings.TrimPrefix(s, "```md")
		s = strings.TrimPrefix(s, "```")
		s = strings.TrimSpace(s)
		if i := strings.LastIndex(s, "```"); i >= 0 {
			s = strings.TrimSpace(s[:i])
		}
	}
	return s
}

func empty(s, fallback string) string {
	if strings.TrimSpace(s) == "" {
		return fallback
	}
	return s
}

func must(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
