// Package health exposes liveness and readiness probes.
package health

import (
	"context"
	"fmt"
	"net/http"

	"github.com/elyosemite/openboard/backend/internal/platform/httpx"
)

// Checker verifies one dependency the API needs before it can serve traffic.
type Checker interface {
	Name() string
	Check(ctx context.Context) error
}

// Status is the body returned by healthy probes.
type Status struct {
	Status string `json:"status"`
}

// Handler serves /healthz and /readyz.
type Handler struct {
	checkers []Checker
}

// NewHandler returns a handler whose readiness depends on all checkers passing.
func NewHandler(checkers ...Checker) *Handler {
	return &Handler{checkers: checkers}
}

// Live reports that the process is up. It never checks dependencies.
func (h *Handler) Live(w http.ResponseWriter, _ *http.Request) error {
	return httpx.WriteJSON(w, http.StatusOK, Status{Status: "ok"})
}

// Ready reports whether every dependency is reachable.
func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) error {
	for _, c := range h.checkers {
		if err := c.Check(r.Context()); err != nil {
			return httpx.ErrServiceUnavailable(fmt.Errorf("%s: %w", c.Name(), err))
		}
	}
	return httpx.WriteJSON(w, http.StatusOK, Status{Status: "ready"})
}
