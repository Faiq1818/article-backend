package article

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"strings"

	pkg "article/internal/pkg"
	requesttype "article/internal/request_type"

	"github.com/google/uuid"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomHash() (string, error) {
	randomBytes := make([]byte, 32)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	// Hitung SHA-256 dari data acak
	hash := sha256.Sum256(randomBytes)

	return string(hex.EncodeToString(hash[:])), nil
}

func (s *Service) SaveArticle(ctx context.Context, req requesttype.SaveArticleRequest, ext string) error {
	// img s3 upload
	srcFile, err := req.Image.Open()
	if err != nil {
		return &pkg.AppError{Message: "Gagal membaca file gambar", Code: 400, Err: err}
	}
	defer srcFile.Close() // prevent memori leak in service layer

	// generate image name
	hash, err := randomHash()
	if err != nil {
		s.Logger.Warn("hashing error")
	}
	objectKey := "articles/" + hash + ext

	// upload image
	imageUrl, errS3 := s.S3Repo.UploadObject(ctx, objectKey, srcFile)
	if errS3 != nil {
		s.Logger.Error("S3 Upload Failed")
		return &pkg.AppError{Message: "Gagal mengupload gambar", Code: 500, Err: errS3}
	}

	// generate slug and title
	slug := strings.ReplaceAll(req.Title, " ", "-")
	slug = strings.ToLower(slug)
	cutHash := hash[:5]
	slugGenerate := slug + "-" + cutHash

	// db push
	u := uuid.New()
	_, err = s.DB.Exec("INSERT INTO article (id, title, slug, description, content, image_url) VALUES ($1, $2, $3, $4, $5, $6);", u, req.Title, slugGenerate, req.Description, req.Content, imageUrl)
	if err != nil {
		statusCode, clientMessage := pkg.ParsePostgresError(err)
		log.Printf("Error inserting article: %v", err)

		return &pkg.AppError{
			Message: clientMessage,
			Code:    statusCode,
			Err:     err,
		}
	}

	s.Logger.Info("Successfully inserting article!")
	return nil
}
