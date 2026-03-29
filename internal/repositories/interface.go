package repositories

import (
	"context"
	"io"

	models "article/internal/models"
	requesttype "article/internal/request_type"

	"github.com/google/uuid"
)

type AuthRepository interface {
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(u *models.User) error
	CheckUserId(id uuid.UUID) (bool, error)
}

type ArticleRepository interface {
	GetManyArticle(limit int, offset int) ([]models.Article, int, error)
	GetArticleBySlug(slug string) (models.Article, error)
	SaveArticle(ctx context.Context, req requesttype.SaveArticleRequest, imgUrl string, slugGenerate string) error
	PutArticle(req requesttype.PutArticleRequest, imgUrl string, slugGenerate string, oldSlug string) error
	DeleteArticle(slug string) error

	AdminGetManyArticle(limit int, offset int) ([]models.Article, int, error)
}

type S3Repository interface {
	// GetUserByEmail(email string) (*models.User, error)
	UploadObject(ctx context.Context, key string, fileBody io.Reader) (string, error)
}
