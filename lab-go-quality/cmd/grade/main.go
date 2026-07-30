package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"skills-laboratory/lab-go-quality/internal/evalset"
	"skills-laboratory/lab-go-quality/internal/grade"
	"skills-laboratory/lab-go-quality/internal/timing"
)

func main() {
	iteration := flag.String("iteration", "", "Path to iteration-N directory")
	evalsPath := flag.String("evals", "", "Path to evals.json (for assertions)")
	providerName := flag.String("provider", "opencode", "Judge provider")
	model := flag.String("model", "digitalocean/deepseek-4-flash", "Model id")
	timeout := flag.Float64("timeout", 300, "Judge timeout seconds")
	indexFilter := flag.Int("index", 0, "Grade only this timing index (0 = all)")
	flag.Parse()

	if *iteration == "" || *evalsPath == "" {
		fmt.Fprintln(os.Stderr, "error: -iteration and -evals are required")
		flag.Usage()
		os.Exit(2)
	}

	set, err := evalset.Load(*evalsPath)
	must(err)
	byID := map[int]evalset.Case{}
	for _, c := range set.Evals {
		byID[c.ID] = c
	}

	stack, err := timing.Load(*iteration)
	must(err)
	if len(stack.Runs) == 0 {
		fmt.Fprintln(os.Stderr, "error: timing stack is empty; run runevals first")
		os.Exit(1)
	}

	to := time.Duration(*timeout * float64(time.Second))
	ctx := context.Background()

	for _, e := range stack.Runs {
		if *indexFilter > 0 && e.Index != *indexFilter {
			continue
		}
		c, ok := byID[e.EvalID]
		if !ok {
			fmt.Fprintf(os.Stderr, "warning: no eval id=%d in evals.json; skip index=%d\n", e.EvalID, e.Index)
			continue
		}
		if len(c.Assertions) == 0 {
			fmt.Printf("  index=%d: no assertions; writing empty grading\n", e.Index)
			rep := grade.Report{
				Index:    e.Index,
				EvalID:   e.EvalID,
				EvalSlug: e.EvalSlug,
				Config:   e.Config,
				Run:      e.Run,
				Path:     e.Path,
			}
			must(grade.Write(filepath.Join(*iteration, e.Path, "grading.json"), rep))
			continue
		}

		runDir := filepath.Join(*iteration, e.Path)
		outputsDir := filepath.Join(runDir, "outputs")
		files, excerpts, err := grade.ListOutputFiles(outputsDir)
		if err != nil && !os.IsNotExist(err) {
			must(err)
		}
		transcript := ""
		if raw, err := os.ReadFile(filepath.Join(runDir, "transcript.log")); err == nil {
			transcript = string(raw)
			if len(transcript) > 6000 {
				transcript = transcript[:6000] + "\n...[truncated]...\n"
			}
		}

		prompt := grade.JudgePrompt(c.ExpectedOutput, c.Assertions, files, excerpts, transcript)
		fmt.Printf("  grading index=%d %s/%s run=%d\n", e.Index, e.EvalSlug, e.Config, e.Run)
		out, err := grade.RunJudge(ctx, *providerName, *model, prompt, to)
		var results []grade.AssertionResult
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: judge failed index=%d: %v\n", e.Index, err)
			rep := grade.EmptyFailReport(c.Assertions, "judge failed: "+err.Error())
			results = rep.AssertionResults
		} else {
			results, err = grade.ParseJudgeJSON(out)
			if err != nil {
				fmt.Fprintf(os.Stderr, "warning: parse judge JSON index=%d: %v\n", e.Index, err)
				_ = os.WriteFile(filepath.Join(runDir, "judge_raw.txt"), []byte(out), 0o644)
				rep := grade.EmptyFailReport(c.Assertions, "could not parse judge JSON: "+err.Error())
				results = rep.AssertionResults
			}
		}

		rep := grade.Report{
			Index:            e.Index,
			EvalID:           e.EvalID,
			EvalSlug:         e.EvalSlug,
			Config:           e.Config,
			Run:              e.Run,
			Path:             e.Path,
			AssertionResults: results,
		}
		must(grade.Write(filepath.Join(runDir, "grading.json"), rep))
		fmt.Printf("    pass_rate=%.2f (%d/%d)\n", rep.Summary.PassRate, rep.Summary.Passed, rep.Summary.Total)
	}

	// write empty feedback.json template if missing
	fb := filepath.Join(*iteration, "feedback.json")
	if _, err := os.Stat(fb); os.IsNotExist(err) {
		m := map[string]string{}
		for _, c := range set.Evals {
			m[evalset.Slug(c)] = ""
		}
		raw, _ := json.MarshalIndent(m, "", "  ")
		raw = append(raw, '\n')
		_ = os.WriteFile(fb, raw, 0o644)
	}
	fmt.Println("done.")
}

func must(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
