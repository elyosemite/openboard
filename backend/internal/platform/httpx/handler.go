package httpx

import (
	"errors"
	"log/slog"
	"net/http"
)

// HandlerFunc is an HTTP handler that reports failures by returning an error,
// so error-to-response mapping lives in one place.
type HandlerFunc func(http.ResponseWriter, *http.Request) error

// Handle adapts h to http.Handler, turning returned errors into problem responses.
func Handle(logger *slog.Logger, h HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			WriteError(logger, w, r, err)
		}
	})
}

// WriteError maps err to a problem response. Errors that are not *Error become
// a generic 500 so internal details never leak; 5xx causes are logged.
func WriteError(logger *slog.Logger, w http.ResponseWriter, r *http.Request, err error) {
	var apiErr *Error
	if !errors.As(err, &apiErr) {
		apiErr = errInternal(err)
	}
	if apiErr.Status >= http.StatusInternalServerError {
		logger.ErrorContext(r.Context(), "request failed",
			"method", r.Method, "path", r.URL.Path, "status", apiErr.Status, "error", err)
	}
	writeProblem(w, r, apiErr)
}

// Route registers h for method+pattern and answers any other method on the
// same pattern with a 405 problem instead of the mux's plain-text reply.
func Route(mux *http.ServeMux, logger *slog.Logger, method, pattern string, h HandlerFunc) {
	mux.Handle(method+" "+pattern, Handle(logger, h))
	allow := method
	if method == http.MethodGet {
		allow = "GET, HEAD"
	}
	mux.Handle(pattern, Handle(logger, func(w http.ResponseWriter, _ *http.Request) error {
		w.Header().Set("Allow", allow)
		return ErrMethodNotAllowed()
	}))
}

// NotFound answers every unmatched route with a 404 problem.
func NotFound(http.ResponseWriter, *http.Request) error {
	return ErrNotFound()
}
