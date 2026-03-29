package handlers

import (
	"errors"
	"net/http"
	"strconv"

	pkg "article/internal/pkg"
	requesttype "article/internal/request_type"
	article "article/internal/services/articles"
)

func AdminSaveArticle(inject *article.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// multipart
		const maxUploadSize = 20 << 20 // 10 MB
		err := r.ParseMultipartForm(maxUploadSize)
		if err != nil {
			pkg.JSONResponse(w, http.StatusBadRequest, pkg.Response{
				Message: "Failed to parse multipart payload",
				Success: false,
			})
			return
		}

		file, header, err := r.FormFile("image")
		if err != nil {
			pkg.JSONResponse(w, http.StatusBadRequest, pkg.Response{
				Message: "Invalid limit parameter",
				Success: false,
			})
			return
		}
		defer file.Close() // preventing file descriptor leak

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
			pkg.JSONResponse(w, http.StatusBadRequest, pkg.Response{
				Message: "Validation error",
				Success: false,
				Data:    pkg.FormatValidationError(err),
			})
			return
		}

		// make dynamic image name extension
		srcFile, err := req.Image.Open()
		if err != nil {
			pkg.JSONResponse(w, http.StatusBadRequest, pkg.Response{
				Message: "Failed to open file",
				Success: false,
			})
			return
		}
		defer srcFile.Close()

		// detect image extension
		ext, err := pkg.DetectImageExtension(srcFile)
		if err != nil {
			pkg.JSONResponse(w, http.StatusBadRequest, pkg.Response{
				Message: "File must be an image (.jpg, .jpeg, .png, .webp, .gif)",
				Success: false,
			})
			return
		}

		// bussiness logic
		err = inject.SaveArticle(ctx, req, ext)
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
			pkg.JSONResponse(w, http.StatusBadRequest, pkg.Response{
				Message: "Invalid limit parameter",
				Success: false,
			})
			return
		}

		page, err := strconv.Atoi(queryParams.Get("page"))
		if err != nil || page < 1 {
			pkg.JSONResponse(w, http.StatusBadRequest, pkg.Response{
				Message: "Invalid page parameter",
				Success: false,
			})
			return
		}

		// bussiness logic
		articles, meta, err := inject.GetArticles(page, limit)
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
			pkg.JSONResponse(w, http.StatusBadRequest, pkg.Response{
				Message: "Article slug is required",
				Success: false,
			})
			return
		}

		// bussiness logic
		articles, err := inject.GetArticleSlug(slug)
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
			pkg.JSONResponse(w, http.StatusBadRequest, pkg.Response{
				Message: "Article slug is required",
				Success: false,
			})
			return
		}

		// multipart
		const maxUploadSize = 20 << 20 // 20 MB
		err := r.ParseMultipartForm(maxUploadSize)
		if err != nil {
			pkg.JSONResponse(w, http.StatusBadRequest, pkg.Response{
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
				pkg.JSONResponse(w, http.StatusBadRequest, pkg.Response{
					Message: "Image file is invalid",
					Success: false,
				})
				return
			}
		} else {
			defer file.Close() // preventing file descriptor leak
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
			pkg.JSONResponse(w, http.StatusBadRequest, pkg.Response{
				Message: "Validation error",
				Success: false,
				Data:    pkg.FormatValidationError(err),
			})
			return
		}

		var ext *string
		if req.Image != nil {
			// make dynamic image name extension
			srcFile, err := req.Image.Open()
			if err != nil {
				pkg.JSONResponse(w, http.StatusBadRequest, pkg.Response{
					Message: "Failed to open file",
					Success: false,
				})
				return
			}
			defer srcFile.Close()

			// detect image extension
			detectedExt, err := pkg.DetectImageExtension(file)
			ext = &detectedExt
			if err != nil {
				pkg.JSONResponse(w, http.StatusBadRequest, pkg.Response{
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
			pkg.JSONResponse(w, http.StatusBadRequest, pkg.Response{
				Message: "Invalid limit parameter",
				Success: false,
			})
			return
		}

		page, err := strconv.Atoi(queryParams.Get("page"))
		if err != nil || page < 1 {
			pkg.JSONResponse(w, http.StatusBadRequest, pkg.Response{
				Message: "Invalid page parameter",
				Success: false,
			})
			return
		}

		// bussiness logic
		articles, meta, err := inject.AdminGetArticlesService(page, limit)
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
			pkg.JSONResponse(w, http.StatusBadRequest, pkg.Response{
				Message: "Article slug is required",
				Success: false,
			})
			return
		}

		err := inject.DeleteArticle(slug)
		if err != nil {
			var appErr *pkg.AppError
			if errors.As(err, &appErr) {
				pkg.JSONResponse(w, appErr.Code, pkg.Response{
					Message: appErr.Message,
					Success: false,
				})
				return
			}

			pkg.JSONResponse(w, http.StatusInternalServerError, pkg.Response{
				Message: "internal server error",
				Success: false,
			})
			return
		}

		pkg.JSONResponse(w, http.StatusOK, pkg.Response{
			Message: "Article successfully deleted",
			Success: true,
		})
	}
}
