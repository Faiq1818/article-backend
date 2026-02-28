package article

import (
	"article/internal/repositories"
	"database/sql"
	"log/slog"

	"github.com/go-playground/validator/v10"
)

type Service struct {
	Repo     repositories.ArticleRepository
	S3Repo   repositories.S3Repository
	DB       *sql.DB
	Validate *validator.Validate
	Logger   *slog.Logger
}
