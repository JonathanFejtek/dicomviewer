package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func parseURLParam(r *http.Request, key string) (string, error) {
	value := chi.URLParam(r, key)
	if value == "" {
		return "", fmt.Errorf("url param %s is empty", key)
	}

	return value, nil
}

func writeJSONResponse(w http.ResponseWriter, i interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(i)
}

func writeJSONError(w http.ResponseWriter, httpStatusCode int, err error) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)

	var errorResp = struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	}

	return json.NewEncoder(w).Encode(errorResp)
}
