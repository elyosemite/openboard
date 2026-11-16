package domain

import (
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

const maxNameLength = 255

var usernamePattern = regexp.MustCompile(`^[a-z0-9._-]{3,32}$`)

// Username is the unique identifier of a User within an Instance, always lowercase.
type Username string

// ParseUsername normalizes raw input to lowercase and validates its format.
func ParseUsername(raw string) (Username, error) {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	if !usernamePattern.MatchString(normalized) {
		return "", ErrInvalidUsername
	}
	return Username(normalized), nil
}

func (u Username) String() string { return string(u) }

// ParseName trims a display name and checks it is neither blank nor too long.
func ParseName(raw string) (string, error) {
	name := strings.TrimSpace(raw)
	if name == "" || utf8.RuneCountInString(name) > maxNameLength {
		return "", ErrInvalidName
	}
	return name, nil
}

// UserID identifies a User.
type UserID string

// User is a person's identity and credentials, global to the Instance.
type User struct {
	ID                 UserID
	Username           Username
	Name               string
	PasswordHash       string
	MustChangePassword bool
	IsAdministrator    bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
