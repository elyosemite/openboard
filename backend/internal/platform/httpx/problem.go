// Package httpx holds the HTTP building blocks shared by every module:
// RFC 9457 problem responses, centralized error handling, routing helpers
// and middleware.
package httpx

import (
	"encoding/json"
	"net/http"
)

// ProblemContentType is the media type of RFC 9457 problem responses.
const ProblemContentType = "application/problem+json"

const problemTypeBase = "https://openboard.dev/problems/"

// Problem is the RFC 9457 body returned for every API error.
type Problem struct {
	Type     string       `json:"type"`
	Title    string       `json:"title"`
	Status   int          `json:"status"`
	Detail   string       `json:"detail,omitempty"`
	Instance string       `json:"instance,omitempty"`
	Errors   []FieldError `json:"errors,omitempty"`
}

// FieldError is one field-level detail of a validation-error problem.
type FieldError struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error is returned by handlers to produce a specific problem response.
// Cause is logged for 5xx responses and never exposed to the client.
type Error struct {
	Status int
	Slug   string
	Title  string
	Detail string
	Errors []FieldError
	Cause  error
}

func (e *Error) Error() string {
	if e.Cause != nil {
		return e.Title + ": " + e.Cause.Error()
	}
	return e.Title
}

func (e *Error) Unwrap() error { return e.Cause }

func (e *Error) problem(r *http.Request) Problem {
	return Problem{
		Type:     problemTypeBase + e.Slug,
		Title:    e.Title,
		Status:   e.Status,
		Detail:   e.Detail,
		Instance: r.URL.Path,
		Errors:   e.Errors,
	}
}

// ErrNotFound reports that the requested resource does not exist.
func ErrNotFound() *Error {
	return &Error{Status: http.StatusNotFound, Slug: "not-found", Title: "Resource not found"}
}

// ErrMethodNotAllowed reports that the resource exists but does not accept the method.
func ErrMethodNotAllowed() *Error {
	return &Error{Status: http.StatusMethodNotAllowed, Slug: "method-not-allowed", Title: "Method not allowed"}
}

// ErrServiceUnavailable reports a dependency that is not ready; cause is logged only.
func ErrServiceUnavailable(cause error) *Error {
	return &Error{
		Status: http.StatusServiceUnavailable,
		Slug:   "service-unavailable",
		Title:  "Service unavailable",
		Cause:  cause,
	}
}

// ErrMalformedRequest reports a request body that cannot be decoded.
func ErrMalformedRequest() *Error {
	return &Error{
		Status: http.StatusBadRequest,
		Slug:   "malformed-request",
		Title:  "Malformed request",
		Detail: "The request body is not valid JSON for this operation.",
	}
}

// ErrUnsupportedMediaType reports a request body not sent as application/json.
func ErrUnsupportedMediaType() *Error {
	return &Error{
		Status: http.StatusUnsupportedMediaType,
		Slug:   "unsupported-media-type",
		Title:  "Unsupported media type",
		Detail: "Send the request body as application/json.",
	}
}

// ErrValidation reports request fields that break the operation's rules.
func ErrValidation(fields ...FieldError) *Error {
	return &Error{
		Status: http.StatusUnprocessableEntity,
		Slug:   "validation-error",
		Title:  "Validation failed",
		Errors: fields,
	}
}

// ErrCrossOrigin reports a cross-origin browser request that tried to change state.
func ErrCrossOrigin() *Error {
	return &Error{
		Status: http.StatusForbidden,
		Slug:   "cross-origin-request",
		Title:  "Cross-origin request rejected",
	}
}

func errInternal(cause error) *Error {
	return &Error{
		Status: http.StatusInternalServerError,
		Slug:   "internal-error",
		Title:  "Internal server error",
		Detail: "An unexpected error occurred.",
		Cause:  cause,
	}
}

func writeProblem(w http.ResponseWriter, r *http.Request, e *Error) {
	w.Header().Set("Content-Type", ProblemContentType)
	w.WriteHeader(e.Status)
	_ = json.NewEncoder(w).Encode(e.problem(r))
}

// WriteJSON writes v as a JSON response with the given status.
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
