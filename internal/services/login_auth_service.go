package services

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	"log"
	"net/http"
)

// body
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// decode body
	var req LoginRequest
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

	// query get user data from req.email
	userData := h.DB.QueryRow("SELECT id, name, email, password FROM users WHERE email = ($1);", req.Email)

	var id, name, email, password string
	err = userData.Scan(&id, &name, &email, &password)
	if err != nil {
		if err == sql.ErrNoRows {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			w.Write([]byte(`{"message":"User tidak ditemukan", "success": false}`))
			log.Printf("No user: %v", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write([]byte(`{"message":"Database error", "success": false}`))
		log.Printf("Database scan error: %v", err)
		return
	}

	// check and compare password
	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(req.Password))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write([]byte(`{"message":"Password salah", "success": false}`))
		fmt.Println(err)
		return
	}

	// response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte(`{"message":"Berhasil", "success": true}`))
}
