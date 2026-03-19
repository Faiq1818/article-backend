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
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomHash() (string, error) {
	randomBytes := make([]byte, 32)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
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

	err = s.Repo.SaveArticle(req, imageUrl, slugGenerate)
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

func (s *Service) PutArticle(ctx context.Context, req requesttype.PutArticleRequest, ext *string, oldSlug string) error {
	// img s3 upload

	// generate image name
	hash, err := randomHash()
	if err != nil {
		s.Logger.Warn("hashing error")
	}

	var imageUrl string
	if ext != nil {
		srcFile, err := req.Image.Open()
		if err != nil {
			return &pkg.AppError{Message: "Gagal membaca file gambar", Code: 400, Err: err}
		}
		defer srcFile.Close() // prevent memori leak in service layer

		objectKey := "articles/" + hash + *ext

		// upload image
		imageUrl, err = s.S3Repo.UploadObject(ctx, objectKey, srcFile)
		if err != nil {
			s.Logger.Error("S3 Upload Failed")
			return &pkg.AppError{Message: "Gagal mengupload gambar", Code: 500, Err: err}
		}

	}

	// generate slug and title
	slug := strings.ReplaceAll(req.Title, " ", "-")
	slug = strings.ToLower(slug)
	cutHash := hash[:5]
	slugGenerate := slug + "-" + cutHash

	err = s.Repo.PutArticle(req, imageUrl, slugGenerate, oldSlug)
	if err != nil {
		statusCode, clientMessage := pkg.ParsePostgresError(err)
		log.Printf("Error updating article: %v", err)

		return &pkg.AppError{
			Message: clientMessage,
			Code:    statusCode,
			Err:     err,
		}
	}

	s.Logger.Info("Successfully updating article!")
	return nil
}
