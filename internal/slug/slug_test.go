package slug

import (
	"strings"
	"testing"
)

func TestGenerate_NoError(t *testing.T) {
	s, err := Generate("Hello World")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if s == "" {
		t.Fatalf("expected non-empty slug")
	}
}

func TestGenerate_Format(t *testing.T) {
	input := "Hello World Golang"
	s, err := Generate(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Expected base slug
	expectedBase := "hello-world-golang"

	// cek prefix slug
	if !strings.HasPrefix(s, expectedBase+"-") {
		t.Fatalf("expected slug to start with %s-, got %s", expectedBase, s)
	}

	// split untuk cek hash
	parts := strings.Split(s, "-")
	lastPart := parts[len(parts)-1]

	// hash dipotong 5 karakter
	if len(lastPart) != 5 {
		t.Fatalf("expected hash length 5, got %d", len(lastPart))
	}
}

func TestGenerate_Lowercase(t *testing.T) {
	s, err := Generate("HeLLo WoRLD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if s != strings.ToLower(s) {
		t.Fatalf("expected slug to be lowercase, got %s", s)
	}
}

func TestGenerate_Unique(t *testing.T) {
	s1, err1 := Generate("same title")
	s2, err2 := Generate("same title")

	if err1 != nil || err2 != nil {
		t.Fatalf("unexpected error: %v %v", err1, err2)
	}

	if s1 == s2 {
		t.Fatalf("expected unique slugs, got identical values")
	}
}
