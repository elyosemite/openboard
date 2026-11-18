package application

import (
	"context"
	"errors"
	"testing"

	"github.com/elyosemite/openboard/backend/internal/organizations/domain"
)

const acmeID = "0b5e3c1a-9f2d-4e6b-8a7c-1d2e3f4a5b6c"

var contributor = domain.Actor{UserID: marinaID}

func TestGetOrganizationScopesByActor(t *testing.T) {
	cases := map[string]struct {
		actor domain.Actor
		want  Scope
	}{
		"Administrator sees all": {administrator, Scope{UserID: adminID, All: true}},
		"others need Membership": {contributor, Scope{UserID: marinaID}},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			repo := &fakeRepository{}
			org, err := NewService(repo).GetOrganization(context.Background(), tc.actor, " 0B5E3C1A-9F2D-4E6B-8A7C-1D2E3F4A5B6C ")
			if err != nil || org.ID != acmeID {
				t.Fatalf("org %+v, err %v", org, err)
			}
			if repo.scope != tc.want || repo.foundID != acmeID {
				t.Errorf("Find(%q, %+v), want (%q, %+v)", repo.foundID, repo.scope, acmeID, tc.want)
			}
		})
	}
}

func TestGetOrganizationWithInvalidIDIsNotFound(t *testing.T) {
	repo := &fakeRepository{}
	_, err := NewService(repo).GetOrganization(context.Background(), administrator, "acme")
	if !errors.Is(err, domain.ErrNotFound) || repo.foundID != "" {
		t.Errorf("err = %v, Find called with %q; want ErrNotFound without a lookup", err, repo.foundID)
	}
}

func TestListOrganizationsScopesByActor(t *testing.T) {
	repo := &fakeRepository{total: 42}
	page := PageRequest{Limit: 20, Offset: 40}
	orgs, total, err := NewService(repo).ListOrganizations(context.Background(), contributor, page)
	if err != nil || total != 42 || len(orgs) != 1 {
		t.Fatalf("orgs %+v, total %d, err %v", orgs, total, err)
	}
	if repo.scope != (Scope{UserID: marinaID}) || repo.page != page {
		t.Errorf("List(%+v, %+v)", repo.scope, repo.page)
	}
}
