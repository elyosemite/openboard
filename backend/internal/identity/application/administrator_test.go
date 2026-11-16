package application

import (
	"context"
	"errors"
	"testing"

	"github.com/elyosemite/openboard/backend/internal/identity/domain"
)

func onlyUser(t *testing.T, h *harness) domain.User {
	t.Helper()
	if len(h.users.byID) != 1 {
		t.Fatalf("want exactly one user, got %d", len(h.users.byID))
	}
	for _, u := range h.users.byID {
		return u
	}
	return domain.User{}
}

func TestEnsureAdministratorDefaultsToTemporaryAdmin(t *testing.T) {
	h := newHarness(t)

	created, err := h.svc.EnsureAdministrator(context.Background(), AdministratorSeed{})
	if err != nil || !created {
		t.Fatalf("created=%v err=%v", created, err)
	}
	u := onlyUser(t, h)
	if u.Username != "admin" || u.PasswordHash != "h:admin" || !u.MustChangePassword || !u.IsAdministrator {
		t.Errorf("unexpected administrator: %+v", u)
	}
}

func TestEnsureAdministratorFromSeed(t *testing.T) {
	h := newHarness(t)
	seed := AdministratorSeed{Username: "Root.Ops", Password: "a long and fine passphrase"}

	if created, err := h.svc.EnsureAdministrator(context.Background(), seed); err != nil || !created {
		t.Fatalf("created=%v err=%v", created, err)
	}
	u := onlyUser(t, h)
	if u.Username != "root.ops" || u.MustChangePassword || !u.IsAdministrator {
		t.Errorf("unexpected administrator: %+v", u)
	}
}

func TestEnsureAdministratorDoesNothingWhenUsersExist(t *testing.T) {
	h := newHarness(t, gabriel())

	created, err := h.svc.EnsureAdministrator(context.Background(), AdministratorSeed{})
	if err != nil || created || len(h.users.byID) != 1 {
		t.Fatalf("created=%v err=%v users=%d", created, err, len(h.users.byID))
	}
}

func TestEnsureAdministratorRejectsBadSeeds(t *testing.T) {
	cases := map[string]struct {
		seed AdministratorSeed
		want error
	}{
		"only username": {AdministratorSeed{Username: "root"}, ErrIncompleteAdministratorSeed},
		"only password": {AdministratorSeed{Password: "a long and fine passphrase"}, ErrIncompleteAdministratorSeed},
		"bad username":  {AdministratorSeed{Username: "r", Password: "a long and fine passphrase"}, domain.ErrInvalidUsername},
		"weak password": {AdministratorSeed{Username: "root", Password: "short"}, domain.ErrWeakPassword},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			h := newHarness(t)
			if _, err := h.svc.EnsureAdministrator(context.Background(), tc.seed); !errors.Is(err, tc.want) {
				t.Fatalf("err = %v, want %v", err, tc.want)
			}
			if len(h.users.byID) != 0 {
				t.Error("no user should be created")
			}
		})
	}
}

func TestEnsureAdministratorInfrastructureFailures(t *testing.T) {
	cases := map[string]func(h *harness){
		"count":  func(h *harness) { h.users.countErr = errBoom },
		"hash":   func(h *harness) { h.hasher.hashErr = errBoom },
		"create": func(h *harness) { h.users.createErr = errBoom },
	}
	for name, breakIt := range cases {
		t.Run(name, func(t *testing.T) {
			h := newHarness(t)
			breakIt(h)
			if _, err := h.svc.EnsureAdministrator(context.Background(), AdministratorSeed{}); !errors.Is(err, errBoom) {
				t.Fatalf("err = %v, want boom", err)
			}
		})
	}
}
