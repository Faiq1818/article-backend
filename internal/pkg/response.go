package pkg

import (
	"encoding/json"
	// "log"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
}

func JSONResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		// log.Printf("error encoding JSON: %v", err)
	}
}

type AppError struct {
	Message string
	Code    int
	Err     error
}

func (e *AppError) Error() string {
	return e.Message
}
