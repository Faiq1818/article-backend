package article

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"log"
	"net/http"
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

	u := uuid.New()
	_, err = h.DB.Exec("INSERT INTO article (id, title, content) VALUES ($1, $2, $3);", u, req.Title, req.Content)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write([]byte(`{"message":"Akun gagal dibuat, pastikan email unik"`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte(`{"message":"Article berhasil dibuat"`))
}
