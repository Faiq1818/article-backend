package main

import (
	handler "article/internal/handlers"
	"log"
	"net/http"
	"os"

	"database/sql"

	"github.com/go-playground/validator/v10"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Set Header Origin
		// Best Practice: Hindari '*' jika menggunakan Cookies/Auth.
		// Gunakan origin spesifik atau validasi dinamis.
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// 2. Set Methods yang diizinkan
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		// 3. Set Headers yang diizinkan (termasuk Authorization untuk JWT)
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// 4. Handle Preflight Request (OPTIONS)
		// Browser mengirim OPTIONS sebelum request sebenarnya untuk keamanan.
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return // Stop eksekusi di sini, jangan lanjut ke handler utama
		}

		// Lanjut ke handler selanjutnya jika bukan OPTIONS
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
	mux := handler.SetupRoutes(db, validate)

	// server listen
	log.Fatal(http.ListenAndServe(":8000", CORSMiddleware(mux)))
}
