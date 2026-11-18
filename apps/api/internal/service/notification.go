package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

// NotificationService defines the contract for sending business alert notifications.
type NotificationService interface {
	NotifyPickupRequest(ctx context.Context, orderNumber, customerName, customerPhone, address, serviceName, schedule string) error
}

type notificationService struct {
	gatewayURL string
	staffPhone string
	httpClient *http.Client
}

// NewNotificationService creates a new NotificationService instance.
func NewNotificationService(gatewayURL, staffPhone string) NotificationService {
	return &notificationService{
		gatewayURL: gatewayURL,
		staffPhone: strings.TrimSpace(staffPhone),
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// NotifyPickupRequest broadcasts a new pickup order alert to the staff WhatsApp number or group.
func (s *notificationService) NotifyPickupRequest(ctx context.Context, orderNumber, customerName, customerPhone, address, serviceName, schedule string) error {
	message := fmt.Sprintf(`🚨 *NOTIFIKASI PICKUP BARU (LaundryOS)* 🚨

Halo Tim Kurir! Ada permintaan penjemputan cucian baru:

• *No. Order:* %s
• *Nama Pelanggan:* %s
• *No. WhatsApp:* https://wa.me/%s
• *Alamat Penjemputan:* %s
• *Paket Layanan:* %s
• *Jadwal Pickup:* %s

Mohon kurir yang bertugas segera memproses dan menuju lokasi pelanggan. Terima kasih! 🛵💨`,
		orderNumber, customerName, customerPhone, address, serviceName, schedule,
	)

	slog.Info("notification: pickup request triggered",
		"order_number", orderNumber,
		"customer_phone", customerPhone,
		"staff_phone", s.staffPhone,
	)

	if s.gatewayURL == "" || s.staffPhone == "" {
		slog.Info("notification: skipping external gateway dispatch (staffPhone or gatewayURL not set)")
		return nil
	}

	payload := map[string]string{
		"phone":   s.staffPhone,
		"message": message,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("notification: marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.gatewayURL+"/send-message", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("notification: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		slog.Error("notification: gateway dispatch failed", "error", err)
		return nil // Non-blocking: order creation should succeed even if staff alert fails
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		slog.Warn("notification: gateway responded with non-2xx status", "status", resp.StatusCode)
	}

	return nil
}
