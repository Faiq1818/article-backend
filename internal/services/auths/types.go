package auths

import (
	"database/sql"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-playground/validator/v10"
)

type Service struct {
	DB       *sql.DB
	Validate *validator.Validate
	S3Client *s3.Client
	Logger   *slog.Logger
}
