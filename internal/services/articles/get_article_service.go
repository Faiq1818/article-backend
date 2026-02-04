package article

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func (h *AuthHandler) GetArticle(w http.ResponseWriter, r *http.Request) {
	// get url params
	queryParams := r.URL.Query()

	// convert limit and page params to integer
	limit, err := strconv.Atoi(queryParams.Get("limit"))
	if err != nil {
		http.Error(w, "invalid limit", http.StatusBadRequest)
		return
	}

	page, err := strconv.Atoi(queryParams.Get("page"))
	if err != nil || page < 1 {
		http.Error(w, "invalid page", http.StatusBadRequest)
		return
	}

	// making the offset
	offset := (page - 1) * limit

	// query select to db
	articleData, err := h.DB.Query("SELECT id, slug, title, content FROM article ORDER BY created_at DESC LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		log.Fatal(err)
	}
	defer articleData.Close()

	type Article struct {
		ID      string `json:"id"`
		Slug    string `json:"slug"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}

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

	type ArticleResponse struct {
		Message string    `json:"message"`
		Data    []Article `json:"data"`
	}

	response := ArticleResponse{
		Message: "Article berhasil dibuat",
		Data:    articles,
	}

	// write and send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
