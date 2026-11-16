package domain

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Password length limits (NIST SP 800-63B): length is the only strength rule.
const (
	MinPasswordLength = 12
	MaxPasswordLength = 128
)

// ValidateNewPassword checks a password a User is choosing. It must respect the
// length limits and differ from the current password and from the Username.
// Composition rules (digits, symbols, case) are deliberately not enforced.
func ValidateNewPassword(newPassword, currentPassword string, username Username) error {
	length := utf8.RuneCountInString(newPassword)
	switch {
	case length < MinPasswordLength:
		return fmt.Errorf("%w: use at least %d characters", ErrWeakPassword, MinPasswordLength)
	case length > MaxPasswordLength:
		return fmt.Errorf("%w: use at most %d characters", ErrWeakPassword, MaxPasswordLength)
	case newPassword == currentPassword:
		return fmt.Errorf("%w: it must differ from the current password", ErrWeakPassword)
	case strings.EqualFold(newPassword, username.String()):
		return fmt.Errorf("%w: it must differ from the username", ErrWeakPassword)
	}
	return nil
}
