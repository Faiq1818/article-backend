package main

import (
	"article/internal/handlers"

	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

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
	validate := validator.New()

	handlers := &handlers.Dependency_Injection{
		DB:       db,
		Validate: validate,
	}

	mux := handlers.SetupRoutes()

	// server listen
	log.Fatal(http.ListenAndServe(":8000", CORSMiddleware(mux)))
}
