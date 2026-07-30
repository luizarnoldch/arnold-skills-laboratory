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

var (
	// totalTokensRe matches explicit total / usage totals.
	totalTokensRe = regexp.MustCompile(`(?i)(?:total[_\s-]?tokens|token[_\s-]?count|tokens?\s*(?:used|usage)|usage["\s:=]+\{[^}]{0,200}?total)\D{0,8}(\d{1,9})`)
	// pairTokensRe finds input/prompt and output/completion token fields for summing.
	inputTokensRe  = regexp.MustCompile(`(?i)(?:input[_\s-]?tokens|prompt[_\s-]?tokens)\D{0,8}(\d{1,9})`)
	outputTokensRe = regexp.MustCompile(`(?i)(?:output[_\s-]?tokens|completion[_\s-]?tokens)\D{0,8}(\d{1,9})`)
)

// ParseTokensBestEffort extracts a token count from transcript text when present.
// Prefers an explicit total; otherwise sums the last input+output (or prompt+completion) pair.
// Returns 0 when the CLI transcript does not expose usage (common for some providers).
func ParseTokensBestEffort(transcript string) int64 {
	if n := lastInt(totalTokensRe, transcript); n > 0 {
		return n
	}
	in := lastInt(inputTokensRe, transcript)
	out := lastInt(outputTokensRe, transcript)
	if in > 0 || out > 0 {
		return in + out
	}
	return 0
}

func lastInt(re *regexp.Regexp, s string) int64 {
	matches := re.FindAllStringSubmatch(s, -1)
	if len(matches) == 0 {
		return 0
	}
	n, err := strconv.ParseInt(matches[len(matches)-1][1], 10, 64)
	if err != nil {
		return 0
	}
	return n
}

// EvalInputDestRel returns the destination path relative to the sandbox for an
// eval input file. Prefer paths under evals/files/<evalName>/; else under
// evals/files/; else basename (legacy flat copy).
func EvalInputDestRel(src, evalName string) string {
	src = filepath.ToSlash(src)
	marker := "/evals/files/"
	idx := strings.LastIndex(src, marker)
	if idx < 0 {
		return filepath.Base(src)
	}
	after := src[idx+len(marker):]
	name := strings.TrimSpace(evalName)
	if name != "" {
		prefix := name + "/"
		if strings.HasPrefix(after, prefix) {
			return filepath.FromSlash(strings.TrimPrefix(after, prefix))
		}
		// exact file named like the eval folder is uncommon; keep after as-is below
	}
	if after == "" || after == "." {
		return filepath.Base(src)
	}
	return filepath.FromSlash(after)
}

// CopyInputs copies input files into workdir using basename only (legacy).
func CopyInputs(workdir string, files []string) ([]string, error) {
	return CopyEvalInputs(workdir, files, "")
}

// CopyEvalInputs copies eval fixture files into workdir, preserving structure
// under evals/files/<evalName>/ when evalName is set.
func CopyEvalInputs(workdir string, files []string, evalName string) ([]string, error) {
	var dests []string
	for _, src := range files {
		if strings.TrimSpace(src) == "" {
			continue
		}
		rel := EvalInputDestRel(src, evalName)
		if rel == "" || rel == "." {
			rel = filepath.Base(src)
		}
		dst := filepath.Join(workdir, rel)
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return nil, err
		}
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
