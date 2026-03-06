package article

import (
	"database/sql"

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
		s.Logger.Error("failed get articles", "error", err)
		return []models.Article{}, &pkg.AppError{
			Message: "Artikel tidak ditemukan",
			Code:    400,
			Err:     err,
		}
	}

	return articles, nil
}

func (s *Service) GetArticleSlug(slug string) (models.Article, error) {
	// query select to db
	article, err := s.Repo.GetArticleBySlug(slug)
	if err != nil {
		if err == sql.ErrNoRows {
			s.Logger.Info("Article not found", "error", err)
			return models.Article{}, &pkg.AppError{
				Message: "Artikel tidak ditemukan",
				Code:    400,
				Err:     err,
			}
		}

		s.Logger.Error("Database scan error", "error", err)
		return models.Article{}, &pkg.AppError{
			Message: "Database error",
			Code:    500,
			Err:     err,
		}
	}

	s.Logger.Info("Get article from slug success")

	return article, nil
}
