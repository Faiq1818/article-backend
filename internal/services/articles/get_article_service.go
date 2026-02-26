package article

import (
	"database/sql"
	"log"

	pkg "article/internal/pkg"
)

type Article struct {
	Updated_at  string  `json:"updated_at"`
	ID          string  `json:"id"`
	Slug        string  `json:"slug"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Content     string  `json:"content"`
	Image_url   *string `json:"image_url"`
}

type ArticleResponse struct {
	Message string    `json:"message"`
	Data    []Article `json:"data"`
}

func (h *Service) GetArticle(page int, limit int) ([]Article, error) {
	// making the offset
	offset := (page - 1) * limit

	// query select to db
	articleData, err := h.DB.Query("SELECT updated_at, id, slug, title, content, description, image_url FROM article ORDER BY created_at DESC LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		log.Fatal(err)
	}
	defer articleData.Close()

	// build the response
	var articles []Article

	for articleData.Next() {
		var article Article

		err := articleData.Scan(
			&article.Updated_at,
			&article.ID,
			&article.Slug,
			&article.Title,
			&article.Content,
			&article.Description,
			&article.Image_url,
		)
		if err != nil {
			log.Println(err)
			continue
		}
		articles = append(articles, article)
	}

	return articles, nil
}

func (h *Service) GetArticleSlug(slug string) (Article, error) {
	// query select to db
	articleData := h.DB.QueryRow("SELECT updated_at, id, slug, title, content, image_url FROM article WHERE slug = $1", slug)

	var article Article
	err := articleData.Scan(
		&article.Updated_at,
		&article.ID,
		&article.Slug,
		&article.Title,
		&article.Content,
		&article.Image_url,
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
