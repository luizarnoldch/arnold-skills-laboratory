package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"skills-laboratory/lab-go-quality/internal/evalset"
	"skills-laboratory/lab-go-quality/internal/runner"
	"skills-laboratory/lab-go-quality/internal/timing"
	"skills-laboratory/lab-go-quality/internal/workspace"
)

func main() {
	evalsPath := flag.String("evals", "", "Path to evals.json")
	skillPath := flag.String("skill-path", "", "Skill directory or SKILL.md")
	skillName := flag.String("skill-name", "", "Skill name (default: from evals.json)")
	ws := flag.String("workspace", "", "Quality workspace root (workspace/quality/<name>)")
	providerName := flag.String("provider", "opencode", "opencode | claude | codex | agent")
	model := flag.String("model", "digitalocean/deepseek-4-flash", "Model id")
	baseline := flag.String("baseline", "none", "none (=without_skill) | snapshot (=old_skill)")
	runs := flag.Int("runs", 1, "Runs per eval × config")
	timeout := flag.Float64("timeout", 600, "Per-run timeout seconds")
	flag.Parse()

	if *evalsPath == "" || *skillPath == "" || *ws == "" {
		fmt.Fprintln(os.Stderr, "error: -evals, -skill-path, -workspace are required")
		flag.Usage()
		os.Exit(2)
	}

	set, err := evalset.Load(*evalsPath)
	must(err)
	name := *skillName
	if name == "" {
		name = set.SkillName
	}

	skillDir, err := runner.ResolveSkillDir(*skillPath)
	must(err)

	evalsAbs, err := filepath.Abs(*evalsPath)
	must(err)
	evalsRoot := filepath.Dir(evalsAbs) // .../evals

	tr, err := runner.New(*providerName)
	must(err)

	iterDir, err := workspace.NextIterationDir(*ws)
	must(err)
	fmt.Printf("iteration: %s\n", iterDir)

	snapDir := filepath.Join(iterDir, "skill-snapshot")
	must(workspace.CopyTree(skillDir, snapDir))

	to := time.Duration(*timeout * float64(time.Second))
	ctx := context.Background()

	configs := []string{"with_skill"}
	switch strings.ToLower(*baseline) {
	case "none", "without", "without_skill":
		configs = append(configs, "without_skill")
	case "snapshot", "old", "old_skill":
		configs = append(configs, "old_skill")
	default:
		fatal(fmt.Errorf("unknown -baseline %q (want none|snapshot)", *baseline))
	}

	for _, c := range set.Evals {
		slug := evalset.Slug(c)
		evalDir := filepath.Join(iterDir, slug)
		must(workspace.EnsureDir(evalDir))

		inputAbs := resolveInputs(evalsRoot, c.Files)

		for _, cfg := range configs {
			configDir := filepath.Join(evalDir, cfg)
			for runN := 1; runN <= *runs; runN++ {
				runDir, err := workspace.NextRunDir(configDir)
				must(err)
				outputsDir := filepath.Join(runDir, "outputs")
				sandbox := filepath.Join(runDir, "sandbox")
				must(workspace.EnsureDir(sandbox))

				copied, err := runner.CopyInputs(sandbox, inputAbs)
				must(err)

				var useSkill string
				switch cfg {
				case "with_skill":
					useSkill = skillDir
					must(tr.InstallSkill(sandbox, name, skillDir))
				case "old_skill":
					useSkill = snapDir
					must(tr.InstallSkill(sandbox, name, snapDir))
				case "without_skill":
					useSkill = ""
				}

				outAbs, _ := filepath.Abs(outputsDir)
				req := runner.TaskRequest{
					Prompt:     c.Prompt,
					SkillName:  name,
					SkillPath:  useSkill,
					InputFiles: copied,
					OutputDir:  outAbs,
					Workdir:    sandbox,
					Model:      *model,
					Timeout:    to,
				}

				fmt.Printf("  [%s/%s run=%d] provider=%s\n", slug, cfg, runN, tr.Name())
				res, err := tr.Run(ctx, req)
				if err != nil {
					fmt.Fprintf(os.Stderr, "warning: run failed: %v\n", err)
				}
				must(os.WriteFile(filepath.Join(runDir, "transcript.log"), []byte(res.Transcript), 0o644))

				rel := workspace.RelPath(iterDir, runDir)
				entry := timing.Entry{
					EvalID:      c.ID,
					EvalSlug:    slug,
					Config:      cfg,
					Run:         runN,
					Path:        rel,
					TotalTokens: res.TotalTokens,
					DurationMS:  res.DurationMS,
					Provider:    tr.Name(),
					ExitCode:    res.ExitCode,
					TimedOut:    res.TimedOut,
				}
				entry, err = timing.Append(iterDir, entry)
				must(err)
				must(timing.WriteRunTiming(runDir, entry))
				fmt.Printf("    index=%d path=%s duration_ms=%d tokens=%d\n",
					entry.Index, entry.Path, entry.DurationMS, entry.TotalTokens)
			}
		}
	}

	fmt.Printf("\ndone. timing stack: %s\n", timing.Path(iterDir))
}

func resolveInputs(evalsRoot string, files []string) []string {
	var out []string
	for _, f := range files {
		f = strings.TrimSpace(f)
		if f == "" {
			continue
		}
		if filepath.IsAbs(f) {
			out = append(out, f)
			continue
		}
		// paths in evals.json are relative to skill dir (evals/files/...)
		// evalsRoot is .../evals, skill dir is parent
		skillDir := filepath.Dir(evalsRoot)
		p := filepath.Join(skillDir, f)
		if _, err := os.Stat(p); err == nil {
			out = append(out, p)
			continue
		}
		// also try relative to evals root
		p2 := filepath.Join(evalsRoot, f)
		out = append(out, p2)
	}
	return out
}

func must(err error) {
	if err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}
