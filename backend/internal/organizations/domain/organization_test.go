package domain

import (
	"errors"
	"strings"
	"testing"
)

func TestParseName(t *testing.T) {
	valid := map[string]string{
		"Acme":                   "Acme",
		"  Padaria Pão Quente  ": "Padaria Pão Quente",
		strings.Repeat("ã", 255): strings.Repeat("ã", 255),
	}
	for raw, want := range valid {
		if got, err := ParseName(raw); err != nil || got != want {
			t.Errorf("ParseName(%q) = %q, %v; want %q", raw, got, err, want)
		}
	}
	for _, raw := range []string{"", "   ", strings.Repeat("a", 256)} {
		if _, err := ParseName(raw); !errors.Is(err, ErrInvalidName) {
			t.Errorf("ParseName(%q) err = %v, want ErrInvalidName", raw, err)
		}
	}
}

func TestParseUserID(t *testing.T) {
	got, err := ParseUserID(" 0B5E3C1A-9F2D-4E6B-8A7C-1D2E3F4A5B6C ")
	if err != nil || got != "0b5e3c1a-9f2d-4e6b-8a7c-1d2e3f4a5b6c" {
		t.Errorf("ParseUserID = %q, %v; want lowercase UUID", got, err)
	}
	for _, raw := range []string{"", "u1", "0b5e3c1a-9f2d-4e6b-8a7c-1d2e3f4a5b6", "0b5e3c1a9f2d4e6b8a7c1d2e3f4a5b6c"} {
		if _, err := ParseUserID(raw); !errors.Is(err, ErrInvalidUserID) {
			t.Errorf("ParseUserID(%q) err = %v, want ErrInvalidUserID", raw, err)
		}
	}
}

func TestOnlyAdministratorsCreateOrganizations(t *testing.T) {
	if !(Actor{IsAdministrator: true}).CanCreateOrganization() {
		t.Error("an Administrator must be able to create Organizations")
	}
	if (Actor{}).CanCreateOrganization() {
		t.Error("a User who is not an Administrator must not create Organizations")
	}
}
