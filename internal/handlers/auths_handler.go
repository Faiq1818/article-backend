package handler

import (
	"net/http"

	"article/internal/services/auths"
)

func AuthsHandler(router *http.ServeMux, authHandler *auth.Handler) {
	router.HandleFunc("POST /auth/register", authHandler.Register)
	router.HandleFunc("POST /auth/login", authHandler.Login)
}
