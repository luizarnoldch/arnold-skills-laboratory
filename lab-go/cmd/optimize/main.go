package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"skills-laboratory/lab-go/internal/eval"
	"skills-laboratory/lab-go/internal/logdir"
	"skills-laboratory/lab-go/internal/result"
	"skills-laboratory/lab-go/internal/skillmd"
)

func main() {
	skillName := flag.String("skill-name", "", "Skill under test")
	skillMD := flag.String("skill-md", "", "Path to SKILL.md to optimize")
	prompts := flag.String("prompts", "", "Train prompts JSON only")
	workdir := flag.String("workdir", "", "Sandbox cwd")
	iterationsDir := flag.String("iterations-dir", "", "Directory for iteration snapshots")
	resultsDir := flag.String("results-dir", "", "Directory for train result JSONs")
	providerName := flag.String("provider", "opencode", "Eval provider")
	model := flag.String("model", "digitalocean/deepseek-4-flash", "Model id")
	runs := flag.Int("runs", 3, "Runs per prompt")
	timeout := flag.Float64("timeout", 60, "Per-run timeout seconds")
	threshold := flag.Float64("threshold", 0.95, "Stop when train accuracy >= this")
	majority := flag.Float64("majority-threshold", 0.5, "Majority trigger threshold")
	maxIters := flag.Int("max-iters", 5, "Maximum optimize iterations")
	logDir := flag.String("log-dir", "", "Base log directory")
	noLogs := flag.Bool("no-logs", false, "Disable per-run logs")
	dryRun := flag.Bool("dry-run", false, "Evaluate once; do not rewrite description")
	flag.Parse()

	if *skillName == "" || *skillMD == "" || *prompts == "" || *workdir == "" || *iterationsDir == "" || *resultsDir == "" {
		fmt.Fprintln(os.Stderr, "error: -skill-name, -skill-md, -prompts, -workdir, -iterations-dir, -results-dir are required")
		flag.Usage()
		os.Exit(2)
	}

	if strings.Contains(strings.ToLower(filepath.Base(*prompts)), "validation") {
		fmt.Fprintln(os.Stderr, "WARNING: prompts path looks like validation data. optimize must only use train.json.")
	}

	skillPath, err := filepath.Abs(*skillMD)
	if err != nil {
		fatal(err)
	}
	if err := os.MkdirAll(*resultsDir, 0o755); err != nil {
		fatal(err)
	}

	to := time.Duration(*timeout * float64(time.Second))

	for iteration := 1; iteration <= *maxIters; iteration++ {
		fmt.Printf("\n######## TRAIN ITERATION %d/%d ########\n", iteration, *maxIters)
		outPath := filepath.Join(*resultsDir, fmt.Sprintf("iter_%03d.json", iteration))

		payload, err := eval.Run(eval.Config{
			SkillName:         *skillName,
			PromptsPath:       *prompts,
			ProviderName:      *providerName,
			Model:             *model,
			Workdir:           *workdir,
			OutPath:           outPath,
			Runs:              *runs,
			Timeout:           to,
			MajorityThreshold: *majority,
			LogDirBase:        *logDir,
			NoLogs:            *noLogs,
		})
		if err != nil {
			fatal(err)
		}

		accuracy := payload.Summary.Accuracy
		fails := result.Failures(payload)

		parts, err := skillmd.Read(skillPath)
		if err != nil {
			fatal(err)
		}
		currentDesc, err := skillmd.GetDescription(parts.Frontmatter)
		if err != nil {
			fatal(err)
		}

		iterDir, err := logdir.NextIterationDir(*iterationsDir)
		if err != nil {
			fatal(err)
		}
		_ = os.WriteFile(filepath.Join(iterDir, "SKILL.description.md"),
			[]byte("# description snapshot\n\n"+currentDesc+"\n"), 0o644)

		decision := "optimize"
		if accuracy >= *threshold {
			decision = "stop_threshold_met"
		}
		if *dryRun {
			decision = "dry_run_stop"
		}
		if len(fails) == 0 && accuracy >= *threshold {
			decision = "stop_perfect"
		}

		metrics := map[string]any{
			"iteration":     iteration,
			"timestamp":     time.Now().UTC().Format(time.RFC3339),
			"accuracy":      accuracy,
			"accuracy_pct":  payload.Summary.AccuracyPct,
			"threshold":     *threshold,
			"failures":      fails,
			"decision":      decision,
			"results_file":  outPath,
		}
		if err := writeJSON(filepath.Join(iterDir, "metrics.json"), metrics); err != nil {
			fatal(err)
		}

		fmt.Printf("Train accuracy: %.2f%% | failures=%d | decision=%s\n",
			payload.Summary.AccuracyPct, len(fails), decision)

		if strings.HasPrefix(decision, "stop") || decision == "dry_run_stop" {
			fmt.Println("Stopping train loop.")
			return
		}

		if *providerName != "opencode" {
			fmt.Fprintf(os.Stderr, "Auto-optimize uses opencode for rewriting descriptions; provider=%s will still evaluate but rewrite requires opencode.\n", *providerName)
		}

		workdirAbs, _ := filepath.Abs(*workdir)
		newDesc, err := proposeDescription(workdirAbs, *model, *skillName, currentDesc, fails, maxDuration(to, 120*time.Second))
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to propose new description:", err)
			metrics["decision"] = "optimize_failed"
			_ = writeJSON(filepath.Join(iterDir, "metrics.json"), metrics)
			os.Exit(1)
		}

		parts.Frontmatter = skillmd.SetDescription(parts.Frontmatter, newDesc)
		if err := skillmd.Write(skillPath, parts); err != nil {
			fatal(err)
		}
		_ = os.WriteFile(filepath.Join(iterDir, "proposed_description.md"), []byte(newDesc+"\n"), 0o644)
		fmt.Printf("Updated description in %s\n", skillPath)
		fmt.Printf("Snapshot: %s\n", iterDir)
	}

	fmt.Println("Reached max_iters without meeting threshold.")
}

