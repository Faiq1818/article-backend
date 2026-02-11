package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	pkg "article/internal/pkg"
	requesttype "article/internal/request_type"
	article "article/internal/services/articles"

	"github.com/go-playground/validator/v10"
)

func (DI *Dependency_Injection) SaveArticle(inject *article.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// decode body
		var req requesttype.SaveArticleRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// validate body
		err = DI.Validate.Struct(req)
		if err != nil {
			errors := err.(validator.ValidationErrors)
			http.Error(w, fmt.Sprintf("Validation error: %s", errors), http.StatusBadRequest)
			return
		}

		// bussiness logic
		err = inject.SaveArticle(req)
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

func (DI *Dependency_Injection) GetArticle(inject *article.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get url params
		queryParams := r.URL.Query()

		// convert limit and page params to integer
		limit, err := strconv.Atoi(queryParams.Get("limit"))
		if err != nil {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}

		page, err := strconv.Atoi(queryParams.Get("page"))
		if err != nil || page < 1 {
			http.Error(w, "invalid page", http.StatusBadRequest)
			return
		}

		// bussiness logic
		articles, err := inject.GetArticle(page, limit)
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
			Message: "Artikel berhasil didapat",
			Success: true,
			Data:    articles,
		})
	}
}
