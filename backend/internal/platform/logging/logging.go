// Package logging builds the structured logger used across the API.
package logging

import (
	"io"
	"log/slog"
)

// New returns a JSON slog logger writing to w at the given minimum level.
func New(w io.Writer, level slog.Level) *slog.Logger {
	return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: level}))
}
