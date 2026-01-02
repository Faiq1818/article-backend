package main

import (
	// "article/internal/handler"
	"article/internal/services"
	"fmt"
	"log"
	"net/http"

	"database/sql"
	_ "github.com/lib/pq"
)

func main() {
	connStr := "postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT name FROM users WHERE age = $1", 2)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println(rows)

	router := http.NewServeMux()

	router.HandleFunc("POST /auth/register", services.Register)

	log.Fatal(http.ListenAndServe(":8000", router))
}
