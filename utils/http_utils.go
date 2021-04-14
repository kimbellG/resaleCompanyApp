package utils

import (
	"encoding/json"
	"net/http"
)

func Respond(w http.ResponseWriter, value map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(value)
}
