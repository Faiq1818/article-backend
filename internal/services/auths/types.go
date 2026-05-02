package auths

import (
	"log/slog"

	"article/internal/repositories"

	"github.com/go-playground/validator/v10"
)

type Service struct {
	Repo     repositories.AuthRepository
	Validate *validator.Validate
	Logger   *slog.Logger
}
