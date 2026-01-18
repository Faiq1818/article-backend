package article

import (
	"database/sql"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	DB       *sql.DB
	Validate *validator.Validate
}
