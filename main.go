package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tcodes0/jail-mcp/internal"
)

var version = "dev"

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

	log, err := internal.NewLogger(cfg.LogFile)
	if err != nil {
		return fmt.Errorf("logger: %w", err)
	}
	defer log.Close()

	defer func() {
		if msg := recover(); msg != nil {
			log.Error("panic", "msg", msg, "stack", string(debug.Stack()))
			log.Close()
			os.Exit(1)
		}
	}()

	log.Info("jail-mcp starting", "version", version, "dirs", cfg.AllowedDirs, "timeout", cfg.Timeout)

	executor := internal.NewExecutor(cfg, log)
	handler := internal.NewHandler(executor, cfg, log)

	s := server.NewMCPServer(
		"jail-mcp",
		version,
		server.WithToolCapabilities(false),
	)

	s.AddTool(
		mcp.NewTool("shell_exec",
			mcp.WithDescription("Execute any shell command inside the container. Returns stdout, stderr, exit code, and duration."),
			mcp.WithString("command", mcp.Required(), mcp.Description("Shell command to execute")),
			mcp.WithString("cwd", mcp.Description("Working directory. Must be one of the allowed dirs or a subpath. Defaults to first allowed dir.")),
		),
		handler.HandleExec,
	)

	s.AddTool(
		mcp.NewTool("list_dirs",
			mcp.WithDescription("List the directories available inside this container."),
		),
		handler.HandleListDirs,
	)

	log.Info("serving on stdio")
	if err := server.ServeStdio(s); err != nil {
		return fmt.Errorf("server: %w", err)
	}

	return nil
}
