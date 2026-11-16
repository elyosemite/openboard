package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/elyosemite/openboard/backend/internal/identity/domain"
)

const (
	defaultAdministratorUsername = "admin"
	defaultAdministratorPassword = "admin"
	administratorName            = "Administrator"
)

// ErrIncompleteAdministratorSeed reports that only one of the seed fields was set.
var ErrIncompleteAdministratorSeed = errors.New("set both the administrator username and password, or neither")

// AdministratorSeed is the installer's optional choice for the first
// Administrator. When both fields are empty, admin/admin is used as a
// Temporary Password.
type AdministratorSeed struct {
	Username string
	Password string
}

// EnsureAdministrator creates the first Administrator of an Instance that has
// no Users yet. It reports whether a User was created.
func (s *Service) EnsureAdministrator(ctx context.Context, seed AdministratorSeed) (bool, error) {
	count, err := s.users.Count(ctx)
	if err != nil {
		return false, fmt.Errorf("count users: %w", err)
	}
	if count > 0 {
		return false, nil
	}

	username, password, temporary, err := resolveSeed(seed)
	if err != nil {
		return false, err
	}
	hash, err := s.hasher.Hash(password)
	if err != nil {
		return false, fmt.Errorf("hash administrator password: %w", err)
	}

	_, err = s.users.Create(ctx, NewUser{
		Username:           username,
		Name:               administratorName,
		PasswordHash:       hash,
		MustChangePassword: temporary,
		IsAdministrator:    true,
	})
	if err != nil {
		return false, fmt.Errorf("create administrator: %w", err)
	}
	return true, nil
}

// resolveSeed returns the Username and password to use and whether the
// password is a Temporary Password. A configured password must follow the
// password rules; the default one is temporary and exempt.
func resolveSeed(seed AdministratorSeed) (domain.Username, string, bool, error) {
	if seed.Username == "" && seed.Password == "" {
		return defaultAdministratorUsername, defaultAdministratorPassword, true, nil
	}
	if seed.Username == "" || seed.Password == "" {
		return "", "", false, ErrIncompleteAdministratorSeed
	}

	username, err := domain.ParseUsername(seed.Username)
	if err != nil {
		return "", "", false, fmt.Errorf("administrator username: %w", err)
	}
	if err := domain.ValidateNewPassword(seed.Password, "", username); err != nil {
		return "", "", false, fmt.Errorf("administrator password: %w", err)
	}
	return username, seed.Password, false, nil
}
