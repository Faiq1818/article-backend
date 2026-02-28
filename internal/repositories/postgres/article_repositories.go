package postgres

import (
	"article/internal/models"
	"database/sql"
)

type ArticleRepository struct {
	DB *sql.DB
}

func NewArticleRepository(db *sql.DB) *ArticleRepository {
	return &ArticleRepository{DB: db}
}

func (r *ArticleRepository) GetManyArticle(limit int, offset int) ([]models.Article, error) {
	// query select to db
	articleData, err := r.DB.Query("SELECT id, slug, title, content, description, image_url, updated_at FROM article ORDER BY created_at DESC LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, err
	}
	defer articleData.Close()

	// build the response
	var articles []models.Article
	for articleData.Next() {
		var article models.Article
		err := articleData.Scan(
			&article.ID,
			&article.Slug,
			&article.Title,
			&article.Content,
			&article.Description,
			&article.Image_url,
			&article.Updated_at,
		)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	return articles, nil
}

func (r *ArticleRepository) GetArticleBySlug(slug string) (models.Article, error) {
	// query select to db
	articleData := r.DB.QueryRow("SELECT updated_at, id, slug, title, content, image_url FROM article WHERE slug = $1", slug)

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
		return models.Article{}, err
	}

	return article, nil
}
