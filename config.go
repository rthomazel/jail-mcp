package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	AllowedDirs []string
	Timeout     time.Duration
	LogFile     string
}

// loadConfig reads configuration from environment variables. See README for details.
func loadConfig() (*Config, error) {
	dirsRaw := os.Getenv("JAIL_MCP_DIRS")
	if dirsRaw == "" {
		return nil, fmt.Errorf("JAIL_MCP_DIRS is required (colon-separated list of allowed dirs)")
	}

	var allowed []string
	for _, d := range strings.Split(dirsRaw, ":") {
		if d = strings.TrimSpace(d); d != "" {
			allowed = append(allowed, d)
		}
	}
	if len(allowed) == 0 {
		return nil, fmt.Errorf("JAIL_MCP_DIRS contained no valid directories")
	}

	timeout := 30 * time.Second
	if raw := os.Getenv("JAIL_MCP_TIMEOUT"); raw != "" {
		d, err := time.ParseDuration(raw)
		if err != nil {
			return nil, fmt.Errorf("JAIL_MCP_TIMEOUT invalid: %w", err)
		}
		timeout = d
	}

	logFile := os.Getenv("JAIL_MCP_LOG")
	if logFile == "" {
		logFile = "/var/log/jail-mcp/jail.log"
	}

	return &Config{
		AllowedDirs: allowed,
		Timeout:     timeout,
		LogFile:     logFile,
	}, nil
}
