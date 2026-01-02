package main

import (
	// "article/internal/handler"
	// "errors"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func register(w http.ResponseWriter, r *http.Request) {
    password := "mySecurePassword"

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        log.Fatal(err)
    }
		
	fmt.Printf(string(hashedPassword))
	w.Write([]byte("List of users"))
	w.Write([]byte(hashedPassword))
}

func main() {
	router := http.NewServeMux()

	router.HandleFunc("POST /auth/register", register)

	log.Fatal(http.ListenAndServe(":8000", router))
}
