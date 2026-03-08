package internal

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	Timeout time.Duration
}

var defaults = Config{
	Timeout: 15 * time.Second,
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		Timeout: defaults.Timeout,
	}

	if raw := os.Getenv("JAIL_MCP_TIMEOUT"); raw != "" {
		d, err := time.ParseDuration(raw)
		if err != nil {
			return nil, fmt.Errorf("JAIL_MCP_TIMEOUT invalid: %w", err)
		}
		cfg.Timeout = d
	}

	return cfg, nil
}
