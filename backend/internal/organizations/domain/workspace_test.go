package domain

import (
	"errors"
	"testing"
)

const sampleID = "0B5E3C1A-9F2D-4E6B-8A7C-1D2E3F4A5B6C"

func TestParseWorkspaceID(t *testing.T) {
	id, err := ParseWorkspaceID(" " + sampleID + " ")
	if err != nil || id != "0b5e3c1a-9f2d-4e6b-8a7c-1d2e3f4a5b6c" {
		t.Errorf("ParseWorkspaceID = %q, %v", id, err)
	}
	if _, err := ParseWorkspaceID("legal"); !errors.Is(err, ErrWorkspaceNotFound) {
		t.Errorf("err = %v, want ErrWorkspaceNotFound", err)
	}
}

func TestParseOrganizationRef(t *testing.T) {
	id, err := ParseOrganizationRef(sampleID)
	if err != nil || id != "0b5e3c1a-9f2d-4e6b-8a7c-1d2e3f4a5b6c" {
		t.Errorf("ParseOrganizationRef = %q, %v", id, err)
	}
	if _, err := ParseOrganizationRef(""); !errors.Is(err, ErrInvalidOrganizationRef) {
		t.Errorf("err = %v, want ErrInvalidOrganizationRef", err)
	}
}

func TestParseVisibility(t *testing.T) {
	cases := map[string]struct {
		want Visibility
		err  error
	}{
		"":           {VisibilityRestricted, nil},
		"restricted": {VisibilityRestricted, nil},
		"open":       {VisibilityOpen, nil},
		"private":    {"", ErrInvalidVisibility},
	}
	for raw, tc := range cases {
		got, err := ParseVisibility(raw)
		if got != tc.want || !errors.Is(err, tc.err) {
			t.Errorf("ParseVisibility(%q) = %q, %v; want %q, %v", raw, got, err, tc.want, tc.err)
		}
	}
}

func TestOnlyOrganizationOwnersCreateWorkspaces(t *testing.T) {
	if !RoleOwner.CanCreateWorkspace() {
		t.Error("an Organization Owner must be able to create Workspaces")
	}
	if RoleContributor.CanCreateWorkspace() {
		t.Error("an Organization Contributor must not create Workspaces")
	}
}
