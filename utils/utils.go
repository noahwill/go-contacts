package utils

import (
	"encoding/json"
	"net/http"
)

// Build JSON messages
func Message(status bool, message string) (map[string]interface{}) {
	return map[string]interface{} {"status": status, "mesage": message}
}

// Return JSON response
func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}