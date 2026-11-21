package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// GatewayHandler proxies requests to the WhatsApp Gateway service.
type GatewayHandler struct {
	gatewayURL string
	client     *http.Client
}

// NewGatewayHandler creates a new GatewayHandler instance.
func NewGatewayHandler(gatewayURL string) *GatewayHandler {
	return &GatewayHandler{
		gatewayURL: gatewayURL,
		client:     &http.Client{Timeout: 5 * time.Second},
	}
}

// GetStatus proxies the WhatsApp connection status and QR code from the gateway.
func (h *GatewayHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	if h.gatewayURL == "" {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"success": true,
			"data": map[string]any{
				"state": "DISCONNECTED",
				"qr":    nil,
				"bot":   nil,
			},
		})
		return
	}

	resp, err := h.client.Get(h.gatewayURL + "/status")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"error":   err.Error(),
			"data": map[string]any{
				"state": "DISCONNECTED",
			},
		})
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}

// GetGroups proxies the list of participating WhatsApp groups from the gateway.
func (h *GatewayHandler) GetGroups(w http.ResponseWriter, r *http.Request) {
	if h.gatewayURL == "" {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"success": true,
			"data":    []any{},
		})
		return
	}

	resp, err := h.client.Get(h.gatewayURL + "/groups")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"error":   err.Error(),
			"data":    []any{},
		})
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}
