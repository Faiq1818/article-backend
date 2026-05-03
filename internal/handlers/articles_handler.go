package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"article/internal/apperror"
	"article/internal/httpx"
	"article/internal/imageutil"
	requesttype "article/internal/request_type"
	article "article/internal/services/articles"
)

func AdminSaveArticle(inject *article.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// multipart
		err := r.ParseMultipartForm(inject.Config.MaxUploadSizeBytes)
		if err != nil {
			httpx.JSONResponse(w, http.StatusBadRequest, httpx.Response{
				Message: "Failed to parse multipart payload",
				Success: false,
			})
			return
		}

		file, header, err := r.FormFile("image")
		if err != nil {
			httpx.JSONResponse(w, http.StatusBadRequest, httpx.Response{
				Message: "Invalid limit parameter",
				Success: false,
			})
			return
		}
		defer func() { _ = file.Close() }()

		// decode body
		req := requesttype.SaveArticleRequest{
			Title:       r.FormValue("title"),
			Content:     r.FormValue("content"),
			Description: r.FormValue("description"),
			Image:       header,
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

		// make dynamic image name extension
		srcFile, err := req.Image.Open()
		if err != nil {
			httpx.JSONResponse(w, http.StatusBadRequest, httpx.Response{
				Message: "Failed to open file",
				Success: false,
			})
			return
		}
		defer func() { _ = srcFile.Close() }()

		// detect image extension
		ext, err := imageutil.DetectExtension(srcFile)
		if err != nil {
			httpx.JSONResponse(w, http.StatusBadRequest, httpx.Response{
				Message: "File must be an image (.jpg, .jpeg, .png, .webp, .gif)",
				Success: false,
			})
			return
		}

		// bussiness logic
		err = inject.SaveArticle(ctx, req, ext)
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
			Message: "Article created successfully",
			Success: true,
		})
	}
}

func GetArticles(inject *article.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get url params
		queryParams := r.URL.Query()

		// convert limit and page params to integer
		limit, err := strconv.Atoi(queryParams.Get("limit"))
		if err != nil {
			httpx.JSONResponse(w, http.StatusBadRequest, httpx.Response{
				Message: "Invalid limit parameter",
				Success: false,
			})
			return
		}

		page, err := strconv.Atoi(queryParams.Get("page"))
		if err != nil || page < 1 {
			httpx.JSONResponse(w, http.StatusBadRequest, httpx.Response{
				Message: "Invalid page parameter",
				Success: false,
			})
			return
		}

		// bussiness logic
		articles, meta, err := inject.GetArticles(page, limit)
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
			Message: "Article retrieved successfully",
			Success: true,
			Data: map[string]any{
				"articles": articles,
				"meta":     meta,
			},
		})
	}
}

func GetArticleSlug(inject *article.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get slug
		slug := r.PathValue("slug")
		if slug == "" {
			httpx.JSONResponse(w, http.StatusBadRequest, httpx.Response{
				Message: "Article slug is required",
				Success: false,
			})
			return
		}

		// bussiness logic
		articles, err := inject.GetArticleSlug(slug)
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
			Message: "Article retrieved successfully",
			Success: true,
			Data:    articles,
		})
	}
}

func AdminPutArticleSlug(inject *article.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// get slug
		oldSlug := r.PathValue("slug")
		if oldSlug == "" {
			httpx.JSONResponse(w, http.StatusBadRequest, httpx.Response{
				Message: "Article slug is required",
				Success: false,
			})
			return
		}

		// multipart
		err := r.ParseMultipartForm(inject.Config.MaxUploadSizeBytes)
		if err != nil {
			httpx.JSONResponse(w, http.StatusBadRequest, httpx.Response{
				Message: "Failed to parse multipart payload",
				Success: false,
			})
			return
		}

		// if there is no image, header is set to nil
		file, header, err := r.FormFile("image")
		if err != nil {
			if errors.Is(err, http.ErrMissingFile) {
				header = nil
			} else {
				httpx.JSONResponse(w, http.StatusBadRequest, httpx.Response{
					Message: "Image file is invalid",
					Success: false,
				})
				return
			}
		} else {
			defer func() { _ = file.Close() }()
		}

		// decode body
		req := requesttype.PutArticleRequest{
			Title:       r.FormValue("title"),
			Content:     r.FormValue("content"),
			Description: r.FormValue("description"),
			Image:       header,
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

		var ext *string
		if req.Image != nil {
			// make dynamic image name extension
			srcFile, err := req.Image.Open()
			if err != nil {
				httpx.JSONResponse(w, http.StatusBadRequest, httpx.Response{
					Message: "Failed to open file",
					Success: false,
				})
				return
			}
			defer func() { _ = srcFile.Close() }()

			// detect image extension
			detectedExt, err := imageutil.DetectExtension(file)
			ext = &detectedExt
			if err != nil {
				httpx.JSONResponse(w, http.StatusBadRequest, httpx.Response{
					Message: "File must be an image (.jpg, .jpeg, .png, .webp, .gif)",
					Success: false,
				})
				return
			}

		} else {
			ext = nil
		}

		// bussiness logic
		err = inject.PutArticle(ctx, req, ext, oldSlug)
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
			Message: "Article updated successfully",
			Success: true,
		})
	}
}

func AdminGetArticles(inject *article.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get url params
		queryParams := r.URL.Query()

		// convert limit and page params to integer
		limit, err := strconv.Atoi(queryParams.Get("limit"))
		if err != nil {
			httpx.JSONResponse(w, http.StatusBadRequest, httpx.Response{
				Message: "Invalid limit parameter",
				Success: false,
			})
			return
		}

		page, err := strconv.Atoi(queryParams.Get("page"))
		if err != nil || page < 1 {
			httpx.JSONResponse(w, http.StatusBadRequest, httpx.Response{
				Message: "Invalid page parameter",
				Success: false,
			})
			return
		}

		// bussiness logic
		articles, meta, err := inject.AdminGetArticlesService(page, limit)
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
			Message: "Article retrieved successfully",
			Success: true,
			Data: map[string]any{
				"articles": articles,
				"meta":     meta,
			},
		})
	}
}

func AdminDeleteArticleSlug(inject *article.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		if slug == "" {
			httpx.JSONResponse(w, http.StatusBadRequest, httpx.Response{
				Message: "Article slug is required",
				Success: false,
			})
			return
		}

		err := inject.DeleteArticle(slug)
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
				Message: "internal server error",
				Success: false,
			})
			return
		}

		httpx.JSONResponse(w, http.StatusOK, httpx.Response{
			Message: "Article successfully deleted",
			Success: true,
		})
	}
}
