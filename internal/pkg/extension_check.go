package pkg

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

func DetectImageExtension(file multipart.File) (string, error) {
	buf := make([]byte, 512)

	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return "", err
	}

	contentType := http.DetectContentType(buf[:n])

	if !strings.HasPrefix(contentType, "image/") {
		return "", errors.New("file bukan gambar")
	}

	var ext string

	switch contentType {
	case "image/jpeg":
		ext = ".jpg"
	case "image/png":
		ext = ".png"
	case "image/webp":
		ext = ".webp"
	case "image/gif":
		ext = ".gif"
	case "image/svg+xml":
		ext = ".svg"
	default:
		ext = ".bin"
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	return ext, nil
}
