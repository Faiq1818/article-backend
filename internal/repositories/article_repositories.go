package repositories

import "article/internal/models"

type ArticleRepository interface {
	GetManyArticle(limit int, offset int) ([]models.Article, error)
	GetArticleBySlug(slug string) (models.Article, error)
}
