package identityhttp

import (
	"context"
	"errors"
	"net/http"

	"github.com/elyosemite/openboard/backend/internal/identity/application"
	"github.com/elyosemite/openboard/backend/internal/identity/domain"
	"github.com/elyosemite/openboard/backend/internal/platform/httpx"
)

// CookieName is the session cookie. The __Host- prefix makes browsers require
// Secure, Path=/ and no Domain, so no subdomain can set or read it.
const CookieName = "__Host-openboard_session"

type principalKey struct{}

// PrincipalFrom returns the authenticated Principal stored by the middleware.
func PrincipalFrom(ctx context.Context) (application.Principal, bool) {
	p, ok := ctx.Value(principalKey{}).(application.Principal)
	return p, ok
}

// Authenticated requires a valid session and no pending Temporary Password.
// Every endpoint outside the identity module is wrapped with it.
func (h *Handler) Authenticated(next httpx.HandlerFunc) httpx.HandlerFunc {
	return h.withSession(func(w http.ResponseWriter, r *http.Request) error {
		if p, _ := PrincipalFrom(r.Context()); p.User.MustChangePassword {
			return errPasswordChangeRequired()
		}
		return next(w, r)
	})
}

// withSession requires a valid session but lets a User with a pending
// Temporary Password through, so they can read their profile and change it.
func (h *Handler) withSession(next httpx.HandlerFunc) httpx.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		token := sessionToken(r)
		p, err := h.svc.Authenticate(r.Context(), token)
		if err != nil {
			if token != "" && errors.Is(err, domain.ErrUnauthenticated) {
				clearSessionCookie(w)
			}
			return mapError(err)
		}
		return next(w, r.WithContext(context.WithValue(r.Context(), principalKey{}, p)))
	}
}

func sessionToken(r *http.Request) string {
	c, err := r.Cookie(CookieName)
	if err != nil {
		return ""
	}
	return c.Value
}

func setSessionCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, newSessionCookie(token, int(domain.SessionLifetime.Seconds())))
}

func clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, newSessionCookie("", -1))
}

func newSessionCookie(value string, maxAge int) *http.Cookie {
	return &http.Cookie{
		Name:     CookieName,
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
}
