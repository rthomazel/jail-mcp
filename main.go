package main

import (
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tcodes0/jail-mcp/handlers"
	"github.com/tcodes0/jail-mcp/internal"
)

const miseShims = "/mise/shims"

// version is set at build time via -ldflags "-X main.version=..."
var version = "local"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	current := os.Getenv("PATH")
	if !strings.Contains(current, miseShims) {
		_ = os.Setenv("PATH", miseShims+":"+current)
	}

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

	slog.Info("jail-mcp starting", "version", version, "timeout", cfg.Timeout, "background_timeout", cfg.BackgroundTimeout)

	h := handlers.New(cfg, version)

	s := server.NewMCPServer(
		"jail-mcp",
		version,
		server.WithToolCapabilities(false),
	)

	s.AddTool(
		mcp.NewTool("context",
			mcp.WithDescription("Returns environment context. Call this at the start of a session to orient yourself."),
		),
		h.HandleContext,
	)

	s.AddTool(
		mcp.NewTool("exec_sync",
			mcp.WithDescription("Execute a shell command. Returns stdout, stderr, exit code, and duration. Times out after "+cfg.Timeout.String()+". Most agents should load this now and defer exec_background."),
			mcp.WithString("command", mcp.Required(), mcp.Description("Shell command to execute")),
			mcp.WithString("cwd", mcp.Description("Working directory. Defaults to /")),
		),
		h.HandleExec,
	)

	s.AddTool(
		mcp.NewTool("exec_background",
			mcp.WithDescription("Execute a very long-running shell command in the background. Returns a job_id immediately. Use exec_status to poll for results. Times out after "+cfg.BackgroundTimeout.String()+"."),
			mcp.WithString("command", mcp.Required(), mcp.Description("Shell command to execute")),
			mcp.WithString("cwd", mcp.Description("Working directory. Defaults to /")),
		),
		h.HandleExecBackground,
	)

	s.AddTool(
		mcp.NewTool("status",
			mcp.WithDescription("Poll the status of a background job. Returns done, stdout, stderr, exit_code (if done), and duration."),
			mcp.WithString("job_id", mcp.Required(), mcp.Description("Job ID returned by exec_background")),
		),
		h.HandleStatus,
	)

	s.AddTool(
		mcp.NewTool("setup",
			mcp.WithDescription("Discover and install dependencies for the given project paths in parallel. Returns a map of project path to job_id or error. Use the status tool to poll results."),
			mcp.WithArray("paths", mcp.Required(), mcp.Description("Project paths to set up.")),
		),
		h.HandleSetup,
	)

	slog.Info("serving on stdio")
	if err := server.ServeStdio(s); err != nil {
		return fmt.Errorf("server: %w", err)
	}

	return nil
}
