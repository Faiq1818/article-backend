package article

import (
	"database/sql"
	"log"

	models "article/internal/models"
	pkg "article/internal/pkg"
)

type ArticleResponse struct {
	Message string           `json:"message"`
	Data    []models.Article `json:"data"`
}

func (s *Service) GetArticle(page int, limit int) ([]models.Article, error) {
	// making the offset
	offset := (page - 1) * limit

	articles, err := s.Repo.GetManyArticle(limit, offset)
	if err != nil {
	}

	return articles, nil
}

func (h *Service) GetArticleSlug(slug string) (models.Article, error) {
	// query select to db
	articleData := h.DB.QueryRow("SELECT updated_at, id, slug, title, content, image_url FROM article WHERE slug = $1", slug)

	var article models.Article
	err := articleData.Scan(
		&article.Updated_at,
		&article.ID,
		&article.Slug,
		&article.Title,
		&article.Content,
		&article.Image_url,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No user: %v", err)
			return models.Article{}, &pkg.AppError{
				Message: "Artikel tidak ditemukan",
				Code:    400,
				Err:     err,
			}
		}
		log.Printf("Database scan error: %v", err)
		return models.Article{}, &pkg.AppError{
			Message: "Database error",
			Code:    500,
			Err:     err,
		}
	}

	return article, nil
}
