package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"time"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
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
			w.WriteHeader(204)
			w.Write([]byte(`{"message":"User tidak ditemukan"}`))
			log.Printf("No user: %v", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write([]byte(`{"message":"Database error"}`))
		log.Printf("Database scan error: %v", err)
		return
	}

	// check and compare password
	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(req.Password))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write([]byte(`{"message":"Password salah"}`))
		fmt.Println(err)
		return
	}

	// making the jwt
	key := []byte(os.Getenv("JWT_SECRET"))
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"role": "admin",
			"name": "Faiq",
			"exp": time.Now().Add(24 * time.Hour).Unix(),
		})

	s, err := t.SignedString(key)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write([]byte(`{"message":"Error saat membuat token"}`))
		fmt.Println(err)
		return
	}

	// response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	res := map[string]any{
		"message": "Berhasil login",
		"jwt":    s,
	}
	jsonData, _ := json.Marshal(res)
	w.Write(jsonData)
}
