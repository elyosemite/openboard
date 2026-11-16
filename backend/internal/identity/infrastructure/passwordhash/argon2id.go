// Package passwordhash implements password hashing with argon2id.
package passwordhash

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// maxKeyLength bounds the derived key read from a stored hash.
const maxKeyLength = 1024

// ErrMalformedHash reports a stored hash that is not a valid argon2id string.
var ErrMalformedHash = errors.New("malformed argon2id hash")

// Params are the argon2id cost parameters.
type Params struct {
	MemoryKiB   uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

// OWASPParams follows the OWASP Password Storage Cheat Sheet minimum for argon2id.
var OWASPParams = Params{MemoryKiB: 19 * 1024, Iterations: 2, Parallelism: 1, SaltLength: 16, KeyLength: 32}

// Argon2id hashes passwords into the standard PHC string format, which keeps
// the parameters next to the hash so they can be raised later.
type Argon2id struct {
	params Params
}

// New returns a hasher using the given parameters.
func New(params Params) Argon2id {
	return Argon2id{params: params}
}

// Hash returns $argon2id$v=19$m=...,t=...,p=...$salt$key with a random salt.
func (a Argon2id) Hash(password string) (string, error) {
	salt := make([]byte, a.params.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generate salt: %w", err)
	}
	p := a.params
	key := argon2.IDKey([]byte(password), salt, p.Iterations, p.MemoryKiB, p.Parallelism, p.KeyLength)
	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, p.MemoryKiB, p.Iterations, p.Parallelism,
		base64.RawStdEncoding.EncodeToString(salt), base64.RawStdEncoding.EncodeToString(key)), nil
}

// Verify reports whether password matches encodedHash, in constant time.
func (a Argon2id) Verify(password, encodedHash string) (bool, error) {
	p, salt, key, err := decode(encodedHash)
	if err != nil {
		return false, err
	}
	candidate := argon2.IDKey([]byte(password), salt, p.Iterations, p.MemoryKiB, p.Parallelism, p.KeyLength)
	return subtle.ConstantTimeCompare(key, candidate) == 1, nil
}

func decode(encoded string) (Params, []byte, []byte, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return Params{}, nil, nil, ErrMalformedHash
	}

	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil || version != argon2.Version {
		return Params{}, nil, nil, ErrMalformedHash
	}

	var p Params
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &p.MemoryKiB, &p.Iterations, &p.Parallelism); err != nil {
		return Params{}, nil, nil, ErrMalformedHash
	}

	salt, saltErr := base64.RawStdEncoding.DecodeString(parts[4])
	key, keyErr := base64.RawStdEncoding.DecodeString(parts[5])
	if saltErr != nil || keyErr != nil {
		return Params{}, nil, nil, ErrMalformedHash
	}
	keyLength := len(key)
	if keyLength <= 0 || keyLength > maxKeyLength {
		return Params{}, nil, nil, ErrMalformedHash
	}
	p.KeyLength = uint32(keyLength)
	return p, salt, key, nil
}
