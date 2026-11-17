package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/elyosemite/openboard/backend/internal/platform/database"
	"github.com/elyosemite/openboard/backend/internal/testsupport/postgrestest"
)

var dbURL string

func TestMain(m *testing.M) { os.Exit(runWithDatabase(m)) }

func runWithDatabase(m *testing.M) int {
	db, err := postgrestest.Start(context.Background())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	defer db.Stop()
	dbURL = db.URL
	return m.Run()
}

type syncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (s *syncBuffer) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.Write(p)
}

func (s *syncBuffer) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.String()
}

// envOf returns a getenv pointing at the test database unless overridden.
func envOf(values map[string]string) func(string) string {
	return func(key string) string {
		if v, ok := values[key]; ok {
			return v
		}
		if key == "OPENBOARD_DATABASE_URL" {
			return dbURL
		}
		return ""
	}
}

func truncateUsers(t *testing.T) {
	t.Helper()
	pool, err := database.Open(context.Background(), dbURL)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	if _, err := pool.Exec(context.Background(), `TRUNCATE users, sessions, organizations CASCADE`); err != nil {
		t.Fatal(err)
	}
}

func startAndStop(t *testing.T, env func(string) string) string {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	out := &syncBuffer{}
	done := make(chan error, 1)
	go func() { done <- run(ctx, env, out) }()

	// Start-up hashes the first Administrator's password (argon2id), which is
	// slow when every test package runs at once. The deadline only stops a hung
	// server from blocking the run; an early exit fails at once with its error.
	deadline := time.After(60 * time.Second)
	for !strings.Contains(out.String(), "http server listening") {
		select {
		case err := <-done:
			cancel()
			t.Fatalf("run exited before listening: %v: %s", err, out.String())
		case <-deadline:
			cancel()
			t.Fatalf("server did not start: %s", out.String())
		case <-time.After(10 * time.Millisecond):
		}
	}
	cancel()
	if err := <-done; err != nil {
		t.Fatalf("run returned %v", err)
	}
	return out.String()
}

func TestRunStartsCreatesAdministratorOnceAndStops(t *testing.T) {
	truncateUsers(t)
	env := envOf(map[string]string{"OPENBOARD_HTTP_ADDR": "127.0.0.1:0"})

	first := startAndStop(t, env)
	if !strings.Contains(first, "first Administrator created as admin/admin") || !strings.Contains(first, "http server stopped") {
		t.Errorf("unexpected first run log: %s", first)
	}
	second := startAndStop(t, env)
	if strings.Contains(second, "first Administrator created") {
		t.Error("the Administrator must only be created when there are no Users")
	}
}

func TestRunCreatesAdministratorFromSeed(t *testing.T) {
	truncateUsers(t)
	out := startAndStop(t, envOf(map[string]string{
		"OPENBOARD_HTTP_ADDR":      "127.0.0.1:0",
		"OPENBOARD_ADMIN_USERNAME": "root",
		"OPENBOARD_ADMIN_PASSWORD": "a long and fine passphrase",
	}))
	if !strings.Contains(out, `"username":"root"`) {
		t.Errorf("seeded administrator not logged: %s", out)
	}
}

func TestRunFailures(t *testing.T) {
	cases := map[string]struct {
		env  map[string]string
		want string
	}{
		"invalid config":  {map[string]string{"OPENBOARD_LOG_LEVEL": "loud"}, "load config"},
		"database down":   {map[string]string{"OPENBOARD_DATABASE_URL": "postgres://x:y@127.0.0.1:1/z?connect_timeout=2"}, "open database"},
		"invalid address": {map[string]string{"OPENBOARD_HTTP_ADDR": "not-an-address"}, "listen on"},
		"bad admin seed":  {map[string]string{"OPENBOARD_ADMIN_USERNAME": "root"}, "create first administrator"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			truncateUsers(t)
			err := run(context.Background(), envOf(tc.env), io.Discard)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("err = %v, want it to mention %q", err, tc.want)
			}
		})
	}
}

type downChecker struct{}

func (downChecker) Name() string                { return "down" }
func (downChecker) Check(context.Context) error { return fmt.Errorf("unreachable") }

func TestRouterServesProbesAndProblems(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cases := []struct {
		router       http.Handler
		method, path string
		status       int
	}{
		{newRouter(logger, nil), http.MethodGet, "/healthz", http.StatusOK},
		{newRouter(logger, nil), http.MethodHead, "/healthz", http.StatusOK},
		{newRouter(logger, nil), http.MethodGet, "/readyz", http.StatusOK},
		{newRouter(logger, nil, downChecker{}), http.MethodGet, "/readyz", http.StatusServiceUnavailable},
		{newRouter(logger, nil), http.MethodPost, "/healthz", http.StatusMethodNotAllowed},
		{newRouter(logger, nil), http.MethodGet, "/api/v1/unknown", http.StatusNotFound},
	}
	for _, tc := range cases {
		rec := httptest.NewRecorder()
		tc.router.ServeHTTP(rec, httptest.NewRequest(tc.method, tc.path, nil))
		if rec.Code != tc.status {
			t.Errorf("%s %s = %d, want %d", tc.method, tc.path, rec.Code, tc.status)
		}
	}
}
