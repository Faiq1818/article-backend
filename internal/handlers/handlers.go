package handler

import (
	"article/internal/services/articles"
	"article/internal/services/auths"

	"database/sql"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func SetupRoutes(db *sql.DB, validate *validator.Validate) *http.ServeMux {
	// Dependency Injection
	authHandler := &auth.Handler{
		DB:       db,
		Validate: validate,
	}
	articleHandler := &article.AuthHandler{
		DB:       db,
		Validate: validate,
	}

	// initiate route
	router := http.NewServeMux()

	// routes
	AuthsHandler(router, authHandler)
	ArticlesHandler(router, articleHandler)

	return router
}
