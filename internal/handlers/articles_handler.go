package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	pkg "article/internal/pkg"
	requesttype "article/internal/request_type"
	article "article/internal/services/articles"
)

func SaveArticle(inject *article.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// decode body
		var req requesttype.SaveArticleRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			pkg.JSONResponse(w, http.StatusBadRequest, pkg.Response{
				Message: "Invalid Body",
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

func GetArticle(inject *article.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get url params
		queryParams := r.URL.Query()

		// convert limit and page params to integer
		limit, err := strconv.Atoi(queryParams.Get("limit"))
		if err != nil {
			pkg.JSONResponse(w, http.StatusBadRequest, pkg.Response{
				Message: "Parameter limit tidak valid",
				Success: false,
			})
			return
		}

		page, err := strconv.Atoi(queryParams.Get("page"))
		if err != nil || page < 1 {
			pkg.JSONResponse(w, http.StatusBadRequest, pkg.Response{
				Message: "Parameter page tidak valid",
				Success: false,
			})
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
