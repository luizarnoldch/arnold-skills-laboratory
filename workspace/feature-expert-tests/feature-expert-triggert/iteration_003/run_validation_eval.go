package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	defaultModel   = "digitalocean/deepseek-4-flash"
	defaultRuns    = 1
	defaultTimeout = 60
	skillName      = "feature-expert"
	inputFile      = "train_queries.json"
	outputFile     = "results/train_results.json"
	logBaseDir     = "log"
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

// getNextRunDir scans logBaseDir for existing run_# folders and returns the next run path
func getNextRunDir(baseDir string) (string, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return "", err
	}

	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`^run_(\d+)$`)
	maxRun := 0

	for _, entry := range entries {
		if entry.IsDir() {
			matches := re.FindStringSubmatch(entry.Name())
			if len(matches) == 2 {
				num, err := strconv.Atoi(matches[1])
				if err == nil && num > maxRun {
					maxRun = num
				}
			}
		}
	}

	nextRunNum := maxRun + 1
	nextRunDir := filepath.Join(baseDir, fmt.Sprintf("run_%d", nextRunNum))
	if err := os.MkdirAll(nextRunDir, 0755); err != nil {
		return "", err
	}

	return nextRunDir, nil
}

func main() {
	// Define command-line flags
	runsFlag := flag.Int("r", defaultRuns, "Default number of runs per query")
	modelFlag := flag.String("m", defaultModel, "Model to evaluate")
	timeoutFlag := flag.Int("s", defaultTimeout, "Timeout per run in seconds")

	flag.Parse()

	model := *modelFlag
	runsDefault := *runsFlag
	runTimeout := time.Duration(*timeoutFlag) * time.Second

	if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
		log.Fatalf("Error creating results directory: %v", err)
	}

	currentLogDir, err := getNextRunDir(logBaseDir)
	if err != nil {
		log.Fatalf("Error setting up run log directory: %v", err)
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
	fmt.Printf("Default Runs: %d | Timeout: %s\n", runsDefault, runTimeout)
	fmt.Printf("Log Directory: ./%s/\n", currentLogDir)
	fmt.Println("==================================================")

	var results []ResultItem

	for _, q := range queries {
		runs := runsDefault
		if q.Runs != nil {
			runs = *q.Runs
		}

		fmt.Printf("\n[ID: %d] Query: %q | Expected: %t | Runs: %d\n", q.ID, q.Query, q.ShouldTrigger, runs)

		triggerCount := 0
		targetSubstring := fmt.Sprintf(`Skill "%s"`, skillName)
		fallbackSubstring := fmt.Sprintf(`"name":"%s"`, skillName)

		for i := 1; i <= runs; i++ {
			var logFilename string
			if runs == 1 {
				logFilename = filepath.Join(currentLogDir, fmt.Sprintf("id_%d.log", q.ID))
			} else {
				logFilename = filepath.Join(currentLogDir, fmt.Sprintf("id_%d_run_%d.log", q.ID, i))
			}

			logFile, err := os.Create(logFilename)
			if err != nil {
				log.Printf("Warning: failed to create log file %s: %v", logFilename, err)
			}

			ctx, cancel := context.WithTimeout(context.Background(), runTimeout)

			pr, pw := io.Pipe()
			var logWriters []io.Writer

			if logFile != nil {
				header := fmt.Sprintf("=== RUN METADATA ===\nID: %d\nTimestamp: %s\nQuery: %s\nModel: %s\nRun: %d/%d\nExpected Trigger: %t\nTimeout: %s\n====================\n\n",
					q.ID, time.Now().Format(time.RFC3339), q.Query, model, i, runs, q.ShouldTrigger, runTimeout)
				_, _ = logFile.WriteString(header)
				logWriters = append(logWriters, logFile)
			}

			multiWriter := io.MultiWriter(append(logWriters, pw)...)

			cmd := exec.CommandContext(ctx, "opencode", "run", "--model", model, q.Query)
			cmd.Stdout = multiWriter
			cmd.Stderr = multiWriter

			if err := cmd.Start(); err != nil {
				log.Printf("Error starting process: %v", err)
				pw.Close()
				cancel()
				if logFile != nil {
					logFile.Close()
				}
				continue
			}

			triggered := false
			doneScanner := make(chan struct{})

			go func() {
				defer close(doneScanner)
				var buf bytes.Buffer
				chunk := make([]byte, 1024)

				for {
					n, err := pr.Read(chunk)
					if n > 0 {
						buf.Write(chunk[:n])
						currentOut := buf.String()

						if !triggered && (strings.Contains(currentOut, targetSubstring) || strings.Contains(currentOut, fallbackSubstring)) {
							triggered = true
							if cmd.Process != nil {
								_ = cmd.Process.Kill()
							}
						}
					}
					if err != nil {
						break
					}
				}
			}()

			err = cmd.Wait()
			pw.Close()
			<-doneScanner

			isTimeout := ctx.Err() == context.DeadlineExceeded
			cancel()

			if logFile != nil {
				logFile.Close()
			}

			if triggered {
				triggerCount++
				fmt.Printf("  └─ Run %d/%d: TRIGGERED (%s) [Terminated Early]\n", i, runs, skillName)
			} else if isTimeout {
				fmt.Printf("  └─ Run %d/%d: TIMED OUT (> %s)\n", i, runs, runTimeout)
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

	fmt.Printf("\nEvaluation complete -> Results saved to %s\nLogs saved in: ./%s/\n", outputFile, currentLogDir)
}
