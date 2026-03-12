package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"

	middlewares "article/internal/middlewares"
	postgres "article/internal/repositories/postgres"
	s3Repo "article/internal/repositories/s3"
	article "article/internal/services/articles"
	auths "article/internal/services/auths"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-playground/validator/v10"
)

func SetupRoutes(db *sql.DB, validate *validator.Validate, s3Client *s3.Client, s3Uploader *manager.Uploader, logger *slog.Logger) *http.ServeMux {
	// Repository db initiate
	authRepo := postgres.NewAuthRepository(db)
	articleRepo := postgres.NewArticleRepository(db)
	s3Repo := s3Repo.NewS3Repository(s3Client, s3Uploader)

	// Dependency Injection
	authInject := &auths.Service{
		Repo:     authRepo,
		Validate: validate,
		S3Client: s3Client,
		Logger:   logger,
	}
	articleInject := &article.Service{
		Repo:     articleRepo,
		DB:       db,
		Validate: validate,
		S3Repo:   s3Repo,
		Logger:   logger,
	}

	authMiddleware := middlewares.AuthMiddleware(logger)
	// initiate route
	router := http.NewServeMux()

	// routes
	// Public
	router.HandleFunc("POST /auth/login", Login(authInject))
	router.HandleFunc("GET /articles", GetArticles(articleInject))
	router.HandleFunc("GET /article/{slug}", GetArticleSlug(articleInject))

	// Admin
	router.Handle("GET /auth/me", authMiddleware(Me(authInject)))
	router.Handle("POST /auth/register", authMiddleware(Register(authInject)))

	router.Handle("GET /admin/articles", authMiddleware(AdminGetArticles(articleInject)))
	router.Handle("GET /admin/article/{slug}", authMiddleware(GetArticleSlug(articleInject)))
	router.Handle("POST /admin/article", authMiddleware(AdminSaveArticle(articleInject)))
	router.Handle("PUT /admin/article/{slug}", authMiddleware(AdminPutArticleSlug(articleInject)))

	return router
}
