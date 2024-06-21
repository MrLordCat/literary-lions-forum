package handlers

import (
	"encoding/json"
	"net/http"
)

func ErrorHandler(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
