package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"skills-laboratory/lab-go/internal/eval"
)

func main() {
	skillName := flag.String("skill-name", "", "Skill name to detect")
	prompts := flag.String("prompts", "", "Prompt JSON (train or validation)")
	providerName := flag.String("provider", "opencode", "CLI provider: opencode|codex|claude|cursor_agent")
	runs := flag.Int("runs", 1, "Runs per prompt (default 1)")
	model := flag.String("model", "digitalocean/deepseek-4-flash", "Model id for the CLI")
	workdir := flag.String("workdir", "", "Sandbox cwd for the CLI")
	out := flag.String("out", "", "Results JSON output path")
	timeout := flag.Float64("timeout", 60, "Per-run timeout seconds")
	majority := flag.Float64("majority-threshold", 0.5, "trigger_rate >= threshold counts as triggered")
	logDir := flag.String("log-dir", "", "Base log directory (creates run_N). Default: <workdir>/../logs")
	noLogs := flag.Bool("no-logs", false, "Disable per-run log files")
	flag.Parse()

	if *skillName == "" || *prompts == "" || *workdir == "" || *out == "" {
		fmt.Fprintln(os.Stderr, "error: -skill-name, -prompts, -workdir, and -out are required")
		flag.Usage()
		os.Exit(2)
	}

	_, err := eval.Run(eval.Config{
		SkillName:         *skillName,
		PromptsPath:       *prompts,
		ProviderName:      *providerName,
		Model:             *model,
		Workdir:           *workdir,
		OutPath:           *out,
		Runs:              *runs,
		Timeout:           time.Duration(*timeout * float64(time.Second)),
		MajorityThreshold: *majority,
		LogDirBase:        *logDir,
		NoLogs:            *noLogs,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
