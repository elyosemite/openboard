package throttle

import (
	"fmt"
	"testing"
	"time"

	"github.com/elyosemite/openboard/backend/internal/identity/domain"
)

type clock struct{ now time.Time }

func (c *clock) Now() time.Time { return c.now }

func newThrottle() (*Memory, *clock) {
	c := &clock{now: time.Date(2026, 11, 15, 12, 0, 0, 0, time.UTC)}
	return NewMemory(5, 5*time.Minute, c.Now), c
}

func TestBlocksAfterMaxConsecutiveFailures(t *testing.T) {
	m, c := newThrottle()
	for i := 0; i < 4; i++ {
		m.RecordFailure("gabriel")
	}
	if !m.Allow("gabriel") {
		t.Fatal("four failures must not block")
	}
	m.RecordFailure("gabriel")
	if m.Allow("gabriel") {
		t.Fatal("the fifth failure must block")
	}
	if !m.Allow("helena") {
		t.Error("other Usernames must not be affected")
	}

	c.now = c.now.Add(5*time.Minute - time.Second)
	if m.Allow("gabriel") {
		t.Error("still blocked just before the lockout ends")
	}
	c.now = c.now.Add(time.Second)
	if !m.Allow("gabriel") {
		t.Error("allowed again when the lockout ends")
	}
}

func TestResetClearsFailures(t *testing.T) {
	m, _ := newThrottle()
	for i := 0; i < 4; i++ {
		m.RecordFailure("gabriel")
	}
	m.Reset("gabriel")
	m.RecordFailure("gabriel")
	if !m.Allow("gabriel") {
		t.Error("failures before a successful sign-in must not count")
	}
}

func TestOldFailuresStopCounting(t *testing.T) {
	m, c := newThrottle()
	for i := 0; i < 4; i++ {
		m.RecordFailure("gabriel")
	}
	c.now = c.now.Add(5 * time.Minute)
	m.RecordFailure("gabriel")
	if !m.Allow("gabriel") {
		t.Error("failures older than the lockout window must not count")
	}
}

func TestPruneDropsStaleEntries(t *testing.T) {
	m, c := newThrottle()
	for i := 0; i <= pruneThreshold; i++ {
		m.RecordFailure(domain.Username(fmt.Sprintf("user%d", i)))
	}
	for i := 0; i < 5; i++ {
		m.RecordFailure("attacked")
	}
	c.now = c.now.Add(time.Minute)
	m.RecordFailure("fresh")
	if len(m.entries) <= pruneThreshold {
		t.Fatalf("recent entries must be kept, got %d", len(m.entries))
	}

	c.now = c.now.Add(10 * time.Minute)
	m.RecordFailure("trigger")
	if len(m.entries) != 1 {
		t.Errorf("stale entries should be pruned, %d left", len(m.entries))
	}
}
