package handlers

import (
	"article/internal/services/articles"
	"article/internal/services/auths"

	"database/sql"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Dependency_Injection struct {
	DB       *sql.DB
	Validate *validator.Validate
}

func (DI *Dependency_Injection) SetupRoutes() *http.ServeMux {
	// Dependency Injection
	authInject := &auths.Handler{
		DB:       DI.DB,
		Validate: DI.Validate,
	}
	articleInject := &article.Handler{
		DB:       DI.DB,
		Validate: DI.Validate,
	}

	// initiate route
	router := http.NewServeMux()

	// routes
	DI.AuthsHandler(router, authInject)
	ArticlesHandler(router, articleInject)

	return router
}
