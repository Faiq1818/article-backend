package article

import (
	appconfig "article/internal/config"
	"article/internal/repositories"
	"log/slog"

	"github.com/go-playground/validator/v10"
)

type Service struct {
	Repo     repositories.ArticleRepository
	S3Repo   repositories.S3Repository
	Validate *validator.Validate
	Logger   *slog.Logger
	Config   *appconfig.Config
}
