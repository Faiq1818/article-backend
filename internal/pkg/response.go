package pkg

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func JSONResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(data)
}

type AppError struct {
	Message string
	Code    int
	Err     error
}

func (e *AppError) Error() string {
	return e.Message
}
