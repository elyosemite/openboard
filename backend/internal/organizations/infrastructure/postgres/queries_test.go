package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/elyosemite/openboard/backend/internal/organizations/application"
	"github.com/elyosemite/openboard/backend/internal/organizations/domain"
)

// seed creates Organizations owned by owner and returns them by name.
func seed(t *testing.T, owner domain.UserID, names ...string) map[string]domain.Organization {
	t.Helper()
	orgs := map[string]domain.Organization{}
	for _, name := range names {
		org, err := NewOrganizations(pool).Create(context.Background(),
			application.NewOrganization{Name: name, OwnerID: owner, DefaultWorkspaceName: "Default"})
		if err != nil {
			t.Fatal(err)
		}
		orgs[name] = org
	}
	return orgs
}

func addUser(t *testing.T, username string) domain.UserID {
	t.Helper()
	var id string
	err := pool.QueryRow(context.Background(),
		`INSERT INTO users (username, name, password_hash) VALUES ($1, $1, 'hash') RETURNING id::text`, username).Scan(&id)
	if err != nil {
		t.Fatal(err)
	}
	return domain.UserID(id)
}

func TestFindSeesOnlyWhatScopeAllows(t *testing.T) {
	marina := createUser(t, "marina")
	julia := addUser(t, "julia")
	acme := seed(t, marina, "Acme")["Acme"]
	repo := NewOrganizations(pool)
	ctx := context.Background()

	cases := map[string]struct {
		scope application.Scope
		found bool
	}{
		"member":           {application.Scope{UserID: marina}, true},
		"not a member":     {application.Scope{UserID: julia}, false},
		"sees all":         {application.Scope{UserID: julia, All: true}, true},
		"all without user": {application.Scope{All: true}, true},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			org, err := repo.Find(ctx, acme.ID, tc.scope)
			if tc.found && (err != nil || org != acme) {
				t.Fatalf("Find = %+v, %v; want %+v", org, err, acme)
			}
			if !tc.found && !errors.Is(err, domain.ErrNotFound) {
				t.Fatalf("err = %v, want ErrNotFound", err)
			}
		})
	}
	if _, err := repo.Find(ctx, "00000000-0000-0000-0000-000000000000", application.Scope{All: true}); !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("unknown id err = %v, want ErrNotFound", err)
	}
}

func TestListPagesByNameWithinScope(t *testing.T) {
	marina := createUser(t, "marina")
	julia := addUser(t, "julia")
	seed(t, marina, "Charlie", "alpha", "Bravo")
	seed(t, julia, "Delta")
	repo := NewOrganizations(pool)
	ctx := context.Background()

	names := func(orgs []domain.Organization) []string {
		var out []string
		for _, o := range orgs {
			out = append(out, o.Name)
		}
		return out
	}
	cases := map[string]struct {
		scope application.Scope
		page  application.PageRequest
		want  []string
		total int
	}{
		"all, first page":   {application.Scope{All: true}, application.PageRequest{Limit: 2}, []string{"alpha", "Bravo"}, 4},
		"all, second page":  {application.Scope{All: true}, application.PageRequest{Limit: 2, Offset: 2}, []string{"Charlie", "Delta"}, 4},
		"member only":       {application.Scope{UserID: julia}, application.PageRequest{Limit: 20}, []string{"Delta"}, 1},
		"past the last one": {application.Scope{UserID: julia}, application.PageRequest{Limit: 20, Offset: 20}, nil, 1},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			orgs, total, err := repo.List(ctx, tc.scope, tc.page)
			if err != nil || total != tc.total || len(orgs) != len(tc.want) {
				t.Fatalf("List = %v, %d, %v; want %v, %d", names(orgs), total, err, tc.want, tc.total)
			}
			for i, name := range names(orgs) {
				if name != tc.want[i] {
					t.Errorf("List = %v, want %v", names(orgs), tc.want)
				}
			}
		})
	}
}

func TestListReportsStorageFailures(t *testing.T) {
	createUser(t, "marina")
	if _, _, err := NewOrganizations(pool).List(context.Background(), application.Scope{All: true}, application.PageRequest{Limit: -1}); err == nil {
		t.Error("a negative LIMIT must fail")
	}
}
