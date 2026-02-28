package repositories

import "article/internal/models"

type ArticleRepository interface {
	GetManyArticle(limit int, offset int) ([]models.Article, error)
}
