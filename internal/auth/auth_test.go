package auth

import (
	"testing"

	"github.com/alexedwards/argon2id"
)

func TestHashPassword(t *testing.T) {
	password := "password123"
	want, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		t.Errorf("Hashing failed: %v", err)
	}
	hashed_password, err := HashPassword(password)
	match, err1 := argon2id.ComparePasswordAndHash(password, hashed_password)
	if err1 != nil {
		t.Errorf("Hashing failed: %v", err)
	}
	if !match || err != nil {
		t.Errorf(`HashPassword(%q) = %q, %v, want %q, <nil>`, password, hashed_password, err, want)
	}
}

func TestComparePasswordHash(t *testing.T) {
	password := "password123"
	hashed_password, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		t.Errorf("Hashing failed: %v", err)
	}
	match, err := CheckPasswordHash(password, hashed_password)
	if !match || err != nil {
		t.Errorf(`CheckPasswordHash("password123", %q) = %t, %v, want true, <nil>`, hashed_password, match, err)
	}
}
