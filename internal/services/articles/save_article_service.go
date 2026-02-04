package article

import (
	"article/internal/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type SaveArticleRequest struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

func (h *AuthHandler) SaveArticle(w http.ResponseWriter, r *http.Request) {
	// decode body
	var req SaveArticleRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// validate body
	err = h.Validate.Struct(req)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		http.Error(w, fmt.Sprintf("Validation error: %s", errors), http.StatusBadRequest)
		return
	}

	// generate a slug
	slug := strings.ReplaceAll(req.Title, " ", "-")
	slug = strings.ToLower(slug)

	u := uuid.New()
	_, err = h.DB.Exec("INSERT INTO article (id, title, slug, content) VALUES ($1, $2, $3, $4);", u, req.Title, slug, req.Content)
	if err != nil {
		fmt.Printf("%#v\n", err)
		statusCode, clientMessage := utils.ParsePostgresError(err)

		log.Printf("Error inserting user: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		w.Write([]byte(`{"message":"` + clientMessage + `"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte(`{"message":"Article berhasil dibuat"}`))
}
