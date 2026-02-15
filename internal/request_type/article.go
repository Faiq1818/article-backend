package requesttype

import "mime/multipart"

type SaveArticleRequest struct {
	Title   string                `json:"title" validate:"required"`
	Content string                `json:"content" validate:"required"`
	Image   *multipart.FileHeader `validate:"required"`
}
