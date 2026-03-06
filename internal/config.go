package internal

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

var defaults = Config{
	Timeout: 30 * time.Second,
	LogFile: "/var/log/jail-mcp/jail.log",
}

func LoadConfig() (*Config, error) {
	dirsRaw := os.Getenv("JAIL_MCP_DIRS")
	if dirsRaw == "" {
		return nil, fmt.Errorf("JAIL_MCP_DIRS is required (colon-separated list of allowed dirs)")
	}

	var allowed []string
	for d := range strings.SplitSeq(dirsRaw, ":") {
		if d = strings.TrimSpace(d); d != "" {
			allowed = append(allowed, d)
		}
	}

	if len(allowed) == 0 {
		return nil, fmt.Errorf("JAIL_MCP_DIRS contained no valid directories")
	}

	cfg := &Config{
		AllowedDirs: allowed,
		Timeout:     defaults.Timeout,
		LogFile:     defaults.LogFile,
	}

	if raw := os.Getenv("JAIL_MCP_TIMEOUT"); raw != "" {
		d, err := time.ParseDuration(raw)
		if err != nil {
			return nil, fmt.Errorf("JAIL_MCP_TIMEOUT invalid: %w", err)
		}
		cfg.Timeout = d
	}

	if logFile := os.Getenv("JAIL_MCP_LOG"); logFile != "" {
		cfg.LogFile = logFile
	}

	return cfg, nil
}
