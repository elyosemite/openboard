package organizationshttp

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/elyosemite/openboard/backend/internal/organizations/application"
	"github.com/elyosemite/openboard/backend/internal/organizations/domain"
	"github.com/elyosemite/openboard/backend/internal/platform/httpx"
)

var administrator = domain.Actor{UserID: "00000000-0000-0000-0000-00000000000a", IsAdministrator: true}

type fakeService struct {
	err   error
	actor domain.Actor
	in    application.CreateOrganizationInput
}

func (f *fakeService) CreateOrganization(_ context.Context, actor domain.Actor, in application.CreateOrganizationInput) (domain.Organization, error) {
	f.actor, f.in = actor, in
	if f.err != nil {
		return domain.Organization{}, f.err
	}
	return domain.Organization{ID: "o1", Name: in.Name, CreatedAt: time.Date(2026, 11, 17, 12, 0, 0, 0, time.UTC)}, nil
}

// fakeAuth lets requests through as administrator unless signedOut is set.
func fakeAuth(signedOut bool) Auth {
	return Auth{
		Authenticated: func(next httpx.HandlerFunc) httpx.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) error {
				if signedOut {
					return &httpx.Error{Status: http.StatusUnauthorized, Slug: "unauthenticated", Title: "Not signed in"}
				}
				return next(w, r)
			}
		},
		Actor: func(context.Context) domain.Actor { return administrator },
	}
}

func post(svc Service, auth Auth, body string) *httptest.ResponseRecorder {
	mux := http.NewServeMux()
	NewHandler(svc, auth, slog.New(slog.NewTextHandler(io.Discard, nil))).Register(mux)
	r := httptest.NewRequest(http.MethodPost, "/api/v1/organizations", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, r)
	return rec
}

func problemOf(t *testing.T, rec *httptest.ResponseRecorder) httpx.Problem {
	t.Helper()
	var p httpx.Problem
	if err := json.NewDecoder(rec.Body).Decode(&p); err != nil {
		t.Fatalf("decode problem: %v", err)
	}
	return p
}

func TestCreateOrganization(t *testing.T) {
	svc := &fakeService{}
	rec := post(svc, fakeAuth(false), `{"name":"Acme","owner_user_id":"u2"}`)

	if rec.Code != http.StatusCreated || rec.Header().Get("Location") != "/api/v1/organizations/o1" {
		t.Fatalf("status %d, Location %q", rec.Code, rec.Header().Get("Location"))
	}
	if svc.actor != administrator || svc.in != (application.CreateOrganizationInput{Name: "Acme", OwnerUserID: "u2"}) {
		t.Errorf("service called with %+v, %+v", svc.actor, svc.in)
	}
	want := `{"id":"o1","name":"Acme","created_at":"2026-11-17T12:00:00Z"}`
	if got := strings.TrimSpace(rec.Body.String()); got != want {
		t.Errorf("body %s, want %s", got, want)
	}
}

func TestCreateOrganizationRequiresSession(t *testing.T) {
	rec := post(&fakeService{}, fakeAuth(true), `{"name":"Acme"}`)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status %d, want 401", rec.Code)
	}
}

func TestCreateOrganizationRejectsMalformedBody(t *testing.T) {
	if rec := post(&fakeService{}, fakeAuth(false), `{"name":`); rec.Code != http.StatusBadRequest {
		t.Errorf("status %d, want 400", rec.Code)
	}
}

func TestCreateOrganizationErrors(t *testing.T) {
	cases := map[string]struct {
		err    error
		status int
		slug   string
		fields []string
	}{
		"forbidden":     {domain.ErrForbidden, http.StatusForbidden, "forbidden", nil},
		"invalid name":  {domain.ErrInvalidName, http.StatusUnprocessableEntity, "validation-error", []string{"name:invalid"}},
		"unknown owner": {domain.ErrOwnerNotFound, http.StatusUnprocessableEntity, "validation-error", []string{"owner_user_id:not_found"}},
		"both invalid": {errors.Join(domain.ErrInvalidName, domain.ErrInvalidUserID), http.StatusUnprocessableEntity,
			"validation-error", []string{"name:invalid", "owner_user_id:invalid"}},
		"internal": {errors.New("store down"), http.StatusInternalServerError, "internal-error", nil},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			rec := post(&fakeService{err: tc.err}, fakeAuth(false), `{"name":"Acme"}`)
			p := problemOf(t, rec)
			if rec.Code != tc.status || !strings.HasSuffix(p.Type, "/"+tc.slug) {
				t.Fatalf("status %d, type %q; want %d %s", rec.Code, p.Type, tc.status, tc.slug)
			}
			var got []string
			for _, fe := range p.Errors {
				got = append(got, fe.Field+":"+fe.Code)
			}
			if strings.Join(got, ",") != strings.Join(tc.fields, ",") {
				t.Errorf("fields %v, want %v", got, tc.fields)
			}
		})
	}
}
