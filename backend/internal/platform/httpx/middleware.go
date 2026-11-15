package httpx

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// Recover turns a panicking handler into a 500 problem response.
func Recover(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			rec := recover()
			if rec == nil {
				return
			}
			if isAbortHandler(rec) {
				panic(rec)
			}
			WriteError(logger, w, r, fmt.Errorf("panic: %v", rec))
		}()
		next.ServeHTTP(w, r)
	})
}

// isAbortHandler reports the sentinel panic net/http uses to abort a response;
// it must propagate so the server can close the connection.
func isAbortHandler(rec any) bool {
	err, ok := rec.(error)
	return ok && errors.Is(err, http.ErrAbortHandler)
}

// AccessLog writes one structured log entry per request.
func AccessLog(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		logger.InfoContext(r.Context(), "http request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(status int) {
	s.status = status
	s.ResponseWriter.WriteHeader(status)
}

func (s *statusRecorder) Unwrap() http.ResponseWriter { return s.ResponseWriter }
