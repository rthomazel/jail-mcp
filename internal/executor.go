package internal

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type Result struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Duration time.Duration
	Error    string
}

type Executor struct {
	cfg *Config
	log *Logger
}

func NewExecutor(cfg *Config, log *Logger) *Executor {
	return &Executor{cfg: cfg, log: log}
}

func (e *Executor) Run(ctx context.Context, command, cwd string) *Result {
	start := time.Now()
	e.log.Info("exec start", "cmd", command, "cwd", cwd)

	ctx, cancel := context.WithTimeout(ctx, e.cfg.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "bash", "-c", command)
	cmd.Dir = cwd

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	duration := time.Since(start)
	exitCode := 0

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			e.log.Error("exec failed to start", "cmd", command, "err", err)
			return &Result{
				Duration: duration,
				ExitCode: -1,
				Error:    fmt.Sprintf("could not start process: %v", err),
			}
		}
	}

	e.log.Info("exec done",
		"cmd", command,
		"exit_code", exitCode,
		"duration", duration.Round(time.Millisecond),
	)

	return &Result{
		Stdout:   strings.TrimRight(stdout.String(), "\n"),
		Stderr:   strings.TrimRight(stderr.String(), "\n"),
		ExitCode: exitCode,
		Duration: duration,
	}
}
