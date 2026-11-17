package application

import (
	"context"
	"errors"
	"testing"

	"github.com/elyosemite/openboard/backend/internal/organizations/domain"
)

const (
	adminID  = domain.UserID("00000000-0000-0000-0000-00000000000a")
	marinaID = "00000000-0000-0000-0000-0000000000b1"
)

var administrator = domain.Actor{UserID: adminID, IsAdministrator: true}

type fakeRepository struct {
	created []NewOrganization
	err     error
}

func (f *fakeRepository) Create(_ context.Context, org NewOrganization) (domain.Organization, error) {
	if f.err != nil {
		return domain.Organization{}, f.err
	}
	f.created = append(f.created, org)
	return domain.Organization{ID: "o1", Name: org.Name}, nil
}

func TestCreateOrganizationMakesCreatorOwnerByDefault(t *testing.T) {
	repo := &fakeRepository{}
	org, err := NewService(repo).CreateOrganization(context.Background(), administrator, CreateOrganizationInput{Name: "  Acme "})
	if err != nil {
		t.Fatal(err)
	}
	want := NewOrganization{Name: "Acme", OwnerID: adminID, DefaultWorkspaceName: "Default"}
	if len(repo.created) != 1 || repo.created[0] != want {
		t.Fatalf("created %+v, want %+v", repo.created, want)
	}
	if org.ID != "o1" || org.Name != "Acme" {
		t.Errorf("returned %+v", org)
	}
}

func TestCreateOrganizationWithNamedOwner(t *testing.T) {
	repo := &fakeRepository{}
	_, err := NewService(repo).CreateOrganization(context.Background(), administrator,
		CreateOrganizationInput{Name: "Acme", OwnerUserID: marinaID})
	if err != nil {
		t.Fatal(err)
	}
	if repo.created[0].OwnerID != marinaID {
		t.Errorf("owner = %q, want %q", repo.created[0].OwnerID, marinaID)
	}
}

func TestCreateOrganizationRejections(t *testing.T) {
	errStore := errors.New("store down")
	cases := map[string]struct {
		actor domain.Actor
		in    CreateOrganizationInput
		repo  error
		want  []error
	}{
		"not an Administrator": {domain.Actor{UserID: adminID}, CreateOrganizationInput{Name: "Acme"}, nil, []error{domain.ErrForbidden}},
		"blank name":           {administrator, CreateOrganizationInput{Name: " "}, nil, []error{domain.ErrInvalidName}},
		"bad owner id":         {administrator, CreateOrganizationInput{Name: "Acme", OwnerUserID: "x"}, nil, []error{domain.ErrInvalidUserID}},
		"both invalid": {administrator, CreateOrganizationInput{OwnerUserID: "x"}, nil,
			[]error{domain.ErrInvalidName, domain.ErrInvalidUserID}},
		"unknown owner": {administrator, CreateOrganizationInput{Name: "Acme", OwnerUserID: marinaID},
			domain.ErrOwnerNotFound, []error{domain.ErrOwnerNotFound}},
		"store failure": {administrator, CreateOrganizationInput{Name: "Acme"}, errStore, []error{errStore}},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			repo := &fakeRepository{err: tc.repo}
			_, err := NewService(repo).CreateOrganization(context.Background(), tc.actor, tc.in)
			for _, want := range tc.want {
				if !errors.Is(err, want) {
					t.Errorf("err = %v, want it to wrap %v", err, want)
				}
			}
		})
	}
}
