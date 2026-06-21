package helper

import (
	"encoding/json"
	"net/http"
)

func BadResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]any{
		"message": message,
	})
}

func GoodResponse(w http.ResponseWriter, statusCode int, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]any{
		"message": message,
		"data":    data,
	})
}
