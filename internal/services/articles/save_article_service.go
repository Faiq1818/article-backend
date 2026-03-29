package article

import (
	"context"
	"log"
	"strings"

	pkg "article/internal/pkg"
	requesttype "article/internal/request_type"
)

func (s *Service) SaveArticle(ctx context.Context, req requesttype.SaveArticleRequest, ext string) error {
	// img s3 upload
	srcFile, err := req.Image.Open()
	if err != nil {
		return &pkg.AppError{Message: "Failed to read image file", Code: 400, Err: err}
	}
	defer srcFile.Close() // prevent memori leak in service layer

	// generate image name
	hash, err := pkg.RandomHash()
	if err != nil {
		s.Logger.Warn("Hashing error")
		return &pkg.AppError{Message: "Failed to generate image key", Code: 500, Err: err}
	}
	objectKey := "articles/" + hash + ext

	// upload image
	imageUrl, err := s.S3Repo.UploadObject(ctx, objectKey, srcFile)
	if err != nil {
		s.Logger.Error("S3 Upload Failed")
		return &pkg.AppError{Message: "Failed to upload image", Code: 500, Err: err}
	}

	// generate slug and title
	slug, err := pkg.SlugGenerate(req.Title)
	if err != nil {
		s.Logger.Error("Slug Generate Error", "err", err)
		return &pkg.AppError{Message: "Failed to generate slug", Code: 500, Err: err}
	}

	err = s.Repo.SaveArticle(ctx, req, imageUrl, slug)
	if err != nil {
		statusCode, clientMessage := pkg.ParsePostgresError(err)
		s.Logger.Error("Error inserting article", "err", err)

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
	hash, err := pkg.RandomHash()
	if err != nil {
		s.Logger.Warn("Hashing error")
	}

	var imageUrl string
	if ext != nil {
		srcFile, err := req.Image.Open()
		if err != nil {
			return &pkg.AppError{Message: "Failed to read image file", Code: 400, Err: err}
		}
		defer srcFile.Close() // prevent memori leak in service layer

		objectKey := "articles/" + hash + *ext

		// upload image
		imageUrl, err = s.S3Repo.UploadObject(ctx, objectKey, srcFile)
		if err != nil {
			s.Logger.Error("S3 Upload Failed")
			return &pkg.AppError{Message: "Failed to upload image", Code: 500, Err: err}
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
