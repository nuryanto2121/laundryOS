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
	svc := service.NewBotService(nil, "")
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
				From:     "628123456789",
				Text:     "Berapa harga cuci kiloan?",
				PushName: "Budi",
			},
			wantCategory: bot.CategoryPrice,
			wantInReply:  "DAFTAR HARGA",
		},
		{
			name: "Operating hours query matches HOURS",
			req: model.IncomingMessageRequest{
				From:     "628123456789",
				Text:     "Toko buka jam berapa?",
				PushName: "Siti",
			},
			wantCategory: bot.CategoryHours,
			wantInReply:  "JAM OPERASIONAL",
		},
		{
			name: "Address query matches ADDRESS",
			req: model.IncomingMessageRequest{
				From:     "628123456789",
				Text:     "Alamat outletnya dimana ya?",
				PushName: "Ahmad",
			},
			wantCategory: bot.CategoryAddress,
			wantInReply:  "LOKASI OUTLET",
		},
		{
			name: "Pickup request matches PICKUP",
			req: model.IncomingMessageRequest{
				From:     "628123456789",
				Text:     "Bisa tolong jemput cucian saya?",
				PushName: "Dewi",
			},
			wantCategory: bot.CategoryPickup,
			wantInReply:  "CARA REQUEST PICKUP",
		},
		{
			name: "Unknown query matches UNKNOWN",
			req: model.IncomingMessageRequest{
				From:     "628123456789",
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
