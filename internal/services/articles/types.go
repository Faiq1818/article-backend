package article

import (
	"database/sql"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-playground/validator/v10"
)

type Service struct {
	DB         *sql.DB
	Validate   *validator.Validate
	S3Client   *s3.Client
	S3Uploader *manager.Uploader
	Logger     *slog.Logger
}
