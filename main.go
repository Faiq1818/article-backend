package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

func checkNumber(num int) (string, error) {
	if num < 0 {
		return "", errors.New("number is negative")
	}
	return "number is positive", nil
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	result, err := checkNumber(-5)
	if err != nil {
		fmt.Println("Error:", err) // Output: Error: number is negative
	} else {
		fmt.Println(result)
	}
	result2, err := checkNumber(-5)
	if err != nil {
		fmt.Println("Error:", err) // Output: Error: number is negative
	} else {
		fmt.Println(result)
	}

    _ = result2

	w.Write([]byte("List of users"))
}

func main() {
	router := http.NewServeMux()

	router.HandleFunc("GET /users/", getUsers)

	log.Fatal(http.ListenAndServe(":8000", router))
}
