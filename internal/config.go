package internal

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	Timeout time.Duration
	LogFile string
}

var defaults = Config{
	Timeout: 15 * time.Second,
	LogFile: "/var/log/jail-mcp/server.log",
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		Timeout: defaults.Timeout,
		LogFile: defaults.LogFile,
	}

	if raw := os.Getenv("JAIL_MCP_TIMEOUT"); raw != "" {
		d, err := time.ParseDuration(raw)
		if err != nil {
			return nil, fmt.Errorf("JAIL_MCP_TIMEOUT invalid: %w", err)
		}
		cfg.Timeout = d
	}

	if logFile := os.Getenv("JAIL_MCP_LOG_FILE"); logFile != "" {
		cfg.LogFile = logFile
	}

	return cfg, nil
}
