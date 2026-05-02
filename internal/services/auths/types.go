package auths

import (
	"log/slog"

	appconfig "article/internal/config"
	"article/internal/repositories"

	"github.com/go-playground/validator/v10"
)

type Service struct {
	Repo     repositories.AuthRepository
	Validate *validator.Validate
	Logger   *slog.Logger
	Config   *appconfig.Config
}
