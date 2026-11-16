package passwordhash

import (
	"errors"
	"strings"
	"testing"
)

// fastParams keeps tests quick; production uses OWASPParams.
var fastParams = Params{MemoryKiB: 64, Iterations: 1, Parallelism: 1, SaltLength: 16, KeyLength: 32}

func TestHashAndVerify(t *testing.T) {
	h := New(fastParams)

	encoded, err := h.Hash("correct horse battery")
	if err != nil {
		t.Fatalf("Hash: %v", err)
	}
	if !strings.HasPrefix(encoded, "$argon2id$v=19$m=64,t=1,p=1$") {
		t.Errorf("unexpected format: %s", encoded)
	}

	if ok, err := h.Verify("correct horse battery", encoded); err != nil || !ok {
		t.Errorf("correct password: ok=%v err=%v", ok, err)
	}
	if ok, err := h.Verify("wrong", encoded); err != nil || ok {
		t.Errorf("wrong password: ok=%v err=%v", ok, err)
	}
}

func TestHashUsesRandomSalt(t *testing.T) {
	h := New(fastParams)
	a, _ := h.Hash("same")
	b, _ := h.Hash("same")
	if a == b {
		t.Error("two hashes of the same password must differ")
	}
}

func TestVerifyUsesParamsStoredInTheHash(t *testing.T) {
	old, _ := New(fastParams).Hash("pw-from-before")
	stronger := New(Params{MemoryKiB: 128, Iterations: 2, Parallelism: 1, SaltLength: 16, KeyLength: 32})
	if ok, err := stronger.Verify("pw-from-before", old); err != nil || !ok {
		t.Errorf("hashes made with older params must still verify: ok=%v err=%v", ok, err)
	}
}

func TestOWASPParams(t *testing.T) {
	if OWASPParams.MemoryKiB != 19456 || OWASPParams.Iterations != 2 || OWASPParams.Parallelism != 1 {
		t.Errorf("unexpected OWASP params: %+v", OWASPParams)
	}
}

func TestVerifyRejectsMalformedHashes(t *testing.T) {
	valid, _ := New(fastParams).Hash("pw")
	parts := strings.Split(valid, "$")
	replace := func(i int, v string) string {
		p := append([]string(nil), parts...)
		p[i] = v
		return strings.Join(p, "$")
	}
	cases := map[string]string{
		"empty":         "",
		"bcrypt":        "$2a$10$abcdefghijklmnopqrstuv",
		"wrong version": replace(2, "v=16"),
		"bad version":   replace(2, "version"),
		"bad params":    replace(3, "m=x"),
		"bad salt":      replace(4, "!!!"),
		"bad key":       replace(5, "!!!"),
		"empty key":     replace(5, ""),
		"oversized key": replace(5, strings.Repeat("A", 1400)),
	}
	for name, encoded := range cases {
		if _, err := New(fastParams).Verify("pw", encoded); !errors.Is(err, ErrMalformedHash) {
			t.Errorf("%s: err = %v, want ErrMalformedHash", name, err)
		}
	}
}
