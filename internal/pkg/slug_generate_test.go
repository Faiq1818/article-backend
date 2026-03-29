package pkg

import (
	"strings"
	"testing"
)

func TestSlugGenerate_NoError(t *testing.T) {
	slug, err := SlugGenerate("Hello World")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if slug == "" {
		t.Fatalf("expected non-empty slug")
	}
}

func TestSlugGenerate_Format(t *testing.T) {
	input := "Hello World Golang"
	slug, err := SlugGenerate(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Expected base slug
	expectedBase := "hello-world-golang"

	// cek prefix slug
	if !strings.HasPrefix(slug, expectedBase+"-") {
		t.Fatalf("expected slug to start with %s-, got %s", expectedBase, slug)
	}

	// split untuk cek hash
	parts := strings.Split(slug, "-")
	lastPart := parts[len(parts)-1]

	// hash dipotong 5 karakter
	if len(lastPart) != 5 {
		t.Fatalf("expected hash length 5, got %d", len(lastPart))
	}
}

func TestSlugGenerate_Lowercase(t *testing.T) {
	slug, err := SlugGenerate("HeLLo WoRLD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if slug != strings.ToLower(slug) {
		t.Fatalf("expected slug to be lowercase, got %s", slug)
	}
}

func TestSlugGenerate_Unique(t *testing.T) {
	slug1, err1 := SlugGenerate("same title")
	slug2, err2 := SlugGenerate("same title")

	if err1 != nil || err2 != nil {
		t.Fatalf("unexpected error: %v %v", err1, err2)
	}

	if slug1 == slug2 {
		t.Fatalf("expected unique slugs, got identical values")
	}
}
