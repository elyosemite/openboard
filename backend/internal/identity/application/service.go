package application

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"time"
)

// touchInterval limits how often a session's last use is written back.
const touchInterval = time.Minute

const timingEqualizerPassword = "openboard-timing-equalizer"

// Deps are the collaborators of the Service. Now and Random default to the
// real clock and crypto/rand when nil.
type Deps struct {
	Users    UserRepository
	Sessions SessionRepository
	Hasher   PasswordHasher
	Throttle SignInThrottle
	Now      func() time.Time
	Random   io.Reader
}

// Service runs the identity use cases.
type Service struct {
	users    UserRepository
	sessions SessionRepository
	hasher   PasswordHasher
	throttle SignInThrottle
	now      func() time.Time
	random   io.Reader
	// dummyHash is verified when a Username does not exist, so a failed
	// sign-in takes the same time whether or not the User exists.
	dummyHash string
}

// NewService builds a Service from its dependencies.
func NewService(d Deps) (*Service, error) {
	dummy, err := d.Hasher.Hash(timingEqualizerPassword)
	if err != nil {
		return nil, fmt.Errorf("prepare timing equalizer: %w", err)
	}
	s := &Service{
		users:     d.Users,
		sessions:  d.Sessions,
		hasher:    d.Hasher,
		throttle:  d.Throttle,
		now:       d.Now,
		random:    d.Random,
		dummyHash: dummy,
	}
	if s.now == nil {
		s.now = time.Now
	}
	if s.random == nil {
		s.random = rand.Reader
	}
	return s, nil
}

func hashToken(token string) []byte {
	sum := sha256.Sum256([]byte(token))
	return sum[:]
}
