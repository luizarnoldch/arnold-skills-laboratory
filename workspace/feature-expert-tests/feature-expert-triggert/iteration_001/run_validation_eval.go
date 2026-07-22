package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	model       = "digitalocean/deepseek-4-flash"
	runsDefault = 1
	skillName   = "feature-expert"
	inputFile   = "validation_queries.json"
	outputFile  = "results/validation_results.json"
	logDir      = "log"
)

type QueryItem struct {
	ID            int    `json:"id"`
	Query         string `json:"query"`
	ShouldTrigger bool   `json:"should_trigger"`
	Runs          *int   `json:"runs,omitempty"`
}

type ResultItem struct {
	ID            int     `json:"id"`
	Query         string  `json:"query"`
	ShouldTrigger bool    `json:"should_trigger"`
	Triggers      bool    `json:"triggers"`
	Runs          int     `json:"runs"`
	TriggerRate   float64 `json:"trigger_rate"`
}

// sanitizeFilename generates a safe filename incorporating the ID and a short hash for collision prevention.
func sanitizeFilename(id int, s string) string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9_\-]+`)
	safe := reg.ReplaceAllString(s, "_")
	safe = strings.Trim(safe, "_")
	if len(safe) > 35 {
		safe = safe[:35]
	}
	if safe == "" {
		safe = "query"
	}

	hash := sha256.Sum256([]byte(s))
	shortHash := hex.EncodeToString(hash[:])[:6]

	return fmt.Sprintf("id_%d_%s_%s", id, strings.ToLower(safe), shortHash)
}

func main() {
	if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
		log.Fatalf("Error creating results directory: %v", err)
	}
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("Error creating log directory: %v", err)
	}

	inputBytes, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Error: %s not found in current directory.", inputFile)
	}

	var queries []QueryItem
	if err := json.Unmarshal(inputBytes, &queries); err != nil {
		log.Fatalf("Error parsing JSON in %s: %v", inputFile, err)
	}

	fmt.Println("==================================================")
	fmt.Printf("Evaluating Validation Queries with model: %s\n", model)
	fmt.Println("==================================================")

	var results []ResultItem

	for _, q := range queries {
		runs := runsDefault
		if q.Runs != nil {
			runs = *q.Runs
		}

		fmt.Printf("\n[ID: %d] Query: %q | Expected: %t | Runs: %d\n", q.ID, q.Query, q.ShouldTrigger, runs)

		triggerCount := 0
		targetSubstring := fmt.Sprintf(`"name":"%s"`, skillName)
		baseLogName := sanitizeFilename(q.ID, q.Query)

		for i := 1; i <= runs; i++ {
			logFilename := filepath.Join(logDir, fmt.Sprintf("%s_run%d.log", baseLogName, i))
			logFile, err := os.Create(logFilename)
			if err != nil {
				log.Printf("Warning: failed to create log file %s: %v", logFilename, err)
			}

			var memBuf bytes.Buffer
			var writers []io.Writer
			writers = append(writers, &memBuf)

			if logFile != nil {
				header := fmt.Sprintf("=== RUN METADATA ===\nID: %d\nTimestamp: %s\nQuery: %s\nModel: %s\nRun: %d/%d\nExpected Trigger: %t\n====================\n\n",
					q.ID, time.Now().Format(time.RFC3339), q.Query, model, i, runs, q.ShouldTrigger)
				_, _ = logFile.WriteString(header)

				writers = append(writers, logFile)
				defer logFile.Close()
			}

			multiWriter := io.MultiWriter(writers...)

			cmd := exec.Command("opencode", "run", "--model", model, q.Query)
			cmd.Stdout = multiWriter
			cmd.Stderr = multiWriter

			_ = cmd.Run()

			outputStr := memBuf.String()

			if strings.Contains(outputStr, targetSubstring) {
				triggerCount++
				fmt.Printf("  └─ Run %d/%d: TRIGGERED (%s)\n", i, runs, skillName)
			} else {
				fmt.Printf("  └─ Run %d/%d: NOT TRIGGERED\n", i, runs)
			}
		}

		triggerRate := float64(triggerCount) / float64(runs)
		triggerRateRounded := float64(int(triggerRate*100+0.5)) / 100.0

		results = append(results, ResultItem{
			ID:            q.ID,
			Query:         q.Query,
			ShouldTrigger: q.ShouldTrigger,
			Triggers:      triggerCount > 0,
			Runs:          runs,
			TriggerRate:   triggerRateRounded,
		})
	}

	outputBytes, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Fatalf("Error encoding output JSON: %v", err)
	}

	if err := os.WriteFile(outputFile, outputBytes, 0644); err != nil {
		log.Fatalf("Error writing results to %s: %v", outputFile, err)
	}

	fmt.Printf("\nEvaluation complete -> Results saved to %s (Logs saved in ./%s/)\n", outputFile, logDir)
}