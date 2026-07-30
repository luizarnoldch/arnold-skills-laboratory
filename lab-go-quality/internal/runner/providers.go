package runner

import (
	"bytes"
	"context"
	"io"
	"os/exec"
	"path/filepath"
	"time"
)

// OpenCode drives: opencode run --model <model> <query>
type OpenCode struct{}

func (o *OpenCode) Name() string { return "opencode" }

func (o *OpenCode) SkillInstallDir(workdir, skillName string) string {
	return filepath.Join(workdir, ".opencode", "skills", skillName)
}

func (o *OpenCode) InstallSkill(workdir, skillName, skillPath string) error {
	return InstallSkillTree(skillPath, o.SkillInstallDir(workdir, skillName))
}

func (o *OpenCode) Run(ctx context.Context, req TaskRequest) (TaskResult, error) {
	return runCLI(ctx, req, "opencode", []string{"run", "--model", req.Model, BuildTaskPrompt(req)})
}

// Claude drives: claude -p <query> --model <model>
type Claude struct{}

func (c *Claude) Name() string { return "claude" }

func (c *Claude) SkillInstallDir(workdir, skillName string) string {
	return filepath.Join(workdir, ".claude", "skills", skillName)
}

func (c *Claude) InstallSkill(workdir, skillName, skillPath string) error {
	return InstallSkillTree(skillPath, c.SkillInstallDir(workdir, skillName))
}

func (c *Claude) Run(ctx context.Context, req TaskRequest) (TaskResult, error) {
	prompt := BuildTaskPrompt(req)
	args := []string{"-p", prompt}
	if req.Model != "" {
		args = append(args, "--model", req.Model)
	}
	return runCLI(ctx, req, "claude", args)
}

// Codex drives: codex exec --model <model> <query>
type Codex struct{}

func (c *Codex) Name() string { return "codex" }

func (c *Codex) SkillInstallDir(workdir, skillName string) string {
	return filepath.Join(workdir, ".codex", "skills", skillName)
}

func (c *Codex) InstallSkill(workdir, skillName, skillPath string) error {
	return InstallSkillTree(skillPath, c.SkillInstallDir(workdir, skillName))
}

func (c *Codex) Run(ctx context.Context, req TaskRequest) (TaskResult, error) {
	prompt := BuildTaskPrompt(req)
	args := []string{"exec"}
	if req.Model != "" {
		args = append(args, "--model", req.Model)
	}
	args = append(args, prompt)
	return runCLI(ctx, req, "codex", args)
}

// CursorAgent drives: agent -p <query> --model <model>
type CursorAgent struct{}

func (c *CursorAgent) Name() string { return "cursor_agent" }

func (c *CursorAgent) SkillInstallDir(workdir, skillName string) string {
	return filepath.Join(workdir, ".cursor", "skills", skillName)
}

func (c *CursorAgent) InstallSkill(workdir, skillName, skillPath string) error {
	return InstallSkillTree(skillPath, c.SkillInstallDir(workdir, skillName))
}

func (c *CursorAgent) Run(ctx context.Context, req TaskRequest) (TaskResult, error) {
	prompt := BuildTaskPrompt(req)
	args := []string{"-p", prompt}
	if req.Model != "" {
		args = append(args, "--model", req.Model)
	}
	return runCLI(ctx, req, "agent", args)
}

func runCLI(ctx context.Context, req TaskRequest, bin string, args []string) (TaskResult, error) {
	timeout := req.Timeout
	if timeout <= 0 {
		timeout = 10 * time.Minute
	}
	runCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	start := time.Now()
	cmd := exec.CommandContext(runCtx, bin, args...)
	cmd.Dir = req.Workdir

	pr, pw := io.Pipe()
	cmd.Stdout = pw
	cmd.Stderr = pw

	if err := cmd.Start(); err != nil {
		_ = pw.Close()
		return TaskResult{}, err
	}

	var buf bytes.Buffer
	done := make(chan struct{})
	go func() {
		defer close(done)
		_, _ = io.Copy(&buf, pr)
	}()

	waitErr := cmd.Wait()
	_ = pw.Close()
	<-done

	duration := time.Since(start).Milliseconds()
	transcript := buf.String()
	timedOut := runCtx.Err() == context.DeadlineExceeded

	code := 0
	if cmd.ProcessState != nil {
		code = cmd.ProcessState.ExitCode()
	} else if waitErr != nil {
		code = -1
	}

	return TaskResult{
		Transcript:  transcript,
		DurationMS:  duration,
		TotalTokens: ParseTokensBestEffort(transcript),
		ExitCode:    code,
		TimedOut:    timedOut,
	}, nil
}
