package runner

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// TaskRequest describes one quality-eval agent run.
type TaskRequest struct {
	Prompt     string
	SkillName  string
	SkillPath  string // empty = without skill
	InputFiles []string
	OutputDir  string
	Workdir    string
	Model      string
	Timeout    time.Duration
}

// TaskResult is the outcome of a full task run.
type TaskResult struct {
	Transcript  string
	DurationMS  int64
	TotalTokens int64
	ExitCode    int
	TimedOut    bool
}

// TaskRunner executes a full agent task (no early kill on skill trigger).
type TaskRunner interface {
	Name() string
	Run(ctx context.Context, req TaskRequest) (TaskResult, error)
	InstallSkill(workdir, skillName, skillPath string) error
	SkillInstallDir(workdir, skillName string) string
}

// BuildTaskPrompt builds the agentskills-style task instructions.
func BuildTaskPrompt(req TaskRequest) string {
	var b strings.Builder
	b.WriteString("Execute this task:\n")
	if req.SkillPath != "" {
		b.WriteString(fmt.Sprintf("- Skill path: %s\n", req.SkillPath))
	} else {
		b.WriteString("- Skill path: (none — do not use a project skill)\n")
	}
	b.WriteString(fmt.Sprintf("- Task: %s\n", strings.TrimSpace(req.Prompt)))
	if len(req.InputFiles) > 0 {
		b.WriteString("- Input files:\n")
		for _, f := range req.InputFiles {
			b.WriteString(fmt.Sprintf("  - %s\n", f))
		}
	}
	b.WriteString(fmt.Sprintf("- Save outputs to: %s\n", req.OutputDir))
	b.WriteString("\nComplete the task. Write any produced files into the Save outputs directory.\n")
	return b.String()
}

var tokenRe = regexp.MustCompile(`(?i)(?:total[_\s-]?tokens|input[_\s-]?tokens|output[_\s-]?tokens|prompt[_\s-]?tokens|completion[_\s-]?tokens)["\s:=]+(\d+)`)

// ParseTokensBestEffort extracts a token count from transcript text when present.
func ParseTokensBestEffort(transcript string) int64 {
	matches := tokenRe.FindAllStringSubmatch(transcript, -1)
	if len(matches) == 0 {
		return 0
	}
	var sum int64
	seen := false
	for _, m := range matches {
		n, err := strconv.ParseInt(m[1], 10, 64)
		if err != nil {
			continue
		}
		sum += n
		seen = true
	}
	if !seen {
		return 0
	}
	// Prefer last "total_tokens" style if present
	last := matches[len(matches)-1]
	if n, err := strconv.ParseInt(last[1], 10, 64); err == nil && strings.Contains(strings.ToLower(last[0]), "total") {
		return n
	}
	return sum
}

// CopyInputs copies input files into workdir, preserving basename.
func CopyInputs(workdir string, files []string) ([]string, error) {
	var dests []string
	for _, src := range files {
		if strings.TrimSpace(src) == "" {
			continue
		}
		base := filepath.Base(src)
		dst := filepath.Join(workdir, base)
		data, err := os.ReadFile(src)
		if err != nil {
			return nil, fmt.Errorf("copy input %s: %w", src, err)
		}
		if err := os.WriteFile(dst, data, 0o644); err != nil {
			return nil, err
		}
		dests = append(dests, dst)
	}
	return dests, nil
}

// InstallSkillTree copies skill directory into installDir.
func InstallSkillTree(skillPath, installDir string) error {
	if err := os.MkdirAll(filepath.Dir(installDir), 0o755); err != nil {
		return err
	}
	_ = os.RemoveAll(installDir)
	return copyDir(skillPath, installDir)
}

func copyDir(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		// skill path may be SKILL.md file — install parent
		return fmt.Errorf("skill path must be a directory containing SKILL.md: %s", src)
	}
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, e := range entries {
		s := filepath.Join(src, e.Name())
		d := filepath.Join(dst, e.Name())
		if e.IsDir() {
			if err := copyDir(s, d); err != nil {
				return err
			}
			continue
		}
		data, err := os.ReadFile(s)
		if err != nil {
			return err
		}
		info, _ := e.Info()
		mode := os.FileMode(0o644)
		if info != nil {
			mode = info.Mode()
		}
		if err := os.WriteFile(d, data, mode); err != nil {
			return err
		}
	}
	return nil
}

// ResolveSkillDir normalizes -skill-path: if pointing at SKILL.md, use parent.
func ResolveSkillDir(skillPath string) (string, error) {
	if skillPath == "" {
		return "", nil
	}
	info, err := os.Stat(skillPath)
	if err != nil {
		return "", err
	}
	if info.IsDir() {
		return skillPath, nil
	}
	if strings.EqualFold(filepath.Base(skillPath), "SKILL.md") {
		return filepath.Dir(skillPath), nil
	}
	return "", fmt.Errorf("skill path must be a skill directory or SKILL.md: %s", skillPath)
}
