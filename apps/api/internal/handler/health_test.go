package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nuryanto2121/laundryos/apps/api/internal/handler"
)

func TestHealthHandler(t *testing.T) {
	h := handler.NewHealthHandler("0.1.0")

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	h(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if body["status"] != "ok" {
		t.Errorf("expected status ok, got %q", body["status"])
	}
	if body["service"] != "laundryos-api" {
		t.Errorf("expected service laundryos-api, got %q", body["service"])
	}
	if body["version"] != "0.1.0" {
		t.Errorf("expected version 0.1.0, got %q", body["version"])
	}
	if body["timestamp"] == "" {
		t.Error("expected non-empty timestamp")
	}
}
