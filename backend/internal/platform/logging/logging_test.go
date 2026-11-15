package logging

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"
)

func TestNewWritesJSONAndRespectsLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := New(&buf, slog.LevelInfo)

	logger.Debug("hidden")
	logger.Info("visible", "key", "value")

	var entry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("expected a single JSON entry, got %q: %v", buf.String(), err)
	}
	if entry["msg"] != "visible" || entry["key"] != "value" {
		t.Errorf("unexpected entry: %v", entry)
	}
}
