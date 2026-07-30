package benchmark

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"

	"skills-laboratory/lab-go-quality/internal/grade"
	"skills-laboratory/lab-go-quality/internal/timing"
)

// Metric is mean + stddev.
type Metric struct {
	Mean   float64 `json:"mean"`
	Stddev float64 `json:"stddev"`
}

// ConfigStats aggregates one configuration (with_skill / without_skill / old_skill).
type ConfigStats struct {
	PassRate    Metric `json:"pass_rate"`
	TimeSeconds Metric `json:"time_seconds"`
	Tokens      Metric `json:"tokens"`
}

// Delta is with_skill - baseline (without_skill preferred, else old_skill).
type Delta struct {
	PassRate    float64 `json:"pass_rate"`
	TimeSeconds float64 `json:"time_seconds"`
	Tokens      float64 `json:"tokens"`
}

// RunSummary is the benchmark payload.
// WithoutSkill is always present (may be zero when baseline was old_skill).
// OldSkill is set when snapshot baseline runs exist.
// Delta is with_skill minus without_skill if that config has data, else old_skill.
type RunSummary struct {
	WithSkill    ConfigStats  `json:"with_skill"`
	WithoutSkill ConfigStats  `json:"without_skill"`
	OldSkill     *ConfigStats `json:"old_skill,omitempty"`
	Delta        Delta        `json:"delta"`
	Baseline     string       `json:"baseline"` // "without_skill" | "old_skill" | "none"
}

// Report is benchmark.json.
type Report struct {
	RunSummary RunSummary `json:"run_summary"`
}

// Compute builds benchmark from iteration timing stack + grading.json files.
func Compute(iterationDir string) (Report, error) {
	stack, err := timing.Load(iterationDir)
	if err != nil {
		return Report{}, err
	}

	passRates := map[string][]float64{}
	times := map[string][]float64{}
	tokens := map[string][]float64{}

	for _, e := range stack.Runs {
		times[e.Config] = append(times[e.Config], float64(e.DurationMS)/1000.0)
		tokens[e.Config] = append(tokens[e.Config], float64(e.TotalTokens))

		gPath := filepath.Join(iterationDir, e.Path, "grading.json")
		if r, err := grade.Load(gPath); err == nil {
			passRates[e.Config] = append(passRates[e.Config], r.Summary.PassRate)
		}
	}

	with := statsFor("with_skill", passRates, times, tokens)
	without := statsFor("without_skill", passRates, times, tokens)
	old := statsFor("old_skill", passRates, times, tokens)

	baselineName := "none"
	baseline := ConfigStats{}
	switch {
	case len(passRates["without_skill"]) > 0 || len(times["without_skill"]) > 0:
		baselineName = "without_skill"
		baseline = without
	case len(passRates["old_skill"]) > 0 || len(times["old_skill"]) > 0:
		baselineName = "old_skill"
		baseline = old
	}

	rep := Report{
		RunSummary: RunSummary{
			WithSkill:    with,
			WithoutSkill: without,
			Baseline:     baselineName,
			Delta: Delta{
				PassRate:    with.PassRate.Mean - baseline.PassRate.Mean,
				TimeSeconds: with.TimeSeconds.Mean - baseline.TimeSeconds.Mean,
				Tokens:      with.Tokens.Mean - baseline.Tokens.Mean,
			},
		},
	}
	if len(passRates["old_skill"]) > 0 || len(times["old_skill"]) > 0 {
		o := old
		rep.RunSummary.OldSkill = &o
	}
	return rep, nil
}

func statsFor(config string, passRates, times, tokens map[string][]float64) ConfigStats {
	return ConfigStats{
		PassRate:    metric(passRates[config]),
		TimeSeconds: metric(times[config]),
		Tokens:      metric(tokens[config]),
	}
}

func metric(vals []float64) Metric {
	if len(vals) == 0 {
		return Metric{}
	}
	var sum float64
	for _, v := range vals {
		sum += v
	}
	mean := sum / float64(len(vals))
	if len(vals) == 1 {
		return Metric{Mean: mean, Stddev: 0}
	}
	var varSum float64
	for _, v := range vals {
		d := v - mean
		varSum += d * d
	}
	return Metric{Mean: mean, Stddev: math.Sqrt(varSum / float64(len(vals)))}
}

// Write saves benchmark.json.
func Write(iterationDir string, r Report) error {
	raw, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	raw = append(raw, '\n')
	return os.WriteFile(filepath.Join(iterationDir, "benchmark.json"), raw, 0o644)
}
