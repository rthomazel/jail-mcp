package main

import (
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tcodes0/jail-mcp/internal"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := internal.LoadConfig()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(
		os.Stderr,
		&slog.HandlerOptions{Level: slog.LevelInfo},
	)))

	defer func() {
		if msg := recover(); msg != nil {
			slog.Error("panic", "msg", msg, "stack", string(debug.Stack()))
			os.Exit(1)
		}
	}()

	slog.Info("jail-mcp starting", "timeout", cfg.Timeout)

	handler := internal.NewHandler(cfg)

	s := server.NewMCPServer(
		"jail-mcp",
		"local",
		server.WithToolCapabilities(false),
	)

	s.AddTool(
		mcp.NewTool("shell_exec",
			mcp.WithDescription("Execute any shell command inside the container. Returns stdout, stderr, exit code, and duration."),
			mcp.WithString("command", mcp.Required(), mcp.Description("Shell command to execute")),
			mcp.WithString("cwd", mcp.Description("Working directory inside the container. Defaults to /")),
		),
		handler.HandleExec,
	)

	s.AddTool(
		mcp.NewTool("context",
			mcp.WithDescription("Returns environment context: mounted projects, OS, available tools, disk space, and log file path. Call this at the start of a session to orient yourself."),
		),
		handler.HandleContext,
	)

	slog.Info("serving on stdio")
	if err := server.ServeStdio(s); err != nil {
		return fmt.Errorf("server: %w", err)
	}

	return nil
}
