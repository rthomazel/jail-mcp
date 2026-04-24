// jail mcp server binary
package main

import (
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rthomazel/jail-mcp/handlers"
	"github.com/rthomazel/jail-mcp/internal"
	"github.com/rthomazel/jail-mcp/internal/pathsnapshot"
)

// version is set at build time via -ldflags "-X main.version=..."
var version = "local"

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

	miseShims := cfg.MiseDir + "/shims"
	homeBin := cfg.Home + "/bin"
	_ = os.MkdirAll(homeBin, 0o755)

	current := os.Getenv("PATH")
	if !strings.Contains(current, miseShims) {
		current = miseShims + ":" + current
	}
	if !strings.Contains(current, homeBin) {
		current = homeBin + ":" + current
	}
	_ = os.Setenv("PATH", current)

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

	pathsnapshot.Diff(cfg.Home)

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
			mcp.WithDescription("Execute one or more shell commands. Returns stdout, stderr, exit code, and duration per command. Times out after "+cfg.Timeout.String()+". Most agents should load this now and defer exec_background."),
			mcp.WithArray("commands", mcp.Required(), mcp.Description("Shell commands to execute.")),
			mcp.WithString("cwd", mcp.Description("Working directory. Defaults to /")),
		),
		h.HandleExec,
	)

	s.AddTool(
		mcp.NewTool("exec_background",
			mcp.WithDescription("Execute one or more shell commands in the background. Returns a job_id per command immediately. Use exec_status to poll for results. Times out after "+cfg.BackgroundTimeout.String()+"."),
			mcp.WithArray("commands", mcp.Required(), mcp.Description("Shell commands to execute.")),
			mcp.WithString("cwd", mcp.Description("Working directory. Defaults to /")),
		),
		h.HandleExecBackground,
	)

	s.AddTool(
		mcp.NewTool("status",
			mcp.WithDescription("Poll the status of one or more background jobs. Returns done, stdout, stderr, exit_code (if done), and duration per job."),
			mcp.WithArray("job_ids", mcp.Required(), mcp.Description("Job IDs returned by exec_background.")),
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
