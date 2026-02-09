package handlers

import (
	"net/http"

	"article/internal/services/auths"
)

func AuthsHandler(router *http.ServeMux, authHandler *auths.Handler) {
	router.HandleFunc("POST /auth/register", authHandler.Register)
	router.HandleFunc("POST /auth/login", authHandler.Login)
}
