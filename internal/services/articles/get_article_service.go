package article

import (
	"database/sql"

	models "article/internal/models"
	pkg "article/internal/pkg"
)

func (s *Service) GetArticles(page int, limit int) ([]models.Article, models.PaginationMeta, error) {
	// set default page and limit if too low
	p := pkg.Pagination{
		Page:  page,
		Limit: limit,
	}

	p.Normalize()
	offset := p.MakeOffset()

	articles, total, err := s.Repo.GetManyArticle(limit, offset)
	if err != nil {
		s.Logger.Error("failed get articles", "error", err)
		return []models.Article{}, models.PaginationMeta{}, &pkg.AppError{
			Message: "Artikel tidak ditemukan",
			Code:    400,
			Err:     err,
		}
	}

	meta := p.MakeMeta(total)

	return articles, meta, nil
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

func (s *Service) AdminGetArticlesService(page int, limit int) ([]models.Article, models.PaginationMeta, error) {
	p := pkg.Pagination{
		Page:  page,
		Limit: limit,
	}

	p.Normalize()
	offset := p.MakeOffset()

	articles, total, err := s.Repo.AdminGetManyArticle(p.Limit, offset)
	if err != nil {
		s.Logger.Error("failed get articles", "error", err)
		return []models.Article{}, models.PaginationMeta{}, &pkg.AppError{
			Message: "Article not found",
			Code:    400,
			Err:     err,
		}
	}

	meta := p.MakeMeta(total)

	return articles, meta, nil
}
