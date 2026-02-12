package article

import (
	"database/sql"
	"log"

	pkg "article/internal/pkg"
)

type Article struct {
	ID      string `json:"id"`
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ArticleResponse struct {
	Message string    `json:"message"`
	Data    []Article `json:"data"`
}

func (h *Handler) GetArticle(page int, limit int) ([]Article, error) {
	// making the offset
	offset := (page - 1) * limit

	// query select to db
	articleData, err := h.DB.Query("SELECT id, slug, title, content FROM article ORDER BY created_at DESC LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		log.Fatal(err)
	}
	defer articleData.Close()

	// build the response
	var articles []Article

	for articleData.Next() {
		var article Article
		if err := articleData.Scan(&article.ID, &article.Slug, &article.Title, &article.Content); err != nil {
			log.Println(err)
			continue
		}
		articles = append(articles, article)
	}

	return articles, nil
}

func (h *Handler) GetArticleSlug(slug string) (Article, error) {
	// query select to db
	articleData := h.DB.QueryRow("SELECT id, slug, title, content FROM article WHERE slug = $1", slug)

	var article Article

	err := articleData.Scan(
		&article.ID,
		&article.Slug,
		&article.Title,
		&article.Content,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No user: %v", err)
			return Article{}, &pkg.AppError{
				Message: "Artikel tidak ditemukan",
				Code:    400,
				Err:     err,
			}
		}
		log.Printf("Database scan error: %v", err)
		return Article{}, &pkg.AppError{
			Message: "Database error",
			Code:    500,
			Err:     err,
		}
	}

	return article, nil
}
