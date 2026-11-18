package service_test

import (
	"context"
	"strings"
	"testing"

	"github.com/nuryanto2121/laundryos/apps/api/internal/bot"
	"github.com/nuryanto2121/laundryos/apps/api/internal/model"
	"github.com/nuryanto2121/laundryos/apps/api/internal/service"
)

func TestBotService_ProcessMessage(t *testing.T) {
	// Standalone bot service without DB repo, tests fallback & matching
	svc := service.NewBotService(service.BotServiceConfig{})
	ctx := context.Background()

	tests := []struct {
		name         string
		req          model.IncomingMessageRequest
		wantCategory string
		wantInReply  string
	}{
		{
			name: "Price query matches PRICE",
			req: model.IncomingMessageRequest{
				From:     "628100000001",
				Text:     "Berapa harga cuci kiloan?",
				PushName: "Budi",
			},
			wantCategory: bot.CategoryPrice,
			wantInReply:  "DAFTAR HARGA",
		},
		{
			name: "Operating hours query matches HOURS",
			req: model.IncomingMessageRequest{
				From:     "628100000002",
				Text:     "Toko buka jam berapa?",
				PushName: "Siti",
			},
			wantCategory: bot.CategoryHours,
			wantInReply:  "JAM OPERASIONAL",
		},
		{
			name: "Address query matches ADDRESS",
			req: model.IncomingMessageRequest{
				From:     "628100000003",
				Text:     "Alamat outletnya dimana ya?",
				PushName: "Ahmad",
			},
			wantCategory: bot.CategoryAddress,
			wantInReply:  "LOKASI OUTLET",
		},
		{
			name: "Pickup request matches PICKUP",
			req: model.IncomingMessageRequest{
				From:     "628100000004",
				Text:     "Bisa tolong jemput cucian saya?",
				PushName: "Dewi",
			},
			wantCategory: bot.CategoryPickup,
			wantInReply:  "REQUEST PENJEMPUTAN",
		},
		{
			name: "Status query matches STATUS",
			req: model.IncomingMessageRequest{
				From:     "628100000006",
				Text:     "Cek status cucian saya dong min",
				PushName: "Lisa",
			},
			wantCategory: bot.CategoryStatus,
			wantInReply:  "STATUS CUCIAN",
		},
		{
			name: "Unknown query matches UNKNOWN",
			req: model.IncomingMessageRequest{
				From:     "628100000005",
				Text:     "Beli tiket pesawat ke Bali",
				PushName: "Rian",
			},
			wantCategory: bot.CategoryUnknown,
			wantInReply:  "belum memahami",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := svc.ProcessMessage(ctx, tt.req)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if res.Category != tt.wantCategory {
				t.Errorf("expected category %q, got %q", tt.wantCategory, res.Category)
			}
			if !strings.Contains(res.Reply, tt.wantInReply) {
				t.Errorf("expected reply to contain %q, got: %s", tt.wantInReply, res.Reply)
			}
		})
	}
}

func TestBotService_MultiTurnPickup(t *testing.T) {
	svc := service.NewBotService(service.BotServiceConfig{})
	ctx := context.Background()
	phone := "628999999999"

	// Turn 1: Trigger pickup
	res1, err := svc.ProcessMessage(ctx, model.IncomingMessageRequest{
		From:     phone,
		Text:     "Kak tolong jemput cucian ya",
		PushName: "Doni",
	})
	if err != nil {
		t.Fatalf("turn 1 error: %v", err)
	}
	if !strings.Contains(res1.Reply, "Alamat Lengkap Penjemputan") {
		t.Errorf("turn 1 expected prompt for address, got: %s", res1.Reply)
	}

	// Turn 2: Provide Address
	res2, err := svc.ProcessMessage(ctx, model.IncomingMessageRequest{
		From: phone,
		Text: "Jl. Sudirman No. 12, Kos Sakura Kamar 10",
	})
	if err != nil {
		t.Fatalf("turn 2 error: %v", err)
	}
	if !strings.Contains(res2.Reply, "pilih paket layanan") {
		t.Errorf("turn 2 expected service choice, got: %s", res2.Reply)
	}

	// Turn 3: Select Service Option 1
	res3, err := svc.ProcessMessage(ctx, model.IncomingMessageRequest{
		From: phone,
		Text: "1",
	})
	if err != nil {
		t.Fatalf("turn 3 error: %v", err)
	}
	if !strings.Contains(res3.Reply, "Kapan kurir kami bisa menjemput") {
		t.Errorf("turn 3 expected schedule prompt, got: %s", res3.Reply)
	}

	// Turn 4: Provide Pickup Schedule
	res4, err := svc.ProcessMessage(ctx, model.IncomingMessageRequest{
		From: phone,
		Text: "Nanti sore jam 16:00 WIB",
	})
	if err != nil {
		t.Fatalf("turn 4 error: %v", err)
	}
	if !strings.Contains(res4.Reply, "RINGKASAN REQUEST PENJEMPUTAN") {
		t.Errorf("turn 4 expected confirmation summary, got: %s", res4.Reply)
	}

	// Turn 5: Confirm with OK
	res5, err := svc.ProcessMessage(ctx, model.IncomingMessageRequest{
		From: phone,
		Text: "OK",
	})
	if err != nil {
		t.Fatalf("turn 5 error: %v", err)
	}
	if !strings.Contains(res5.Reply, "REQUEST PICKUP BERHASIL DIBUAT") {
		t.Errorf("turn 5 expected success confirmation, got: %s", res5.Reply)
	}

	// Verify session was cleared after completion
	res6, err := svc.ProcessMessage(ctx, model.IncomingMessageRequest{
		From: phone,
		Text: "Halo",
	})
	if err != nil {
		t.Fatalf("turn 6 error: %v", err)
	}
	if res6.Category != bot.CategoryHelp {
		t.Errorf("turn 6 expected idle state matching HELP, got %q", res6.Category)
	}
}

func TestBotService_PickupCancellation(t *testing.T) {
	svc := service.NewBotService(service.BotServiceConfig{})
	ctx := context.Background()
	phone := "628888888888"

	// Turn 1: Trigger pickup
	_, _ = svc.ProcessMessage(ctx, model.IncomingMessageRequest{
		From:     phone,
		Text:     "pickup",
		PushName: "Andi",
	})

	// Turn 2: Cancel
	res, err := svc.ProcessMessage(ctx, model.IncomingMessageRequest{
		From: phone,
		Text: "BATAL",
	})
	if err != nil {
		t.Fatalf("cancel error: %v", err)
	}
	if !strings.Contains(res.Reply, "dibatalkan") {
		t.Errorf("expected cancellation message, got: %s", res.Reply)
	}
}
