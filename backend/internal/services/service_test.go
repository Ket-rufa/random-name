package services

import (
	"testing"
)

func TestHashToken(t *testing.T) {
	token := "my-secret-token"
	hashed := HashToken(token)

	if hashed == "" {
		t.Error("Expected hashed token to not be empty")
	}

	if hashed == token {
		t.Error("Expected hashed token to be different from original token")
	}

	// Ensure hash is deterministic
	if HashToken(token) != hashed {
		t.Error("Expected hashing the same token twice to yield identical hashes")
	}
}

func TestGenerateRandomString(t *testing.T) {
	length := 16
	str, err := generateRandomString(length)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(str) != length {
		t.Errorf("Expected string of length %d, got: %d", length, len(str))
	}
}
