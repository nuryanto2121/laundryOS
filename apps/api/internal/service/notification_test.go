package service_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/nuryanto2121/laundryos/apps/api/internal/service"
)

type roundTripFunc func(req *http.Request) *http.Response

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func TestNotificationService_NotifyCustomerOrderCreated(t *testing.T) {
	var receivedPayload map[string]string

	mockClient := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) *http.Response {
			if req.URL.Path != "/send-message" {
				t.Errorf("unexpected path: %s", req.URL.Path)
			}
			_ = json.NewDecoder(req.Body).Decode(&receivedPayload)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"status":"sent"}`)),
				Header:     make(http.Header),
			}
		}),
	}

	svc := service.NewNotificationServiceWithClient("http://gateway.test", "6281234567890", mockClient)

	err := svc.NotifyCustomerOrderCreated(
		context.Background(),
		"6289876543210",
		"Budi",
		"ORD261119000001",
		"Cuci Setrika (5.0 KG)",
		35000,
		35000,
		"Hari ini, 22-11-2026 15:00 WIB",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if receivedPayload["phone"] != "6289876543210" {
		t.Errorf("expected customer phone 6289876543210, got %s", receivedPayload["phone"])
	}

	msg := receivedPayload["message"]
	if !strings.Contains(msg, "#261119000001") {
		t.Errorf("expected # order number in message, got: %s", msg)
	}
	if !strings.Contains(msg, "LUNAS ✅") {
		t.Errorf("expected LUNAS in message, got: %s", msg)
	}
	if !strings.Contains(msg, "Rp 35.000") {
		t.Errorf("expected Rp 35.000 in message, got: %s", msg)
	}
}

func TestNotificationService_NotifyCustomerOrderStatus(t *testing.T) {
	var receivedPayload map[string]string

	mockClient := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) *http.Response {
			_ = json.NewDecoder(req.Body).Decode(&receivedPayload)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"status":"sent"}`)),
				Header:     make(http.Header),
			}
		}),
	}

	svc := service.NewNotificationServiceWithClient("http://gateway.test", "6281234567890", mockClient)

	// Status: READY / SIAP AMBIL
	err := svc.NotifyCustomerOrderStatus(
		context.Background(),
		"6289876543210",
		"Siti",
		"ORD261119000002",
		"READY",
		"A-02",
		50000,
		20000,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	msg := receivedPayload["message"]
	if !strings.Contains(msg, "PESANAN SIAP DIAMBIL") {
		t.Errorf("expected ready message, got: %s", msg)
	}
	if !strings.Contains(msg, "Rak A-02") {
		t.Errorf("expected rack note in message, got: %s", msg)
	}
	if !strings.Contains(msg, "Rp 30.000") { // 50000 - 20000 = 30000 sisa
		t.Errorf("expected remaining balance Rp 30.000, got: %s", msg)
	}
}
