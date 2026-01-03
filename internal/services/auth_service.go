package services

import (
	"encoding/json" // Wajib ada untuk parsing body
	"fmt"
	"github.com/go-playground/validator/v10"
	// "github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

// body
type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// decode body
	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return

	}
	
	// validate body
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		http.Error(w, fmt.Sprintf("Validation error: %s", errors), http.StatusBadRequest)
		return
	}
	
	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// u := uuid.New()
	// rows, err := h.DB.Exec("INSERT INTO users (id, name, email) VALUES ($1, $2, $3);", u, "faiq", "faiq@gmail.com")
	// if err != nil {
	// 	log.Printf("Error inserting user: %v", err)
	// 	http.Error(w, "Failed to register user", http.StatusInternalServerError)
	// 	return
	// }
	// rows.Close()

	// fmt.Printf("%s", rows)
	fmt.Printf("%s", hashedPassword)
	w.Write([]byte(hashedPassword))
}
