package provider

import (
	"bytes"
	"context"
	"io"
	"os/exec"
	"sync"
)

// OpenCode drives the opencode CLI.
type OpenCode struct {
	base
}

func (o *OpenCode) Name() string { return "opencode" }

func (o *OpenCode) Run(ctx context.Context, query string) (RunResult, error) {
	runCtx, cancel := context.WithTimeout(ctx, o.cfg.Timeout)
	defer cancel()

	cmd := exec.CommandContext(runCtx, "opencode", "run", "--model", o.cfg.Model, query)
	cmd.Dir = o.cfg.Workdir

	pr, pw := io.Pipe()
	cmd.Stdout = pw
	cmd.Stderr = pw

	if err := cmd.Start(); err != nil {
		_ = pw.Close()
		return RunResult{}, err
	}

	var (
		mu        sync.Mutex
		buf       bytes.Buffer
		triggered bool
	)

	done := make(chan struct{})
	go func() {
		defer close(done)
		chunk := make([]byte, 1024)
		for {
			n, err := pr.Read(chunk)
			if n > 0 {
				mu.Lock()
				buf.Write(chunk[:n])
				out := buf.String()
				if !triggered && DetectTrigger(out, o.cfg.SkillName) {
					triggered = true
					if cmd.Process != nil {
						_ = cmd.Process.Kill()
					}
				}
				mu.Unlock()
			}
			if err != nil {
				return
			}
		}
	}()

	waitErr := cmd.Wait()
	_ = pw.Close()
	<-done

	mu.Lock()
	stdoutStr := buf.String()
	wasTriggered := triggered
	mu.Unlock()

	timedOut := runCtx.Err() == context.DeadlineExceeded && !wasTriggered
	if !wasTriggered {
		wasTriggered = DetectTrigger(stdoutStr, o.cfg.SkillName)
	}

	code := 0
	if cmd.ProcessState != nil {
		code = cmd.ProcessState.ExitCode()
	} else if waitErr != nil {
		code = -1
	}

	return RunResult{
		Stdout:    stdoutStr,
		Triggered: wasTriggered,
		TimedOut:  timedOut,
		Code:      code,
	}, nil
}
