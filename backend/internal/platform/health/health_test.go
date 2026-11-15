package health

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elyosemite/openboard/backend/internal/platform/httpx"
)

type fakeChecker struct{ err error }

func (fakeChecker) Name() string                  { return "fake" }
func (f fakeChecker) Check(context.Context) error { return f.err }

func call(t *testing.T, h httpx.HandlerFunc) (*httptest.ResponseRecorder, error) {
	t.Helper()
	rec := httptest.NewRecorder()
	err := h(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	return rec, err
}

func decodeStatus(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()
	var s Status
	if err := json.NewDecoder(rec.Body).Decode(&s); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return s.Status
}

func TestLiveIgnoresDependencies(t *testing.T) {
	h := NewHandler(fakeChecker{err: errors.New("down")})
	rec, err := call(t, h.Live)
	if err != nil || rec.Code != http.StatusOK || decodeStatus(t, rec) != "ok" {
		t.Fatalf("Live = %d, %v", rec.Code, err)
	}
}

func TestReadyWhenAllCheckersPass(t *testing.T) {
	h := NewHandler(fakeChecker{}, fakeChecker{})
	rec, err := call(t, h.Ready)
	if err != nil || rec.Code != http.StatusOK || decodeStatus(t, rec) != "ready" {
		t.Fatalf("Ready = %d, %v", rec.Code, err)
	}
}

func TestReadyFailsWhenACheckerFails(t *testing.T) {
	cause := errors.New("connection refused")
	h := NewHandler(fakeChecker{}, fakeChecker{err: cause})
	_, err := call(t, h.Ready)

	var apiErr *httpx.Error
	if !errors.As(err, &apiErr) || apiErr.Status != http.StatusServiceUnavailable {
		t.Fatalf("err = %v, want 503 *httpx.Error", err)
	}
	if !errors.Is(err, cause) {
		t.Error("cause should be wrapped")
	}
}
