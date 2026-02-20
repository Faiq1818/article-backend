package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	pkg "article/internal/pkg"
	requesttype "article/internal/request_type"
	article "article/internal/services/articles"
)

func SaveArticle(inject *article.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// multipart
		const maxUploadSize = 10 << 20 // 10 MB
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
				Message: "Image file is missing or invalid",
				Success: false,
			})
			return
		}
		defer file.Close() // preventing file descriptor leak

		// decode body
		req := requesttype.SaveArticleRequest{
			Title:   r.FormValue("title"),
			Content: r.FormValue("content"),
			Image:   header,
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

		buf := make([]byte, 512)
		n, _ := srcFile.Read(buf)

		contentType := http.DetectContentType(buf[:n])

		//// check if the content type prefix is not an image (image/)
		if !strings.HasPrefix(contentType, "image/") {
			pkg.JSONResponse(w, http.StatusBadRequest, pkg.Response{
				Message: "File harus berupa gambar (.jpg, .jpeg, .png, .webp, .gif",
				Success: false,
			})
			return
		}

		ext := ".bin"
		switch contentType {
		case "image/jpeg":
			ext = ".jpg"

		case "image/png":
			ext = ".png"

		case "image/webp":
			ext = ".webp"

		case "image/gif":
			ext = ".gif"

		case "image/svg+xml":
			ext = ".svg"

		default:
			ext = ".bin"
		}
		srcFile.Seek(0, 0)

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

func GetArticleSlug(inject *article.Handler) http.HandlerFunc {
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
			Message: "Artikel berhasil didapat",
			Success: true,
			Data:    articles,
		})
	}
}
