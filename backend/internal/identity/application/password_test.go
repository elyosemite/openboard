package application

import (
	"context"
	"errors"
	"testing"

	"github.com/elyosemite/openboard/backend/internal/identity/domain"
)

const newPassword = "a brand new passphrase"

func principalWithTwoSessions(t *testing.T) (*harness, Principal) {
	t.Helper()
	h, token := signedIn(t)
	if _, err := h.svc.SignIn(context.Background(), "gabriel", "correct horse battery"); err != nil {
		t.Fatal(err)
	}
	p, err := h.svc.Authenticate(context.Background(), token)
	if err != nil {
		t.Fatal(err)
	}
	return h, p
}

func TestChangePasswordClearsTemporaryAndEndsOtherSessions(t *testing.T) {
	h, p := principalWithTwoSessions(t)
	user := h.users.byID["u-gabriel"]
	user.MustChangePassword = true
	h.users.byID["u-gabriel"] = user
	p.User = user

	if err := h.svc.ChangePassword(context.Background(), p, "correct horse battery", newPassword); err != nil {
		t.Fatalf("ChangePassword: %v", err)
	}
	updated := h.users.byID["u-gabriel"]
	if updated.PasswordHash != "h:"+newPassword || updated.MustChangePassword {
		t.Errorf("user not updated: %+v", updated)
	}
	if _, kept := h.sessions.byID[string(p.SessionID)]; !kept || len(h.sessions.byID) != 1 {
		t.Errorf("only the current session should remain, got %d", len(h.sessions.byID))
	}
}

func TestChangePasswordRejections(t *testing.T) {
	h, p := principalWithTwoSessions(t)

	if err := h.svc.ChangePassword(context.Background(), p, "wrong", newPassword); !errors.Is(err, domain.ErrWrongCurrentPassword) {
		t.Errorf("wrong current: err = %v", err)
	}
	if err := h.svc.ChangePassword(context.Background(), p, "correct horse battery", "short"); !errors.Is(err, domain.ErrWeakPassword) {
		t.Errorf("weak: err = %v", err)
	}
	if h.users.byID["u-gabriel"].PasswordHash != "h:correct horse battery" {
		t.Error("password must not change on rejection")
	}
}

func TestChangePasswordInfrastructureFailures(t *testing.T) {
	cases := map[string]func(h *harness){
		"verify":         func(h *harness) { h.hasher.verifyErr = errBoom },
		"hash":           func(h *harness) { h.hasher.hashErr = errBoom },
		"update":         func(h *harness) { h.users.updateErr = errBoom },
		"other sessions": func(h *harness) { h.sessions.deleteOthErr = errBoom },
	}
	for name, breakIt := range cases {
		t.Run(name, func(t *testing.T) {
			h, p := principalWithTwoSessions(t)
			breakIt(h)
			if err := h.svc.ChangePassword(context.Background(), p, "correct horse battery", newPassword); !errors.Is(err, errBoom) {
				t.Fatalf("err = %v, want boom", err)
			}
		})
	}
}
