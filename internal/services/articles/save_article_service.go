package article

import (
	"context"
	"crypto/rand"
	"log"
	"os"
	"strings"

	pkg "article/internal/pkg"
	requesttype "article/internal/request_type"
	s3helpers "article/internal/s3_helpers"

	"github.com/google/uuid"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomHash(n int) (string, error) {
	bytes := make([]byte, n)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	for i := range bytes {
		bytes[i] = letters[int(bytes[i])%len(letters)]
	}

	return string(bytes), nil
}

func (h *Handler) SaveArticle(ctx context.Context, req requesttype.SaveArticleRequest) error {
	// generate a slug
	slug := strings.ReplaceAll(req.Title, " ", "-")
	slug = strings.ToLower(slug)

	// img s3 upload
	srcFile, err := req.Image.Open()
	if err != nil {
		return &pkg.AppError{Message: "Gagal membaca file gambar", Code: 400, Err: err}
	}
	defer srcFile.Close() // Mencegah memori leak di layer service

	s3Actor := s3helpers.S3Actions{
		S3Client:  h.S3Client,
		S3Manager: h.S3Uploader,
	}

	hash, err := randomHash(5)
	if err != nil {
		log.Printf("hash error")
	}

	objectKey := slug + hash + ".jpg" // Catatan: Idealnya ekstensi file diekstrak dinamis
	_, errS3 := s3Actor.UploadObject(ctx, os.Getenv("S3_BUCKET_NAME"), objectKey, srcFile)

	if errS3 != nil {
		log.Printf("S3 Upload Failed: %v", errS3)
		return &pkg.AppError{Message: "Gagal mengupload gambar", Code: 500, Err: errS3}
	}

	// slug make
	slugGenerate := slug + "ioi" + hash
	titleGenerate := req.Title + "-" + hash

	// db push
	u := uuid.New()
	_, err = h.DB.Exec("INSERT INTO article (id, title, slug, content) VALUES ($1, $2, $3, $4);", u, titleGenerate, slugGenerate, req.Content)
	if err != nil {
		statusCode, clientMessage := pkg.ParsePostgresError(err)
		log.Printf("Error inserting user: %v", err)

		return &pkg.AppError{
			Message: clientMessage,
			Code:    statusCode,
			Err:     err,
		}
	}

	return nil
}
