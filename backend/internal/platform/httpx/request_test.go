package httpx

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func decodeStatus(t *testing.T, contentType, body string) (map[string]string, int) {
	t.Helper()
	var got map[string]string
	h := Handle(discardLogger(), func(w http.ResponseWriter, r *http.Request) error {
		if err := DecodeJSON(w, r, &got); err != nil {
			return err
		}
		w.WriteHeader(http.StatusNoContent)
		return nil
	})
	req := httptest.NewRequest(http.MethodPost, "/things", strings.NewReader(body))
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return got, rec.Code
}

func TestDecodeJSONAcceptsJSONWithParameters(t *testing.T) {
	got, status := decodeStatus(t, "application/json; charset=utf-8", `{"name":"x"}`)
	if status != http.StatusNoContent || got["name"] != "x" {
		t.Fatalf("status %d, got %v", status, got)
	}
}

func TestDecodeJSONRejections(t *testing.T) {
	cases := map[string]struct {
		contentType, body string
		status            int
	}{
		"missing content type": {"", `{}`, http.StatusUnsupportedMediaType},
		"form":                 {"application/x-www-form-urlencoded", `a=b`, http.StatusUnsupportedMediaType},
		"text plain":           {"text/plain", `{}`, http.StatusUnsupportedMediaType},
		"invalid json":         {"application/json", `{`, http.StatusBadRequest},
		"wrong shape":          {"application/json", `[1]`, http.StatusBadRequest},
		"too large":            {"application/json", `{"a":"` + strings.Repeat("x", maxBodyBytes) + `"}`, http.StatusBadRequest},
	}
	for name, tc := range cases {
		if _, status := decodeStatus(t, tc.contentType, tc.body); status != tc.status {
			t.Errorf("%s: status %d, want %d", name, status, tc.status)
		}
	}
}

func TestValidationProblemListsFields(t *testing.T) {
	h := Handle(discardLogger(), func(http.ResponseWriter, *http.Request) error {
		return ErrValidation(FieldError{Field: "name", Code: "required", Message: "name is required"})
	})
	rec := serve(h, http.MethodPost, "/things")
	p := decodeProblem(t, rec)
	if p.Status != http.StatusUnprocessableEntity || len(p.Errors) != 1 || p.Errors[0].Field != "name" {
		t.Fatalf("unexpected problem: %+v", p)
	}
}

func TestCrossOriginRejectsUnsafeCrossSiteRequests(t *testing.T) {
	h := CrossOrigin(discardLogger(), http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	cases := []struct {
		method, site string
		status       int
	}{
		{http.MethodPost, "cross-site", http.StatusForbidden},
		{http.MethodPost, "same-origin", http.StatusNoContent},
		{http.MethodGet, "cross-site", http.StatusNoContent},
		{http.MethodPost, "", http.StatusNoContent},
	}
	for _, tc := range cases {
		req := httptest.NewRequest(tc.method, "/things", nil)
		if tc.site != "" {
			req.Header.Set("Sec-Fetch-Site", tc.site)
		}
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != tc.status {
			t.Errorf("%s with Sec-Fetch-Site %q = %d, want %d", tc.method, tc.site, rec.Code, tc.status)
		}
		if rec.Code == http.StatusForbidden {
			if p := decodeProblem(t, rec); p.Type != "https://openboard.dev/problems/cross-origin-request" {
				t.Errorf("problem type %q", p.Type)
			}
		}
	}
}
