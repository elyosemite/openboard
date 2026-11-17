// Package organizationshttp exposes the Organization use cases over HTTP.
// Authentication comes from the identity module, injected by the caller, so
// this package never imports another module.
package organizationshttp

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/elyosemite/openboard/backend/internal/organizations/application"
	"github.com/elyosemite/openboard/backend/internal/organizations/domain"
	"github.com/elyosemite/openboard/backend/internal/platform/httpx"
)

// Service is the part of the Organization application the HTTP layer uses.
type Service interface {
	CreateOrganization(ctx context.Context, actor domain.Actor, in application.CreateOrganizationInput) (domain.Organization, error)
}

// Auth authenticates requests and tells who is behind them.
type Auth struct {
	// Authenticated rejects requests without a valid session.
	Authenticated func(next httpx.HandlerFunc) httpx.HandlerFunc
	// Actor returns the signed-in User of a request Authenticated let through.
	Actor func(ctx context.Context) domain.Actor
}

// Handler serves the Organization endpoints.
type Handler struct {
	svc    Service
	auth   Auth
	logger *slog.Logger
}

// NewHandler builds a Handler.
func NewHandler(svc Service, auth Auth, logger *slog.Logger) *Handler {
	return &Handler{svc: svc, auth: auth, logger: logger}
}

// Register adds the Organization routes to mux.
func (h *Handler) Register(mux *http.ServeMux) {
	httpx.Route(mux, h.logger, http.MethodPost, "/api/v1/organizations", h.auth.Authenticated(h.create))
}

type createRequest struct {
	Name        string `json:"name"`
	OwnerUserID string `json:"owner_user_id"`
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) error {
	var req createRequest
	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		return err
	}
	org, err := h.svc.CreateOrganization(r.Context(), h.auth.Actor(r.Context()),
		application.CreateOrganizationInput{Name: req.Name, OwnerUserID: req.OwnerUserID})
	if err != nil {
		return mapError(err)
	}
	w.Header().Set("Location", "/api/v1/organizations/"+string(org.ID))
	return httpx.WriteJSON(w, http.StatusCreated, organizationResponse{ID: org.ID, Name: org.Name, CreatedAt: org.CreatedAt})
}

type organizationResponse struct {
	ID        domain.OrganizationID `json:"id"`
	Name      string                `json:"name"`
	CreatedAt time.Time             `json:"created_at"`
}

// fieldErrors maps each domain validation error to the request field it concerns.
var fieldErrors = []struct {
	err   error
	field httpx.FieldError
}{
	{domain.ErrInvalidName, httpx.FieldError{Field: "name", Code: "invalid", Message: domain.ErrInvalidName.Error()}},
	{domain.ErrInvalidUserID, httpx.FieldError{Field: "owner_user_id", Code: "invalid", Message: domain.ErrInvalidUserID.Error()}},
	{domain.ErrOwnerNotFound, httpx.FieldError{Field: "owner_user_id", Code: "not_found", Message: domain.ErrOwnerNotFound.Error()}},
}

// mapError turns Organization domain errors into problem responses; anything
// else passes through and becomes a 500.
func mapError(err error) error {
	if errors.Is(err, domain.ErrForbidden) {
		return &httpx.Error{Status: http.StatusForbidden, Slug: "forbidden", Title: "Not allowed"}
	}
	var fields []httpx.FieldError
	for _, fe := range fieldErrors {
		if errors.Is(err, fe.err) {
			fields = append(fields, fe.field)
		}
	}
	if len(fields) > 0 {
		return httpx.ErrValidation(fields...)
	}
	return err
}
