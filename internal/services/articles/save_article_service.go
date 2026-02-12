package article

import (
	"log"
	"strings"

	pkg "article/internal/pkg"
	requesttype "article/internal/request_type"

	"github.com/google/uuid"
)

func (h *Handler) SaveArticle(req requesttype.SaveArticleRequest) error {
	// generate a slug
	slug := strings.ReplaceAll(req.Title, " ", "-")
	slug = strings.ToLower(slug)

	u := uuid.New()
	_, err := h.DB.Exec("INSERT INTO article (id, title, slug, content) VALUES ($1, $2, $3, $4);", u, req.Title, slug, req.Content)
	if err != nil {
		statusCode, clientMessage := pkg.ParsePostgresError(err)
		log.Printf("Error inserting user: %v", err)

		return &pkg.AppError{
			Message: clientMessage,
			Code:    statusCode,
			Err:     err,
		}
	}

	return nil
}
