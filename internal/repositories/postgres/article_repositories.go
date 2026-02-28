package postgres

import (
	"article/internal/models"
	"database/sql"
	"log"
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
		log.Fatal(err)
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
			// s.Logger.Error("Error when append article data from db to struct")
			continue
		}
		articles = append(articles, article)
	}

	return articles, nil
}
