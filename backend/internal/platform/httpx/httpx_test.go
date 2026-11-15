package httpx

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func decodeProblem(t *testing.T, rec *httptest.ResponseRecorder) Problem {
	t.Helper()
	if ct := rec.Header().Get("Content-Type"); ct != ProblemContentType {
		t.Fatalf("Content-Type = %q, want %q", ct, ProblemContentType)
	}
	var p Problem
	if err := json.NewDecoder(rec.Body).Decode(&p); err != nil {
		t.Fatalf("decode problem: %v", err)
	}
	return p
}

func serve(h http.Handler, method, path string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(method, path, nil))
	return rec
}

func TestHandleSuccessWritesNothingExtra(t *testing.T) {
	h := Handle(discardLogger(), func(w http.ResponseWriter, _ *http.Request) error {
		return WriteJSON(w, http.StatusCreated, map[string]string{"id": "1"})
	})
	rec := serve(h, http.MethodPost, "/things")
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d", rec.Code)
	}
	if got := strings.TrimSpace(rec.Body.String()); got != `{"id":"1"}` {
		t.Errorf("body = %s", got)
	}
}

func TestHandleMapsAPIError(t *testing.T) {
	h := Handle(discardLogger(), NotFound)
	rec := serve(h, http.MethodGet, "/missing")
	p := decodeProblem(t, rec)
	if rec.Code != http.StatusNotFound || p.Status != http.StatusNotFound {
		t.Fatalf("status = %d / %d", rec.Code, p.Status)
	}
	if p.Type != "https://openboard.dev/problems/not-found" || p.Instance != "/missing" {
		t.Errorf("unexpected problem: %+v", p)
	}
}

func TestHandleHidesUnknownErrorsAndLogsThem(t *testing.T) {
	var logs bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logs, nil))
	h := Handle(logger, func(http.ResponseWriter, *http.Request) error {
		return errors.New("pq: secret table exploded")
	})
	rec := serve(h, http.MethodGet, "/boom")
	p := decodeProblem(t, rec)
	if rec.Code != http.StatusInternalServerError || p.Title != "Internal server error" {
		t.Fatalf("unexpected response %d %+v", rec.Code, p)
	}
	if strings.Contains(rec.Body.String(), "secret") {
		t.Error("internal error detail leaked to client")
	}
	if !strings.Contains(logs.String(), "secret table exploded") {
		t.Error("internal error was not logged")
	}
}

func TestServiceUnavailableKeepsCauseOutOfBody(t *testing.T) {
	cause := errors.New("db down")
	err := ErrServiceUnavailable(cause)
	if !errors.Is(err, cause) {
		t.Error("cause should be reachable through Unwrap")
	}
	if err.Error() != "Service unavailable: db down" {
		t.Errorf("Error() = %q", err.Error())
	}
	rec := serve(Handle(discardLogger(), func(http.ResponseWriter, *http.Request) error { return err }), http.MethodGet, "/readyz")
	if strings.Contains(rec.Body.String(), "db down") {
		t.Error("cause leaked to client")
	}
}

func TestErrorWithoutCause(t *testing.T) {
	if got := ErrNotFound().Error(); got != "Resource not found" {
		t.Errorf("Error() = %q", got)
	}
}

func TestRouteMatchesMethodAndRejectsOthers(t *testing.T) {
	mux := http.NewServeMux()
	ok := func(w http.ResponseWriter, _ *http.Request) error {
		w.WriteHeader(http.StatusNoContent)
		return nil
	}
	Route(mux, discardLogger(), http.MethodGet, "/ping", ok)
	Route(mux, discardLogger(), http.MethodPost, "/items", ok)

	if rec := serve(mux, http.MethodGet, "/ping"); rec.Code != http.StatusNoContent {
		t.Errorf("GET /ping = %d", rec.Code)
	}
	cases := map[string]string{"/ping": "GET, HEAD", "/items": "POST"}
	for path, allow := range cases {
		rec := serve(mux, http.MethodDelete, path)
		p := decodeProblem(t, rec)
		if p.Status != http.StatusMethodNotAllowed || rec.Header().Get("Allow") != allow {
			t.Errorf("DELETE %s = %d, Allow %q", path, p.Status, rec.Header().Get("Allow"))
		}
	}
}

func TestRecoverTurnsPanicInto500(t *testing.T) {
	h := Recover(discardLogger(), http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic("kaboom")
	}))
	rec := serve(h, http.MethodGet, "/panic")
	if p := decodeProblem(t, rec); p.Status != http.StatusInternalServerError {
		t.Errorf("status = %d", p.Status)
	}
}

func TestRecoverRethrowsAbortHandler(t *testing.T) {
	h := Recover(discardLogger(), http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic(http.ErrAbortHandler)
	}))
	defer func() {
		rec := recover()
		if !isAbortHandler(rec) {
			t.Errorf("recovered %v, want ErrAbortHandler", rec)
		}
	}()
	serve(h, http.MethodGet, "/abort")
}

func TestRecoverPassesThrough(t *testing.T) {
	h := Recover(discardLogger(), http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))
	if rec := serve(h, http.MethodGet, "/"); rec.Code != http.StatusAccepted {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestAccessLogRecordsStatus(t *testing.T) {
	var logs bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logs, nil))
	h := AccessLog(logger, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))
	serve(h, http.MethodGet, "/tea")

	var entry map[string]any
	if err := json.Unmarshal(logs.Bytes(), &entry); err != nil {
		t.Fatalf("decode log: %v", err)
	}
	if entry["status"] != float64(http.StatusTeapot) || entry["path"] != "/tea" || entry["method"] != "GET" {
		t.Errorf("unexpected log entry: %v", entry)
	}
	if _, ok := entry["duration_ms"]; !ok {
		t.Error("missing duration_ms")
	}
}

func TestStatusRecorderUnwrap(t *testing.T) {
	inner := httptest.NewRecorder()
	rec := &statusRecorder{ResponseWriter: inner}
	if rec.Unwrap() != inner {
		t.Error("Unwrap should return the wrapped writer")
	}
}
