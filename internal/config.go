package internal

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	Timeout           time.Duration
	BackgroundTimeout time.Duration
	Home              string
	MiseDir           string
}

var defaults = Config{
	Timeout:           15 * time.Second,
	BackgroundTimeout: 5 * time.Minute,
}

func LoadConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("home directory: %w", err)
	}

	if raw := os.Getenv("JAIL_MCP_HOME"); raw != "" {
		home = raw
	}

	miseDir := "/mise"
	if raw := os.Getenv("JAIL_MCP_MISE_DIR"); raw != "" {
		miseDir = raw
	}

	cfg := &Config{
		Timeout:           defaults.Timeout,
		BackgroundTimeout: defaults.BackgroundTimeout,
		Home:              home,
		MiseDir:           miseDir,
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
