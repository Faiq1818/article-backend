package middlewares

import (
	"article/internal/pkg"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func ValidateToken(tokenString string) (*jwt.Token, error) {
	key := []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}

func GetClaims(token *jwt.Token) (jwt.MapClaims, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("failed to cast claims")
	}
	return claims, nil
}

func AuthMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookieToken, err := r.Cookie("session_token")
			if err != nil {
				logger.Warn("Unauthorized access attempt: missing cookie",
					"method", r.Method,
					"path", r.URL.Path,
					"remote_addr", r.RemoteAddr,
					"error", err,
				)

				if err == http.ErrNoCookie {
					// Cookie not found
					pkg.JSONResponse(w, http.StatusUnauthorized, pkg.Response{
						Message: "No session cookie found",
						Success: false,
					})
					return
				}

				// Other error
				pkg.JSONResponse(w, http.StatusBadRequest, pkg.Response{
					Message: "Error reading cookie",
					Success: false,
				})
				return
			}

			token, err := ValidateToken(cookieToken.Value)
			if err != nil {
				logger.Warn("Unauthorized access attempt: invalid token",
					"method", r.Method,
					"path", r.URL.Path,
					"remote_addr", r.RemoteAddr,
					"error", err,
				)

				pkg.JSONResponse(w, http.StatusUnauthorized, pkg.Response{
					Message: "Unauthorized: " + err.Error(),
					Success: false,
				})
				return
			}

			claims, _ := GetClaims(token)
			ctx := context.WithValue(r.Context(), "user_info", claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
