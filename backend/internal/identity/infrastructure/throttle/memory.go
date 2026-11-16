// Package throttle blocks Usernames after repeated failed sign-ins.
package throttle

import (
	"sync"
	"time"

	"github.com/elyosemite/openboard/backend/internal/identity/domain"
)

// pruneThreshold is how many tracked Usernames trigger a cleanup of stale ones.
const pruneThreshold = 10_000

type entry struct {
	failures     int
	lastFailure  time.Time
	blockedUntil time.Time
}

// Memory is an in-process SignInThrottle. It suits a single API instance;
// several instances would need a shared store.
type Memory struct {
	mu          sync.Mutex
	maxFailures int
	lockout     time.Duration
	now         func() time.Time
	entries     map[domain.Username]*entry
}

// NewMemory blocks a Username for lockout after maxFailures consecutive
// failures. Failures older than lockout no longer count.
func NewMemory(maxFailures int, lockout time.Duration, now func() time.Time) *Memory {
	return &Memory{
		maxFailures: maxFailures,
		lockout:     lockout,
		now:         now,
		entries:     map[domain.Username]*entry{},
	}
}

// Allow reports whether username may attempt to sign in now.
func (m *Memory) Allow(username domain.Username) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	e, ok := m.entries[username]
	return !ok || !m.now().Before(e.blockedUntil)
}

// RecordFailure counts a failed sign-in and blocks the Username when needed.
func (m *Memory) RecordFailure(username domain.Username) {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := m.now()
	e, ok := m.entries[username]
	if !ok || now.Sub(e.lastFailure) >= m.lockout {
		e = &entry{}
		m.entries[username] = e
	}
	e.failures++
	e.lastFailure = now
	if e.failures >= m.maxFailures {
		e.blockedUntil = now.Add(m.lockout)
		e.failures = 0
	}
	m.prune(now)
}

// Reset forgets the failures of username after a successful sign-in.
func (m *Memory) Reset(username domain.Username) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.entries, username)
}

func (m *Memory) prune(now time.Time) {
	if len(m.entries) <= pruneThreshold {
		return
	}
	for username, e := range m.entries {
		if now.Sub(e.lastFailure) >= m.lockout && !now.Before(e.blockedUntil) {
			delete(m.entries, username)
		}
	}
}
