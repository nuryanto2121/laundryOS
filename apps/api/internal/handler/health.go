package handler

import (
	"encoding/json"
	"net/http"
	"time"
)

// healthResponse defines the standard health check response shape.
type healthResponse struct {
	Status    string `json:"status"`
	Service   string `json:"service"`
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
}

// NewHealthHandler returns an http.HandlerFunc that responds with the
// service health status. Version is injected at construction time from config.
func NewHealthHandler(version string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := healthResponse{
			Status:    "ok",
			Service:   "laundryos-api",
			Version:   version,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}
