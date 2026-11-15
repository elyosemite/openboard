// Package config loads the API configuration from environment variables.
package config

import (
	"fmt"
	"log/slog"
	"time"
)

const (
	defaultHTTPAddr        = ":8080"
	defaultShutdownTimeout = 10 * time.Second
)

// Config is the runtime configuration of the API process.
type Config struct {
	HTTPAddr        string
	LogLevel        slog.Level
	ShutdownTimeout time.Duration
}

// Load reads the configuration using getenv (usually os.Getenv), applying
// defaults for unset variables and rejecting invalid values.
func Load(getenv func(string) string) (Config, error) {
	cfg := Config{
		HTTPAddr:        defaultHTTPAddr,
		LogLevel:        slog.LevelInfo,
		ShutdownTimeout: defaultShutdownTimeout,
	}

	if raw := getenv("OPENBOARD_HTTP_ADDR"); raw != "" {
		cfg.HTTPAddr = raw
	}

	if raw := getenv("OPENBOARD_LOG_LEVEL"); raw != "" {
		if err := cfg.LogLevel.UnmarshalText([]byte(raw)); err != nil {
			return Config{}, fmt.Errorf("OPENBOARD_LOG_LEVEL: %w", err)
		}
	}

	if raw := getenv("OPENBOARD_SHUTDOWN_TIMEOUT"); raw != "" {
		timeout, err := parsePositiveDuration(raw)
		if err != nil {
			return Config{}, fmt.Errorf("OPENBOARD_SHUTDOWN_TIMEOUT: %w", err)
		}
		cfg.ShutdownTimeout = timeout
	}

	return cfg, nil
}

func parsePositiveDuration(raw string) (time.Duration, error) {
	d, err := time.ParseDuration(raw)
	if err != nil {
		return 0, err
	}
	if d <= 0 {
		return 0, fmt.Errorf("duration must be positive, got %q", raw)
	}
	return d, nil
}
