package json

import (
	"encoding/json"
	"net/http"
)

func Write(w http.ResponseWriter, status int, data any) {
	// Set header
	w.Header().Set("Content-Type", "application/json")
	// Set Status
	w.WriteHeader(status)
	// Return json
	json.NewEncoder(w).Encode(data)
}

func Read(r *http.Request, data any) error {
	decoder := json.NewDecoder(r.Body)
	// Don't allow other fields - for security purposes
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)

}
