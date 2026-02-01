package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// parseIDFromPath extracts integer ID from path after the given prefix (e.g. "/api/categories/").
// Returns (0, false) if parsing fails.
func parseIDFromPath(path, prefix string) (int, bool) {
	idStr := strings.TrimPrefix(path, prefix)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, false
	}
	return id, true
}

// writeJSON sets Content-Type and encodes v as JSON with status code.
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// writeError writes a JSON error response.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{
		"status":  "error",
		"message": message,
	})
}
