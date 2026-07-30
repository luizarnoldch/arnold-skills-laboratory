package eval

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"skills-laboratory/lab-go/internal/logdir"
	"skills-laboratory/lab-go/internal/prompt"
	"skills-laboratory/lab-go/internal/provider"
	"skills-laboratory/lab-go/internal/result"
)

// Config configures a full prompt-set evaluation.
type Config struct {
	SkillName         string
	PromptsPath       string
	ProviderName      string
	Model             string
	Workdir           string
	OutPath           string
	Runs              int
	Timeout           time.Duration
	MajorityThreshold float64
	LogDirBase        string // empty => <workdir>/../logs
	NoLogs            bool
}

// Run evaluates all prompts and writes the payload JSON to OutPath.
func Run(cfg Config) (result.Payload, error) {
	if cfg.Runs < 1 {
		cfg.Runs = 1
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 60 * time.Second
	}
	if cfg.MajorityThreshold <= 0 {
		cfg.MajorityThreshold = 0.5
	}

	items, err := prompt.Load(cfg.PromptsPath)
	if err != nil {
		return result.Payload{}, err
	}

	workdir, err := filepath.Abs(cfg.Workdir)
	if err != nil {
		return result.Payload{}, err
	}
	info, err := os.Stat(workdir)
	if err != nil || !info.IsDir() {
		return result.Payload{}, fmt.Errorf("workdir does not exist: %s", workdir)
	}

	prov, err := provider.New(cfg.ProviderName, provider.Config{
		Model:     cfg.Model,
		SkillName: cfg.SkillName,
		Workdir:   workdir,
		Timeout:   cfg.Timeout,
	})
	if err != nil {
		return result.Payload{}, err
	}

	var logDir string
	if !cfg.NoLogs {
		base := cfg.LogDirBase
		if base == "" {
			base = filepath.Join(filepath.Dir(workdir), "logs")
		}
		base, err = filepath.Abs(base)
		if err != nil {
			return result.Payload{}, err
		}
		logDir, err = logdir.NextRunDir(base)
		if err != nil {
			return result.Payload{}, err
		}
	}

	promptsAbs, _ := filepath.Abs(cfg.PromptsPath)

	fmt.Println("==================================================")
	fmt.Printf("Evaluating with provider=%s model=%s\n", prov.Name(), cfg.Model)
	fmt.Printf("Skill: %s | Default runs: %d | Timeout: %s\n", cfg.SkillName, cfg.Runs, cfg.Timeout)
	fmt.Printf("Prompts: %s (%d items)\n", cfg.PromptsPath, len(items))
	if logDir != "" {
		fmt.Printf("Log directory: %s\n", logDir)
	}
	fmt.Println("==================================================")

	started := time.Now()
	results := make([]result.Item, 0, len(items))

	for _, item := range items {
		runCount := cfg.Runs
		if item.Runs != nil && *item.Runs > 0 {
			runCount = *item.Runs
		}
		fmt.Printf("\n[ID: %d] Query: %q | Expected: %t | Runs: %d\n",
			item.ID, item.Query, item.ShouldTrigger, runCount)

		ri, err := evaluatePrompt(prov, item, runCount, logDir, cfg.MajorityThreshold)
		if err != nil {
			return result.Payload{}, err
		}
		results = append(results, ri)
	}

	summary := result.BuildSummary(results)
	payload := result.Payload{
		SkillName:         cfg.SkillName,
		Provider:          prov.Name(),
		Model:             cfg.Model,
		PromptsFile:       promptsAbs,
		Workdir:           workdir,
		RunsDefault:       cfg.Runs,
		MajorityThreshold: cfg.MajorityThreshold,
		Summary:           summary,
		ElapsedSeconds:    float64(time.Since(started).Milliseconds()) / 1000.0,
		Results:           results,
	}

	if err := result.WriteJSON(cfg.OutPath, payload); err != nil {
		return result.Payload{}, err
	}

	fmt.Printf("\nAccuracy: %d/%d (%.2f%%)\n", summary.Correct, summary.Total, summary.AccuracyPct)
	fmt.Printf("Results saved to %s\n", cfg.OutPath)
	if logDir != "" {
		fmt.Printf("Logs saved in %s\n", logDir)
	}
	return payload, nil
}

func evaluatePrompt(
	prov provider.Provider,
	item prompt.Item,
	runs int,
	logDir string,
	majorityThreshold float64,
) (result.Item, error) {
	triggerCount := 0
	for i := 1; i <= runs; i++ {
		rr, err := prov.Run(context.Background(), item.Query)
		if err != nil {
			return result.Item{}, fmt.Errorf("id %d run %d: %w", item.ID, i, err)
		}
		status := "NOT TRIGGERED"
		if rr.Triggered {
			triggerCount++
			status = "TRIGGERED"
		} else if rr.TimedOut {
			status = "TIMED OUT"
		}
		fmt.Printf("  └─ Run %d/%d: %s\n", i, runs, status)

		if logDir != "" {
			name := fmt.Sprintf("id_%d.log", item.ID)
			if runs > 1 {
				name = fmt.Sprintf("id_%d_run_%d.log", item.ID, i)
			}
			header := fmt.Sprintf(
				"=== RUN METADATA ===\nID: %d\nTimestamp: %s\nQuery: %s\nModel: %s\nProvider: %s\nRun: %d/%d\nExpected Trigger: %t\nTriggered: %t\nTimed Out: %t\n====================\n\n",
				item.ID, time.Now().UTC().Format(time.RFC3339), item.Query, prov.Model(), prov.Name(),
				i, runs, item.ShouldTrigger, rr.Triggered, rr.TimedOut,
			)
			if err := os.WriteFile(filepath.Join(logDir, name), []byte(header+rr.Stdout), 0o644); err != nil {
				return result.Item{}, err
			}
		}
	}

	rate := result.RoundRate(triggerCount, runs)
	return result.Item{
		ID:            item.ID,
		Query:         item.Query,
		ShouldTrigger: item.ShouldTrigger,
		SkillName:     prov.SkillName(),
		Triggers:      triggerCount,
		Runs:          runs,
		TriggerRate:   rate,
		Correct:       result.IsCorrect(rate, item.ShouldTrigger, majorityThreshold),
	}, nil
}
