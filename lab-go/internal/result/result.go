package result

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
)

// Item is the per-prompt evaluation outcome.
type Item struct {
	ID            int     `json:"id"`
	Query         string  `json:"query"`
	ShouldTrigger bool    `json:"should_trigger"`
	SkillName     string  `json:"skill_name"`
	Triggers      int     `json:"triggers"`
	Runs          int     `json:"runs"`
	TriggerRate   float64 `json:"trigger_rate"`
	Correct       bool    `json:"correct"`
}

// Summary aggregates accuracy over a prompt set.
type Summary struct {
	Total       int     `json:"total"`
	Correct     int     `json:"correct"`
	Accuracy    float64 `json:"accuracy"`
	AccuracyPct float64 `json:"accuracy_pct"`
}

// Payload is the full evaluate output document.
type Payload struct {
	SkillName          string  `json:"skill_name"`
	Provider           string  `json:"provider"`
	Model              string  `json:"model"`
	PromptsFile        string  `json:"prompts_file"`
	Workdir            string  `json:"workdir"`
	RunsDefault        int     `json:"runs_default"`
	MajorityThreshold  float64 `json:"majority_threshold"`
	Summary            Summary `json:"summary"`
	ElapsedSeconds     float64 `json:"elapsed_seconds"`
	Results            []Item  `json:"results"`
}

// BuildSummary computes accuracy from results.
func BuildSummary(results []Item) Summary {
	correct := 0
	for _, r := range results {
		if r.Correct {
			correct++
		}
	}
	total := len(results)
	acc := 0.0
	if total > 0 {
		acc = round4(float64(correct) / float64(total))
	}
	return Summary{
		Total:       total,
		Correct:     correct,
		Accuracy:    acc,
		AccuracyPct: round2(acc * 100),
	}
}

// RoundRate rounds trigger_rate to 2 decimals (matches Python).
func RoundRate(triggers, runs int) float64 {
	if runs == 0 {
		return 0
	}
	return round2(float64(triggers) / float64(runs))
}

// IsCorrect applies majority threshold prediction vs expected.
func IsCorrect(triggerRate float64, shouldTrigger bool, majorityThreshold float64) bool {
	predicted := triggerRate >= majorityThreshold
	return predicted == shouldTrigger
}

// Failures returns incorrect result items.
func Failures(payload Payload) []Item {
	var out []Item
	for _, r := range payload.Results {
		if !r.Correct {
			out = append(out, r)
		}
	}
	return out
}

// WriteJSON writes the payload to path.
func WriteJSON(path string, payload Payload) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}

// LoadJSON reads an evaluate payload from path.
func LoadJSON(path string) (Payload, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Payload{}, err
	}
	var p Payload
	if err := json.Unmarshal(data, &p); err != nil {
		return Payload{}, err
	}
	return p, nil
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

func round4(v float64) float64 {
	return math.Round(v*10000) / 10000
}
