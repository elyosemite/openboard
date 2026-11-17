package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestAdministratorCreatesOrganization(t *testing.T) {
	c := newAPIClient(t)
	c.expect(http.MethodPost, "/api/v1/organizations", `{"name":"Acme"}`, http.StatusUnauthorized)

	c.expect(http.MethodPost, "/api/v1/sessions", `{"username":"admin","password":"admin"}`, http.StatusCreated)
	c.expect(http.MethodPost, "/api/v1/organizations", `{"name":"Acme"}`, http.StatusForbidden)
	c.expect(http.MethodPut, "/api/v1/users/me/password",
		`{"current_password":"admin","new_password":"a long and fine passphrase"}`, http.StatusNoContent)

	c.expect(http.MethodPost, "/api/v1/organizations",
		`{"name":"Acme","owner_user_id":"00000000-0000-0000-0000-000000000000"}`, http.StatusUnprocessableEntity)
	rec := c.expect(http.MethodPost, "/api/v1/organizations", `{"name":"  Acme  "}`, http.StatusCreated)

	var org struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&org); err != nil {
		t.Fatal(err)
	}
	if org.Name != "Acme" || !strings.HasSuffix(rec.Header().Get("Location"), "/"+org.ID) {
		t.Errorf("org %+v, Location %q", org, rec.Header().Get("Location"))
	}
}
