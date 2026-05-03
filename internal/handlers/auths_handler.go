package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"article/internal/apperror"
	"article/internal/httpx"
	"article/internal/middlewares"
	requesttype "article/internal/request_type"
	auths "article/internal/services/auths"

	"github.com/golang-jwt/jwt/v5"
)

func Register(inject *auths.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// decode body
		var req requesttype.RegisterRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			httpx.JSONResponse(w, http.StatusBadRequest, httpx.Response{
				Message: "Invalid request Body",
				Success: false,
			})
			return
		}

		// validate body
		err = inject.Validate.Struct(req)
		if err != nil {
			httpx.JSONResponse(w, http.StatusBadRequest, httpx.Response{
				Message: "Validation error",
				Success: false,
				Data:    httpx.FormatValidationError(err),
			})
			return
		}

		// bussiness logic
		err = inject.Register(req)
		if err != nil {
			var appErr *apperror.AppError
			if errors.As(err, &appErr) {
				httpx.JSONResponse(w, appErr.Code, httpx.Response{
					Message: appErr.Message,
					Success: false,
				})
				return
			}

			// fallback unknown error
			httpx.JSONResponse(w, http.StatusInternalServerError, httpx.Response{
				Message: "internal server error",
				Success: false,
			})
			return
		}

		// success response
		httpx.JSONResponse(w, http.StatusOK, httpx.Response{
			Message: "Account created successfully",
			Success: true,
		})
	}
}

func Login(inject *auths.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// decode body
		var req requesttype.LoginRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			httpx.JSONResponse(w, http.StatusBadRequest, httpx.Response{
				Message: "Invalid request Body",
				Success: false,
			})
			return
		}

		// validate body
		err = inject.Validate.Struct(req)
		if err != nil {
			httpx.JSONResponse(w, http.StatusBadRequest, httpx.Response{
				Message: "Validation error",
				Success: false,
				Data:    httpx.FormatValidationError(err),
			})
			return
		}

		// bussiness logic
		token, err := inject.Login(req)
		if err != nil {
			var appErr *apperror.AppError
			if errors.As(err, &appErr) {
				httpx.JSONResponse(w, appErr.Code, httpx.Response{
					Message: appErr.Message,
					Success: false,
				})
				return
			}

			// fallback unknown error
			httpx.JSONResponse(w, http.StatusInternalServerError, httpx.Response{
				Message: "internal server error",
				Success: false,
			})
			return
		}

		cookie := &http.Cookie{
			Name:     "session_token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
		}

		http.SetCookie(w, cookie)

		// success response
		httpx.JSONResponse(w, http.StatusOK, httpx.Response{
			Message: "Login success",
			Success: true,
		})
	}
}

func Me(inject *auths.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userInfo := r.Context().Value(middlewares.UserInfoKey).(jwt.MapClaims)

		userProfile, err := inject.CheckUserAuthorization(userInfo)
		if err != nil {
			var appErr *apperror.AppError
			if errors.As(err, &appErr) {
				httpx.JSONResponse(w, appErr.Code, httpx.Response{
					Message: appErr.Message,
					Success: false,
				})
				return
			}

			httpx.JSONResponse(w, http.StatusInternalServerError, httpx.Response{
				Message: "Internal server error",
				Success: false,
			})
			return
		}

		httpx.JSONResponse(w, http.StatusOK, httpx.Response{
			Message: "Successfully retrieved user profile",
			Success: true,
			Data:    userProfile,
		})
	}
}
