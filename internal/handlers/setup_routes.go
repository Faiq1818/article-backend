package handlers

import (
	"database/sql"
	"net/http"

	article "article/internal/services/articles"
	auths "article/internal/services/auths"

	"github.com/go-playground/validator/v10"
)

func SetupRoutes(db *sql.DB, validate *validator.Validate) *http.ServeMux {
	// Dependency Injection
	authInject := &auths.Handler{
		DB:       db,
		Validate: validate,
	}
	articleInject := &article.Handler{
		DB:       db,
		Validate: validate,
	}

	// initiate route
	router := http.NewServeMux()

	// routes
	router.HandleFunc("POST /auth/register", Register(authInject))
	router.HandleFunc("POST /auth/login", Login(authInject))
	router.HandleFunc("POST /article", SaveArticle(articleInject))
	router.HandleFunc("GET /article", GetArticle(articleInject))
	router.HandleFunc("GET /article/{slug}", GetArticleSlug(articleInject))
	return router
}
