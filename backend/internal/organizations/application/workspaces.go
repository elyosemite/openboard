package application

import (
	"context"
	"errors"

	"github.com/elyosemite/openboard/backend/internal/organizations/domain"
)

// NewWorkspace holds what is needed to create a Workspace.
type NewWorkspace struct {
	OrganizationID domain.OrganizationID
	Name           string
	Visibility     domain.Visibility
	CreatorID      domain.UserID
}

// WorkspaceRepository stores Workspaces.
type WorkspaceRepository interface {
	// RoleIn returns the role scope acts with in the Organization, or
	// domain.ErrUnknownOrganization when scope does not see it.
	RoleIn(ctx context.Context, id domain.OrganizationID, scope Scope) (domain.Role, error)
	// Create stores the Workspace and, when the creator holds an Organization
	// Membership, makes them its Workspace Owner, atomically. It returns
	// domain.ErrWorkspaceNameTaken when the Organization already uses the name.
	Create(ctx context.Context, workspace NewWorkspace) (domain.Workspace, error)
	// Find returns the Workspace with id if scope can enter it, or
	// domain.ErrWorkspaceNotFound.
	Find(ctx context.Context, id domain.WorkspaceID, scope Scope) (domain.Workspace, error)
	// List returns one page of the Organization's Workspaces scope can enter,
	// ordered by name ignoring case, and how many there are in total.
	List(ctx context.Context, id domain.OrganizationID, scope Scope, page PageRequest) ([]domain.Workspace, int, error)
}

// CreateWorkspaceInput is the raw request to create a Workspace.
type CreateWorkspaceInput struct {
	OrganizationID string
	Name           string
	Visibility     string
}

// WorkspaceService runs the Workspace use cases.
type WorkspaceService struct {
	workspaces WorkspaceRepository
}

// NewWorkspaceService builds a WorkspaceService from its repository.
func NewWorkspaceService(workspaces WorkspaceRepository) *WorkspaceService {
	return &WorkspaceService{workspaces: workspaces}
}

// CreateWorkspace creates a Workspace in an Organization the Actor manages.
// Invalid fields are all reported at once, joined.
func (s *WorkspaceService) CreateWorkspace(ctx context.Context, actor domain.Actor, in CreateWorkspaceInput) (domain.Workspace, error) {
	organizationID, refErr := domain.ParseOrganizationRef(in.OrganizationID)
	name, nameErr := domain.ParseName(in.Name)
	visibility, visibilityErr := domain.ParseVisibility(in.Visibility)
	if err := errors.Join(refErr, nameErr, visibilityErr); err != nil {
		return domain.Workspace{}, err
	}
	role, err := s.workspaces.RoleIn(ctx, organizationID, scopeOf(actor))
	if err != nil {
		return domain.Workspace{}, err
	}
	if !role.CanCreateWorkspace() {
		return domain.Workspace{}, domain.ErrForbidden
	}
	return s.workspaces.Create(ctx, NewWorkspace{
		OrganizationID: organizationID,
		Name:           name,
		Visibility:     visibility,
		CreatorID:      actor.UserID,
	})
}

// GetWorkspace returns a Workspace the Actor can enter. Unknown IDs and
// Workspaces the Actor cannot enter both give domain.ErrWorkspaceNotFound.
func (s *WorkspaceService) GetWorkspace(ctx context.Context, actor domain.Actor, rawID string) (domain.Workspace, error) {
	id, err := domain.ParseWorkspaceID(rawID)
	if err != nil {
		return domain.Workspace{}, err
	}
	return s.workspaces.Find(ctx, id, scopeOf(actor))
}

// ListWorkspaces returns one page of the Workspaces the Actor can enter in an
// Organization they can see, and how many there are in total.
func (s *WorkspaceService) ListWorkspaces(ctx context.Context, actor domain.Actor, rawOrganizationID string, page PageRequest) ([]domain.Workspace, int, error) {
	organizationID, err := domain.ParseOrganizationRef(rawOrganizationID)
	if err != nil {
		return nil, 0, err
	}
	scope := scopeOf(actor)
	if _, err := s.workspaces.RoleIn(ctx, organizationID, scope); err != nil {
		return nil, 0, err
	}
	return s.workspaces.List(ctx, organizationID, scope, page)
}
