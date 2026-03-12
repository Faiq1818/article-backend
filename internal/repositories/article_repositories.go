package repositories

import (
	"article/internal/models"
	requesttype "article/internal/request_type"
)

type ArticleRepository interface {
	GetManyArticle(limit int, offset int) ([]models.Article, int, error)
	GetArticleBySlug(slug string) (models.Article, error)
	SaveArticle(req requesttype.SaveArticleRequest, imgUrl string, slugGenerate string) error
	PutArticle(req requesttype.PutArticleRequest, imgUrl string, slugGenerate string, oldSlug string) error

	AdminGetManyArticle(limit int, offset int) ([]models.Article, int, error)
}
