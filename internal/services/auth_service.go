package services

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	password := "mySecurePassword"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Printf("%s", hashedPassword)
	w.Write([]byte(hashedPassword))
}