func proposeDescription(workdir, model, skillName, current string, fails []result.Item, timeout time.Duration) (string, error) {
	var failLines []string
	for _, f := range fails {
		kind := "false_positive"
		if f.ShouldTrigger {
			kind = "false_negative"
		}
		failLines = append(failLines, fmt.Sprintf("- id=%d (%s) trigger_rate=%v: %q", f.ID, kind, f.TriggerRate, f.Query))
	}
	failBlock := "(none)"
	if len(failLines) > 0 {
		failBlock = strings.Join(failLines, "\n")
	}

	prompt := fmt.Sprintf(`You are optimizing the YAML frontmatter `+"`description`"+` of an agent skill named %q.

Current description:
"""%s"""

These train prompts were classified incorrectly (skill trigger vs expected):
%s

Rewrite ONLY the description text so the model more reliably triggers on should_trigger=true
prompts and does NOT trigger on should_trigger=false prompts.

Rules:
- Return ONLY the new description plain text (no YAML, no markdown fences, no quotes wrapper).
- Keep it to 1-3 sentences / one folded paragraph.
- Include clear positive triggers AND explicit Do NOT / Does NOT negatives when helpful.
`, skillName, current, failBlock)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "opencode", "run", "--model", model, prompt)
	cmd.Dir = workdir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil && ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("description optimization timed out")
	}

	text := stdout.String() + "\n" + stderr.String()
	var lines []string
	for _, ln := range strings.Split(text, "\n") {
		ln = strings.TrimSpace(ln)
		if ln != "" {
			lines = append(lines, ln)
		}
	}
	if len(lines) == 0 {
		return "", fmt.Errorf("optimizer returned empty output")
	}

	var cleaned []string
	for _, ln := range lines {
		if strings.HasPrefix(ln, "Skill ") || strings.Contains(ln, `"name":`) || strings.HasPrefix(ln, "===") {
			continue
		}
		cleaned = append(cleaned, ln)
	}
	src := cleaned
	if len(src) == 0 {
		src = lines
	}
	start := 0
	if len(src) > 3 {
		start = len(src) - 3
	}
	candidate := strings.TrimSpace(strings.Join(src[start:], " "))
	candidate = strings.Trim(candidate, "`\"'")
	if len(candidate) < 40 {
		return "", fmt.Errorf("optimizer output too short: %q", candidate)
	}
	return candidate, nil
}

func writeJSON(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}

func maxDuration(a, b time.Duration) time.Duration {
	if a > b {
		return a
	}
	return b
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}
