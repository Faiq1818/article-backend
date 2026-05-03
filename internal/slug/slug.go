package slug

import (
	"strings"

	"article/internal/hashutil"
)

func Generate(baseText string) (string, error) {
	hash, err := hashutil.RandomHash()
	if err != nil {
		return "", err
	}

	s := strings.ReplaceAll(baseText, " ", "-")
	s = strings.ToLower(s)
	cutHash := hash[:5]
	return s + "-" + cutHash, nil
}
