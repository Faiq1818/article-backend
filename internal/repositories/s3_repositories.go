package repositories

import (
	"context"
	"io"
)

type S3Repository interface {
	// GetUserByEmail(email string) (*models.User, error)
	UploadObject(ctx context.Context, key string, fileBody io.Reader) (string, error)
}
