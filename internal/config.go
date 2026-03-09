package internal

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	Timeout           time.Duration
	BackgroundTimeout time.Duration
}

var defaults = Config{
	Timeout:           15 * time.Second,
	BackgroundTimeout: 5 * time.Minute,
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		Timeout:           defaults.Timeout,
		BackgroundTimeout: defaults.BackgroundTimeout,
	}

	if raw := os.Getenv("JAIL_MCP_TIMEOUT"); raw != "" {
		d, err := time.ParseDuration(raw)
		if err != nil {
			return nil, fmt.Errorf("JAIL_MCP_TIMEOUT invalid: %w", err)
		}
		cfg.Timeout = d
	}

	if raw := os.Getenv("JAIL_MCP_BACKGROUND_TIMEOUT"); raw != "" {
		d, err := time.ParseDuration(raw)
		if err != nil {
			return nil, fmt.Errorf("JAIL_MCP_BACKGROUND_TIMEOUT invalid: %w", err)
		}
		cfg.BackgroundTimeout = d
	}

	return cfg, nil
}
