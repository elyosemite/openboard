package main

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

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

func envOf(values map[string]string) func(string) string {
	return func(key string) string { return values[key] }
}

func TestRunStartsAndStopsCleanly(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	out := &syncBuffer{}
	done := make(chan error, 1)
	go func() {
		done <- run(ctx, envOf(map[string]string{"OPENBOARD_HTTP_ADDR": "127.0.0.1:0"}), out)
	}()

	deadline := time.Now().Add(5 * time.Second)
	for !strings.Contains(out.String(), "http server listening") {
		if time.Now().After(deadline) {
			t.Fatal("server did not start")
		}
		time.Sleep(10 * time.Millisecond)
	}

	cancel()
	if err := <-done; err != nil {
		t.Fatalf("run returned %v", err)
	}
	if !strings.Contains(out.String(), "http server stopped") {
		t.Error("missing shutdown log")
	}
}

func TestRunFailsOnInvalidConfig(t *testing.T) {
	err := run(context.Background(), envOf(map[string]string{"OPENBOARD_LOG_LEVEL": "loud"}), io.Discard)
	if err == nil || !strings.Contains(err.Error(), "load config") {
		t.Fatalf("err = %v", err)
	}
}

func TestRunFailsOnInvalidAddress(t *testing.T) {
	err := run(context.Background(), envOf(map[string]string{"OPENBOARD_HTTP_ADDR": "not-an-address"}), io.Discard)
	if err == nil || !strings.Contains(err.Error(), "listen on") {
		t.Fatalf("err = %v", err)
	}
}

func TestRouterServesProbesAndProblems(t *testing.T) {
	router := newRouter(slog.New(slog.NewTextHandler(io.Discard, nil)))
	cases := []struct {
		method, path string
		status       int
	}{
		{http.MethodGet, "/healthz", http.StatusOK},
		{http.MethodHead, "/healthz", http.StatusOK},
		{http.MethodGet, "/readyz", http.StatusOK},
		{http.MethodPost, "/healthz", http.StatusMethodNotAllowed},
		{http.MethodGet, "/api/v1/unknown", http.StatusNotFound},
	}
	for _, tc := range cases {
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, httptest.NewRequest(tc.method, tc.path, nil))
		if rec.Code != tc.status {
			t.Errorf("%s %s = %d, want %d", tc.method, tc.path, rec.Code, tc.status)
		}
	}
}
