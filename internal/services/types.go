package services

import (
	"database/sql"
)

type AuthHandler struct {
	DB *sql.DB
}

type ErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
