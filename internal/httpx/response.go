package httpx

import (
	"encoding/json"
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

	_ = json.NewEncoder(w).Encode(data)
}
