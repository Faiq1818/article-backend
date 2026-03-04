package repositories

import (
	"article/internal/models"
	requesttype "article/internal/request_type"
)

type ArticleRepository interface {
	GetManyArticle(limit int, offset int) ([]models.Article, error)
	GetArticleBySlug(slug string) (models.Article, error)
	SaveArticle(req requesttype.SaveArticleRequest, imgUrl string, slugGenerate string) error
}
