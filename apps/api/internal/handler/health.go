package handler

import (
	"encoding/json"
	"net/http"
	"time"
)

// HealthResponse defines the standard health check response shape.
type HealthResponse struct {
	Status    string `json:"status" example:"ok"`
	Service   string `json:"service" example:"laundryos-api"`
	Version   string `json:"version" example:"0.1.0"`
	Timestamp string `json:"timestamp" example:"2026-11-18T10:00:00Z"`
}

// Health godoc
// @Summary Health Check
// @Description Returns the service health status and version
// @Tags Health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func NewHealthHandler(version string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := HealthResponse{
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
