package main

import (
	// "article/internal/handler"
	// "errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func register(w http.ResponseWriter, r *http.Request) {
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

func main() {
	router := http.NewServeMux()

	router.HandleFunc("POST /auth/register", register)

	log.Fatal(http.ListenAndServe(":8000", router))
}
