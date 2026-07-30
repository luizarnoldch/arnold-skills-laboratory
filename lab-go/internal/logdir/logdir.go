package logdir

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

var runRe = regexp.MustCompile(`^run_(\d+)$`)

// NextRunDir creates and returns the next run_N directory under base.
func NextRunDir(base string) (string, error) {
	if err := os.MkdirAll(base, 0o755); err != nil {
		return "", err
	}
	entries, err := os.ReadDir(base)
	if err != nil {
		return "", err
	}
	maxRun := 0
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		m := runRe.FindStringSubmatch(e.Name())
		if len(m) != 2 {
			continue
		}
		n, err := strconv.Atoi(m[1])
		if err == nil && n > maxRun {
			maxRun = n
		}
	}
	next := filepath.Join(base, fmt.Sprintf("run_%d", maxRun+1))
	if err := os.MkdirAll(next, 0o755); err != nil {
		return "", err
	}
	return next, nil
}

var iterRe = regexp.MustCompile(`^(\d+)$`)

// NextIterationDir creates and returns the next NNN directory under base.
func NextIterationDir(base string) (string, error) {
	if err := os.MkdirAll(base, 0o755); err != nil {
		return "", err
	}
	entries, err := os.ReadDir(base)
	if err != nil {
		return "", err
	}
	maxN := 0
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		m := iterRe.FindStringSubmatch(e.Name())
		if len(m) != 2 {
			continue
		}
		n, err := strconv.Atoi(m[1])
		if err == nil && n > maxN {
			maxN = n
		}
	}
	next := filepath.Join(base, fmt.Sprintf("%03d", maxN+1))
	if err := os.MkdirAll(next, 0o755); err != nil {
		return "", err
	}
	return next, nil
}
