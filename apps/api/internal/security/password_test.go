package security

import (
	"testing"
)

func TestArgon2PasswordHashing(t *testing.T) {
	password := "SecretPassword123!"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if hash == "" {
		t.Fatal("Expected non-empty hash string")
	}

	match, err := ComparePassword(password, hash)
	if err != nil {
		t.Fatalf("ComparePassword failed: %v", err)
	}
	if !match {
		t.Fatal("Expected password to match hash")
	}

	wrongMatch, err := ComparePassword("WrongPassword!", hash)
	if err != nil {
		t.Fatalf("ComparePassword failed on wrong password: %v", err)
	}
	if wrongMatch {
		t.Fatal("Expected wrong password NOT to match hash")
	}
}
