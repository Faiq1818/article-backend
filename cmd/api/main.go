package main

import (
	// "article/internal/handler"
	"article/internal/services"
	"log"
	"net/http"
	"os"

	"database/sql"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// get env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// postgres db connection
	connStr := os.Getenv("DATABASE_URL_CLIENT")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	authHandler := &services.AuthHandler{
		DB: db,
	}
	
	// initiate route
	router := http.NewServeMux()

	router.HandleFunc("POST /auth/register", authHandler.Register)
	
	// server listen
	log.Fatal(http.ListenAndServe(":8000", router))
}
