package handler

import (
	"net/http"

	"article/internal/services/articles"
)

func ArticlesHandler(router *http.ServeMux, articleHandler *article.AuthHandler) {
	router.HandleFunc("POST /article", articleHandler.SaveArticle)
	router.HandleFunc("GET /article", articleHandler.GetArticle)
}
