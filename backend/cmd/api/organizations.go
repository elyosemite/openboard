package main

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	identityhttp "github.com/elyosemite/openboard/backend/internal/identity/http"
	"github.com/elyosemite/openboard/backend/internal/organizations/application"
	"github.com/elyosemite/openboard/backend/internal/organizations/domain"
	organizationshttp "github.com/elyosemite/openboard/backend/internal/organizations/http"
	"github.com/elyosemite/openboard/backend/internal/organizations/infrastructure/postgres"
)

// setupOrganizations wires the Organization module, authenticating its
// requests through the identity module.
func setupOrganizations(logger *slog.Logger, pool *pgxpool.Pool, identity *identityhttp.Handler) *organizationshttp.Handler {
	svc := application.NewService(postgres.NewOrganizations(pool))
	auth := organizationshttp.Auth{Authenticated: identity.Authenticated, Actor: actorOf}
	return organizationshttp.NewHandler(svc, auth, logger)
}

// setupWorkspaces wires the Workspace endpoints of the Organization module.
func setupWorkspaces(logger *slog.Logger, pool *pgxpool.Pool, identity *identityhttp.Handler) *organizationshttp.WorkspaceHandler {
	svc := application.NewWorkspaceService(postgres.NewWorkspaces(pool))
	auth := organizationshttp.Auth{Authenticated: identity.Authenticated, Actor: actorOf}
	return organizationshttp.NewWorkspaceHandler(svc, auth, logger)
}

// actorOf translates the identity Principal into the Organization module's Actor.
func actorOf(ctx context.Context) domain.Actor {
	p, _ := identityhttp.PrincipalFrom(ctx)
	return domain.Actor{UserID: domain.UserID(p.User.ID), IsAdministrator: p.User.IsAdministrator}
}
