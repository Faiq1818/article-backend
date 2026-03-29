package pkg

import "strings"

func SlugGenerate(baseText string) (string, error) {
	hash, err := RandomHash()
	if err != nil {
		return "", err
	}

	// generate slug and title
	slug := strings.ReplaceAll(baseText, " ", "-")
	slug = strings.ToLower(slug)
	cutHash := hash[:5]
	slugGenerate := slug + "-" + cutHash

	return slugGenerate, nil
}
