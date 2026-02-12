package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	handlers "article/internal/handlers"
	middlewares "article/internal/middlewares"

	"github.com/go-playground/validator/v10"
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

	// validator initiate
	validate := validator.New(validator.WithRequiredStructEnabled())

	mux := handlers.SetupRoutes(db, validate)

	// server listen
	// Wrapping up the mux inside the corsMiddleware so it can smuggle the cors header
	log.Fatal(http.ListenAndServe(":8000", middlewares.CORSMiddleware(mux)))
}
