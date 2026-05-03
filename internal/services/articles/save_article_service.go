package article

import (
	"context"
	"log"
	"strings"

	"article/internal/apperror"
	"article/internal/hashutil"
	"article/internal/repositories/postgres"
	requesttype "article/internal/request_type"
	"article/internal/slug"
)

func (s *Service) SaveArticle(ctx context.Context, req requesttype.SaveArticleRequest, ext string) error {
	// img s3 upload
	srcFile, err := req.Image.Open()
	if err != nil {
		return &apperror.AppError{Message: "Failed to read image file", Code: 400, Err: err}
	}
	defer func() { _ = srcFile.Close() }()

	// generate image name
	hash, err := hashutil.RandomHash()
	if err != nil {
		s.Logger.Warn("Hashing error")
		return &apperror.AppError{Message: "Failed to generate image key", Code: 500, Err: err}
	}
	objectKey := "articles/" + hash + ext

	// upload image
	imageUrl, err := s.S3Repo.UploadObject(ctx, objectKey, srcFile)
	if err != nil {
		s.Logger.Error("S3 Upload Failed")
		return &apperror.AppError{Message: "Failed to upload image", Code: 500, Err: err}
	}

	// generate slug and title
	slugStr, err := slug.Generate(req.Title)
	if err != nil {
		s.Logger.Error("Slug Generate Error", "err", err)
		return &apperror.AppError{Message: "Failed to generate slug", Code: 500, Err: err}
	}

	err = s.Repo.SaveArticle(ctx, req, imageUrl, slugStr)
	if err != nil {
		statusCode, clientMessage := postgres.ParseError(err)
		s.Logger.Error("Error inserting article", "err", err)

		return &apperror.AppError{
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
	hash, err := hashutil.RandomHash()
	if err != nil {
		s.Logger.Warn("Hashing error")
	}

	var imageUrl string
	if ext != nil {
		srcFile, err := req.Image.Open()
		if err != nil {
			return &apperror.AppError{Message: "Failed to read image file", Code: 400, Err: err}
		}
		defer func() { _ = srcFile.Close() }()

		objectKey := "articles/" + hash + *ext

		// upload image
		imageUrl, err = s.S3Repo.UploadObject(ctx, objectKey, srcFile)
		if err != nil {
			s.Logger.Error("S3 Upload Failed")
			return &apperror.AppError{Message: "Failed to upload image", Code: 500, Err: err}
		}

	}

	// generate slug and title
	slugBase := strings.ReplaceAll(req.Title, " ", "-")
	slugBase = strings.ToLower(slugBase)
	cutHash := hash[:5]
	slugGenerate := slugBase + "-" + cutHash

	err = s.Repo.PutArticle(req, imageUrl, slugGenerate, oldSlug)
	if err != nil {
		statusCode, clientMessage := postgres.ParseError(err)
		log.Printf("Error updating article: %v", err)

		return &apperror.AppError{
			Message: clientMessage,
			Code:    statusCode,
			Err:     err,
		}
	}

	s.Logger.Info("Successfully updating article!")
	return nil
}
