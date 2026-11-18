package application

import (
	"context"

	"github.com/elyosemite/openboard/backend/internal/organizations/domain"
)

// GetOrganization returns an Organization the Actor can see. Unknown IDs and
// Organizations the Actor cannot see both give domain.ErrNotFound.
func (s *Service) GetOrganization(ctx context.Context, actor domain.Actor, rawID string) (domain.Organization, error) {
	id, err := domain.ParseOrganizationID(rawID)
	if err != nil {
		return domain.Organization{}, err
	}
	return s.organizations.Find(ctx, id, scopeOf(actor))
}

// ListOrganizations returns one page of the Organizations the Actor can see,
// ordered by name, and how many they see in total.
func (s *Service) ListOrganizations(ctx context.Context, actor domain.Actor, page PageRequest) ([]domain.Organization, int, error) {
	return s.organizations.List(ctx, scopeOf(actor), page)
}
