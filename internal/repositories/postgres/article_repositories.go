package postgres

import (
	"article/internal/models"
	"database/sql"

	requesttype "article/internal/request_type"

	"github.com/google/uuid"
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
	articleData := r.DB.QueryRow("SELECT updated_at, id, slug, title, description, content, image_url FROM article WHERE slug = $1", slug)

	var article models.Article
	err := articleData.Scan(
		&article.Updated_at,
		&article.ID,
		&article.Slug,
		&article.Title,
		&article.Description,
		&article.Content,
		&article.Image_url,
	)
	if err != nil {
		return models.Article{}, err
	}

	return article, nil
}

func (r *ArticleRepository) SaveArticle(req requesttype.SaveArticleRequest, imageUrl string, slugGenerate string) error {
	// db push
	u := uuid.New()
	_, err := r.DB.Exec("INSERT INTO article (id, title, slug, description, content, image_url) VALUES ($1, $2, $3, $4, $5, $6);", u, req.Title, slugGenerate, req.Description, req.Content, imageUrl)
	if err != nil {
		return err
	}

	return nil
}

func (r *ArticleRepository) PutArticle(req requesttype.PutArticleRequest, imageUrl string, slugGenerate string, oldSlug string) error {
	// db push
	_, err := r.DB.Exec(`
		UPDATE article 
		SET title = $1, 
			description = $2, 
			content = $3, 
			image_url = $4
		WHERE slug = $5;
	`, req.Title, req.Description, req.Content, imageUrl, oldSlug)
	if err != nil {
		return err
	}

	return nil
}
