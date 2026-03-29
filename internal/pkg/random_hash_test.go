package pkg

import (
	"testing"
)

func TestRandomHash_NoError(t *testing.T) {
	hash, err := RandomHash()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if hash == "" {
		t.Fatalf("expected non-empty hash")
	}
}

func TestRandomHash_Length(t *testing.T) {
	hash, err := RandomHash()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// SHA-256 = 32 bytes → hex encoding = 64 chars
	if len(hash) != 64 {
		t.Fatalf("expected hash length 64, got %d", len(hash))
	}
}

func TestRandomHash_Unique(t *testing.T) {
	hash1, err1 := RandomHash()
	hash2, err2 := RandomHash()

	if err1 != nil || err2 != nil {
		t.Fatalf("unexpected error: %v %v", err1, err2)
	}

	if hash1 == hash2 {
		t.Fatalf("expected different hashes, got same value")
	}
}
