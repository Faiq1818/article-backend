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
	authHandler := &auth.AuthHandler{
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
	// routes/auth
	router.HandleFunc("POST /auth/register", authHandler.Register)
	router.HandleFunc("POST /auth/login", authHandler.Login)

	// routes/article
	router.HandleFunc("POST /article", articleHandler.SaveArticle)

	return router
}
