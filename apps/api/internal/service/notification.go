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

	"github.com/nuryanto2121/laundryos/apps/api/internal/bot"
)

// NotificationService defines the contract for sending business alert notifications.
type NotificationService interface {
	NotifyPickupRequest(ctx context.Context, orderNumber, customerName, customerPhone, address, serviceName, schedule string) error
	NotifyCustomerOrderCreated(ctx context.Context, phone, name, orderNumber, services string, totalPrice, paidAmount float64, estCompletion string) error
	NotifyCustomerOrderStatus(ctx context.Context, phone, name, orderNumber, newStatus, rackNumber string, totalPrice, paidAmount float64) error
	SendMessage(ctx context.Context, phone, message string) error
}

type notificationService struct {
	gatewayURL string
	staffPhone string
	httpClient *http.Client
}

// NewNotificationService creates a new NotificationService instance with default http client.
func NewNotificationService(gatewayURL, staffPhone string) NotificationService {
	return NewNotificationServiceWithClient(gatewayURL, staffPhone, &http.Client{
		Timeout: 15 * time.Second,
	})
}

// NewNotificationServiceWithClient creates a NotificationService instance with a custom HTTP client.
func NewNotificationServiceWithClient(gatewayURL, staffPhone string, client *http.Client) NotificationService {
	if client == nil {
		client = &http.Client{
			Timeout: 15 * time.Second,
		}
	}
	return &notificationService{
		gatewayURL: gatewayURL,
		staffPhone: strings.TrimSpace(staffPhone),
		httpClient: client,
	}
}

// SendMessage sends an arbitrary WhatsApp message to the specified recipient phone number.
func (s *notificationService) SendMessage(ctx context.Context, phone, message string) error {
	phone = strings.TrimSpace(phone)
	if s.gatewayURL == "" || phone == "" {
		slog.Info("notification: skipping external gateway dispatch (phone or gatewayURL empty)", "phone", phone)
		return nil
	}

	payload := map[string]string{
		"phone":   phone,
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
		slog.Error("notification: gateway dispatch failed", "error", err, "phone", phone)
		return nil // Non-blocking: fail-soft
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		slog.Warn("notification: gateway responded with non-2xx status", "status", resp.StatusCode, "phone", phone)
	}

	return nil
}

// NotifyPickupRequest broadcasts a new pickup order alert to the staff WhatsApp number or group.
func (s *notificationService) NotifyPickupRequest(ctx context.Context, orderNumber, customerName, customerPhone, address, serviceName, schedule string) error {
	displayOrderNum := bot.FormatDisplayOrderNumber(orderNumber)
	message := fmt.Sprintf("🚨 *NOTIFIKASI PICKUP BARU (FAN Laundry)* 🚨\n\n"+
		"Halo Tim Kurir! Ada permintaan penjemputan cucian baru:\n\n"+
		"• *No. Order:* %s\n"+
		"• *Nama Pelanggan:* %s\n"+
		"• *No. WhatsApp:* https://wa.me/%s\n"+
		"• *Alamat Penjemputan:* %s\n"+
		"• *Paket Layanan:* %s\n"+
		"• *Jadwal Pickup:* %s\n\n"+
		"Mohon kurir yang bertugas segera me-reply notifikasi ini untuk mengambil pesanan dan menuju lokasi pelanggan. Terima kasih! 🛵💨",
		displayOrderNum, customerName, customerPhone, address, serviceName, schedule,
	)

	return s.SendMessage(ctx, s.staffPhone, message)
}

// NotifyCustomerOrderCreated notifies a walk-in customer when their order is registered at the counter.
func (s *notificationService) NotifyCustomerOrderCreated(ctx context.Context, phone, name, orderNumber, services string, totalPrice, paidAmount float64, estCompletion string) error {
	displayOrderNum := bot.FormatDisplayOrderNumber(orderNumber)
	statusPembayaran := "BELUM LUNAS"
	if paidAmount >= totalPrice && totalPrice > 0 {
		statusPembayaran = "LUNAS ✅"
	}

	msg := fmt.Sprintf("🧺 *FAN Laundry — Nota Masuk Transaksi*\n\n"+
		"Halo Kak %s! Pesanan laundry Kakak telah kami terima:\n\n"+
		"• *No. Order:* %s\n"+
		"• *Layanan:* %s\n"+
		"• *Total Tagihan:* Rp %s\n"+
		"• *Status Pembayaran:* %s\n"+
		"• *Estimasi Selesai:* %s\n\n"+
		"Pakaian Kakak sedang masuk antrian cuci. Kami akan mengirimkan notifikasi saat cucian Kakak telah selesai & siap diambil. Terima kasih telah mempercayakan pakaian Kakak di FAN Laundry! 🙏✨",
		name, displayOrderNum, services, bot.FormatIDR(totalPrice), statusPembayaran, estCompletion,
	)

	return s.SendMessage(ctx, phone, msg)
}

// NotifyCustomerOrderStatus notifies the customer when their order moves to PROSES, SIAP AMBIL, or SELESAI.
func (s *notificationService) NotifyCustomerOrderStatus(ctx context.Context, phone, name, orderNumber, newStatus, rackNumber string, totalPrice, paidAmount float64) error {
	displayOrderNum := bot.FormatDisplayOrderNumber(orderNumber)
	sisaBayar := totalPrice - paidAmount
	if sisaBayar < 0 {
		sisaBayar = 0
	}

	var msg string
	switch newStatus {
	case "PROSES":
		msg = fmt.Sprintf("🫧 *Update Pesanan FAN Laundry*\n\n"+
			"Halo Kak %s! Pakaian Kakak dengan nomor order %s saat ini sedang dalam proses pencucian & perawatan oleh tim kami. Mohon ditunggu ya! 🙏",
			name, displayOrderNum,
		)
	case "READY", "SIAP AMBIL":
		paymentNote := "Status: *LUNAS* ✅"
		if sisaBayar > 0 {
			paymentNote = fmt.Sprintf("Sisa Pembayaran: *Rp %s* (Bisa tunai, transfer BCA atau QRIS)", bot.FormatIDR(sisaBayar))
		}
		rackNote := ""
		if rackNumber != "" {
			rackNote = fmt.Sprintf("• *Lokasi Simpan:* Rak %s\n", rackNumber)
		}

		msg = fmt.Sprintf("✨ *PESANAN SIAP DIAMBIL (FAN Laundry)* ✨\n\n"+
			"Kabar baik Kak %s! Pakaian Kakak sudah selesai dicuci, dikeringkan, dan dipacking rapi dengan wangi segar.\n\n"+
			"• *No. Order:* %s\n"+
			"%s"+
			"• %s\n\n"+
			"Silakan datang ke outlet FAN Laundry Kalibaru untuk pengambilan, atau balas pesan ini jika membutuhkan bantuan pengantaran kurir. Terima kasih! 🧺❤️",
			name, displayOrderNum, rackNote, paymentNote,
		)
	case "SELESAI":
		msg = fmt.Sprintf("🎉 *Terima Kasih dari FAN Laundry*\n\n"+
			"Pakaian untuk nomor order %s telah diserahkan. Terima kasih banyak Kak %s sudah mencuci di FAN Laundry Kalibaru! Semoga hari Kakak menyenangkan. Sampai jumpa di cucian berikutnya! 🌸✨",
			displayOrderNum, name,
		)
	default:
		return nil
	}

	return s.SendMessage(ctx, phone, msg)
}
