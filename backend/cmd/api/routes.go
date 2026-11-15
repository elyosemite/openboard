package main

import (
	"log/slog"
	"net/http"

	"github.com/elyosemite/openboard/backend/internal/platform/health"
	"github.com/elyosemite/openboard/backend/internal/platform/httpx"
)

// newRouter wires every route of the API. Operational probes live outside
// /api/v1 because they are not domain resources.
func newRouter(logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	probes := health.NewHandler()
	httpx.Route(mux, logger, http.MethodGet, "/healthz", probes.Live)
	httpx.Route(mux, logger, http.MethodGet, "/readyz", probes.Ready)

	mux.Handle("/", httpx.Handle(logger, httpx.NotFound))

	return httpx.AccessLog(logger, httpx.Recover(logger, mux))
}
