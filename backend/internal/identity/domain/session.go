package domain

import "time"

// Session lifetime rules (ADR 0006).
const (
	SessionIdleTimeout = 7 * 24 * time.Hour
	SessionLifetime    = 30 * 24 * time.Hour
)

// Session is a signed-in User's server-side session. ID is the SHA-256 digest
// of the token held in the cookie, never the token itself.
type Session struct {
	ID         []byte
	UserID     UserID
	CreatedAt  time.Time
	LastUsedAt time.Time
	ExpiresAt  time.Time
}

// NewSession starts a session that lasts at most SessionLifetime.
func NewSession(id []byte, userID UserID, now time.Time) Session {
	return Session{
		ID:         id,
		UserID:     userID,
		CreatedAt:  now,
		LastUsedAt: now,
		ExpiresAt:  now.Add(SessionLifetime),
	}
}

// IsExpired reports whether the session passed its lifetime or sat idle too long.
func (s Session) IsExpired(now time.Time) bool {
	return !now.Before(s.ExpiresAt) || now.Sub(s.LastUsedAt) >= SessionIdleTimeout
}
