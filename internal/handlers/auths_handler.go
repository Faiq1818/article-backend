package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	pkg "article/internal/pkg"
	requesttype "article/internal/request_type"
	auths "article/internal/services/auths"

	"github.com/go-playground/validator/v10"
)

func Register(inject *auths.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// decode body
		var req requesttype.RegisterRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// validate body
		err = inject.Validate.Struct(req)
		if err != nil {
			errors := err.(validator.ValidationErrors)
			http.Error(w, fmt.Sprintf("Validation error: %s", errors), http.StatusBadRequest)
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

func Login(inject *auths.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// decode body
		var req requesttype.LoginRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// validate body
		err = inject.Validate.Struct(req)
		if err != nil {
			errors := err.(validator.ValidationErrors)
			http.Error(w, fmt.Sprintf("Validation error: %s", errors), http.StatusBadRequest)
			return
		}

		// bussiness logic
		err = inject.Login(req)
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
