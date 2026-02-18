package repository

import (
	"article/internal/models"
	"context"
	"database/sql"
)

type ArticleRepository interface {
	Create(ctx context.Context, article *models.Article) error
	FindAll(ctx context.Context, limit, offset int) ([]models.Article, error)
	FindBySlug(ctx context.Context, slug string) (*models.Article, error)
}

type articleRepository struct {
	db *sql.DB
}

func NewArticleRepository(db *sql.DB) ArticleRepository {
	return &articleRepository{db: db}
}

func (r *articleRepository) Create(ctx context.Context, article *models.Article) error {
	query := "INSERT INTO article (id, title, slug, content) VALUES ($1, $2, $3, $4)"
	_, err := r.db.ExecContext(ctx, query, article.ID, article.Title, article.Slug, article.Content)
	return err
}

func (r *articleRepository) FindAll(ctx context.Context, limit, offset int) ([]models.Article, error) {
	query := "SELECT id, slug, title, content, updated_at FROM article ORDER BY created_at DESC LIMIT $1 OFFSET $2"
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var a models.Article
		if err := rows.Scan(&a.ID, &a.Slug, &a.Title, &a.Content, &a.UpdatedAt); err != nil {
			return nil, err
		}
		articles = append(articles, a)
	}
	return articles, nil
}

func (r *articleRepository) FindBySlug(ctx context.Context, slug string) (*models.Article, error) {
	query := "SELECT id, slug, title, content, updated_at FROM article WHERE slug = $1"
	row := r.db.QueryRowContext(ctx, query, slug)

	var a models.Article
	err := row.Scan(&a.ID, &a.Slug, &a.Title, &a.Content, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}
