package handlers

import (
	"database/sql"
	"net/http"

	article "article/internal/services/articles"
	auths "article/internal/services/auths"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-playground/validator/v10"
)

func SetupRoutes(db *sql.DB, validate *validator.Validate, s3Client *s3.Client) *http.ServeMux {
	// Dependency Injection
	authInject := &auths.Handler{
		DB:       db,
		Validate: validate,
		S3Client: s3Client,
	}
	articleInject := &article.Handler{
		DB:       db,
		Validate: validate,
		S3Client: s3Client,
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
