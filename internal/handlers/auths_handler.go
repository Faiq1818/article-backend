package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	pkg "article/internal/pkg"
	requesttype "article/internal/request_type"
	auths "article/internal/services/auths"
)

func Register(inject *auths.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// decode body
		var req requesttype.RegisterRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			pkg.JSONResponse(w, http.StatusBadRequest, pkg.Response{
				Message: "Invalid request Body",
				Success: false,
			})
			return
		}

		// validate body
		err = inject.Validate.Struct(req)
		if err != nil {
			pkg.JSONResponse(w, http.StatusBadRequest, pkg.Response{
				Message: "Validation error",
				Success: false,
				Data:    pkg.FormatValidationError(err),
			})
			return
		}

		// bussiness logic
		err = inject.Register(req)
		if err != nil {
			var appErr *pkg.AppError
			if errors.As(err, &appErr) {
				pkg.JSONResponse(w, appErr.Code, pkg.Response{
					Message: appErr.Message,
					Success: false,
				})
				return
			}

			// fallback unknown error
			pkg.JSONResponse(w, http.StatusInternalServerError, pkg.Response{
				Message: "internal server error",
				Success: false,
			})
			return
		}

		// success response
		pkg.JSONResponse(w, http.StatusOK, pkg.Response{
			Message: "Akun berhasil dibuat",
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
			pkg.JSONResponse(w, http.StatusBadRequest, pkg.Response{
				Message: "Invalid request Body",
				Success: false,
			})
			return
		}

		// validate body
		err = inject.Validate.Struct(req)
		if err != nil {
			pkg.JSONResponse(w, http.StatusBadRequest, pkg.Response{
				Message: "Validation error",
				Success: false,
				Data:    pkg.FormatValidationError(err),
			})
			return
		}

		// bussiness logic
		token, err := inject.Login(req)
		if err != nil {
			var appErr *pkg.AppError
			if errors.As(err, &appErr) {
				pkg.JSONResponse(w, appErr.Code, pkg.Response{
					Message: appErr.Message,
					Success: false,
				})
				return
			}

			// fallback unknown error
			pkg.JSONResponse(w, http.StatusInternalServerError, pkg.Response{
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
		pkg.JSONResponse(w, http.StatusOK, pkg.Response{
			Message: "Login success",
			Success: true,
		})
	}
}
