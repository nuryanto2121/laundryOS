package service_test

import (
	"context"
	"strings"
	"testing"

	"github.com/nuryanto2121/laundryos/apps/api/internal/bot"
	"github.com/nuryanto2121/laundryos/apps/api/internal/model"
	"github.com/nuryanto2121/laundryos/apps/api/internal/service"
	"github.com/stretchr/testify/mock"
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

func TestBotService_ConversationalRealWorldScenario(t *testing.T) {
	svc := service.NewBotService(service.BotServiceConfig{
		GatewayURL: "http://mock-gateway",
	})

	ctx := context.Background()
	phone := "6281299998888"

	// Step 1: Customer says "saya mau laundry mba"
	res1, err := svc.ProcessMessage(ctx, model.IncomingMessageRequest{
		From:     phone,
		Text:     "saya mau laundry mba",
		PushName: "Nuryanto",
	})
	if err != nil {
		t.Fatalf("turn 1 error: %v", err)
	}
	if res1.Category != bot.CategoryPickup {
		t.Errorf("expected CategoryPickup, got: %s", res1.Category)
	}
	if !strings.Contains(res1.Reply, "Alamat") {
		t.Errorf("expected address prompt, got: %s", res1.Reply)
	}

	// Step 2: Customer asks "Harga"
	res2, err := svc.ProcessMessage(ctx, model.IncomingMessageRequest{
		From:     phone,
		Text:     "Harga",
		PushName: "Nuryanto",
	})
	if err != nil {
		t.Fatalf("turn 2 error: %v", err)
	}
	if !strings.Contains(res2.Reply, "DAFTAR HARGA") {
		t.Errorf("expected price list, got: %s", res2.Reply)
	}

	// Step 3: Customer chooses service "cuci regular iah"
	res3, err := svc.ProcessMessage(ctx, model.IncomingMessageRequest{
		From:     phone,
		Text:     "cuci regular iah",
		PushName: "Nuryanto",
	})
	if err != nil {
		t.Fatalf("turn 3 error: %v", err)
	}
	if !strings.Contains(res3.Reply, "Cuci Komplit Reguler") {
		t.Errorf("expected confirmation of service choice, got: %s", res3.Reply)
	}
	if !strings.Contains(res3.Reply, "Alamat") {
		t.Errorf("expected address prompt, got: %s", res3.Reply)
	}

	// Step 4: Customer provides address "Jl. Mawar No. 12"
	res4, err := svc.ProcessMessage(ctx, model.IncomingMessageRequest{
		From:     phone,
		Text:     "Jl. Mawar No. 12",
		PushName: "Nuryanto",
	})
	if err != nil {
		t.Fatalf("turn 4 error: %v", err)
	}
	if !strings.Contains(res4.Reply, "jam") && !strings.Contains(res4.Reply, "Jam") {
		t.Errorf("expected schedule/time prompt, got: %s", res4.Reply)
	}

	// Step 5: Customer provides schedule "jam 2 siang ya"
	res5, err := svc.ProcessMessage(ctx, model.IncomingMessageRequest{
		From:     phone,
		Text:     "jam 2 siang ya",
		PushName: "Nuryanto",
	})
	if err != nil {
		t.Fatalf("turn 5 error: %v", err)
	}
	if !strings.Contains(res5.Reply, "RINGKASAN") {
		t.Errorf("expected summary, got: %s", res5.Reply)
	}

	// Step 6: Customer confirms with casual "oke"
	res6, err := svc.ProcessMessage(ctx, model.IncomingMessageRequest{
		From:     phone,
		Text:     "oke",
		PushName: "Nuryanto",
	})
	if err != nil {
		t.Fatalf("turn 6 error: %v", err)
	}
	if !strings.Contains(res6.Reply, "BERHASIL") {
		t.Errorf("expected success order creation, got: %s", res6.Reply)
	}
	if !strings.Contains(res6.Reply, "• *No. Order:* #") {
		t.Errorf("expected hash ready-to-copy order number, got: %s", res6.Reply)
	}
}

func TestBotService_StaffGroupPickupClaim(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockOrderRepository)
	orderSvc := service.NewOrderService(mockRepo)

	mockRepo.On("GetByOrderNumber", ctx, "ORD261118000001").Return(model.OrderTrackingDetail{
		ID:           "test-order-id",
		OrderNumber:  "ORD261118000001",
		CustomerID:   "test-cust-id",
		CustomerName: "Budi",
		ServiceName:  "Cuci Setrika Regular",
		Status:       model.OrderStatusWaitingPickup,
		Notes:        "Jadwal: jam 2 siang | Alamat: Jl. Mawar 12",
	}, nil)

	mockRepo.On("Update", ctx, "test-order-id", mock.Anything).Return(model.Order{
		ID:          "test-order-id",
		OrderNumber: "ORD261118000001",
		Notes:       "Jadwal: jam 2 siang | Alamat: Jl. Mawar 12 | Kurir: Rian",
	}, nil)
	mockRepo.On("UpdateStatus", ctx, "test-order-id", model.OrderStatusConfirmed).Return(model.Order{
		ID:          "test-order-id",
		OrderNumber: "ORD261118000001",
		Status:      model.OrderStatusConfirmed,
	}, nil)

	svc := service.NewBotService(service.BotServiceConfig{
		OrderService: orderSvc,
		StaffPhone:   "120363416194985610@g.us",
	})

	// Scenario 1: Staff quotes the pickup notification and replies "saya ambil yaa"
	req := model.IncomingMessageRequest{
		From:        "120363416194985610@g.us",
		Text:        "saya ambil yaa",
		PushName:    "Rian",
		IsGroup:     true,
		GroupJID:    "120363416194985610@g.us",
		Participant: "6285712345678",
		QuotedText:  "🚨 *NOTIFIKASI PICKUP BARU (FAN Laundry)* 🚨\n• *No. Order:* #261118000001\n• *Nama Pelanggan:* Budi",
	}

	res, err := svc.ProcessMessage(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Category != "GROUP_STAFF_CLAIM" {
		t.Errorf("expected GROUP_STAFF_CLAIM, got: %s", res.Category)
	}
	if !strings.Contains(res.Reply, "PENJEMPUTAN DICONFIRM") {
		t.Errorf("expected ACK confirmation, got: %s", res.Reply)
	}
	if !strings.Contains(res.Reply, "#261118000001") {
		t.Errorf("expected hash order number in ACK, got: %s", res.Reply)
	}
	if !strings.Contains(res.Reply, "Rian") {
		t.Errorf("expected staff name in ACK, got: %s", res.Reply)
	}

	// Scenario 2: Random group banter without order number is ignored
	randomReq := model.IncomingMessageRequest{
		From:     "120363416194985610@g.us",
		Text:     "siapa yang piket besok?",
		IsGroup:  true,
		GroupJID: "120363416194985610@g.us",
	}
	res2, err := svc.ProcessMessage(ctx, randomReq)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res2.Category != "GROUP_IGNORED" {
		t.Errorf("expected GROUP_IGNORED, got: %s", res2.Category)
	}
}

func TestBotService_StaffGroupWeighing_And_Completion(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockOrderRepository)
	orderSvc := service.NewOrderService(mockRepo)

	mockRepo.On("GetByOrderNumber", ctx, "ORD261118000002").Return(model.OrderTrackingDetail{
		ID:           "order-kiloan-id",
		OrderNumber:  "ORD261118000002",
		CustomerID:   "cust-2",
		CustomerName: "Siti",
		ServiceName:  "Cuci Setrika Regular",
		DurationDay:  3,
		Status:       model.OrderStatusWaitingPickup,
		Notes:        "Alamat: Kalibaru",
	}, nil)

	mockRepo.On("Update", ctx, "order-kiloan-id", mock.Anything).Return(model.Order{
		ID:          "order-kiloan-id",
		OrderNumber: "ORD261118000002",
	}, nil)
	mockRepo.On("UpdateStatus", ctx, "order-kiloan-id", model.OrderStatusInProgress).Return(model.Order{
		ID:          "order-kiloan-id",
		OrderNumber: "ORD261118000002",
		Status:      model.OrderStatusInProgress,
	}, nil)
	mockRepo.On("UpdateStatus", ctx, "order-kiloan-id", model.OrderStatusCompleted).Return(model.Order{
		ID:          "order-kiloan-id",
		OrderNumber: "ORD261118000002",
		Status:      model.OrderStatusCompleted,
	}, nil)

	svc := service.NewBotService(service.BotServiceConfig{
		OrderService: orderSvc,
		StaffPhone:   "120363416194985610@g.us",
	})

	// 1. TIMBANG #261118000002 5.5
	weighReq := model.IncomingMessageRequest{
		From:     "120363416194985610@g.us",
		Text:     "TIMBANG #261118000002 5.5",
		PushName: "Admin Outlet",
		IsGroup:  true,
		GroupJID: "120363416194985610@g.us",
	}
	weighRes, err := svc.ProcessMessage(ctx, weighReq)
	if err != nil {
		t.Fatalf("weigh error: %v", err)
	}
	if weighRes.Category != "GROUP_STAFF_WEIGH_SUCCESS" {
		t.Errorf("expected GROUP_STAFF_WEIGH_SUCCESS, got: %s", weighRes.Category)
	}
	if !strings.Contains(weighRes.Reply, "PENIMBANGAN BERHASIL DICATAT") {
		t.Errorf("expected weigh ACK, got: %s", weighRes.Reply)
	}

	// 2. SELESAI #261118000002
	compReq := model.IncomingMessageRequest{
		From:     "120363416194985610@g.us",
		Text:     "SELESAI #261118000002",
		PushName: "Admin Outlet",
		IsGroup:  true,
		GroupJID: "120363416194985610@g.us",
	}
	compRes, err := svc.ProcessMessage(ctx, compReq)
	if err != nil {
		t.Fatalf("completion error: %v", err)
	}
	if compRes.Category != "GROUP_STAFF_COMPLETE_SUCCESS" {
		t.Errorf("expected GROUP_STAFF_COMPLETE_SUCCESS, got: %s", compRes.Category)
	}
	if !strings.Contains(compRes.Reply, "ORDER COMPLETED") {
		t.Errorf("expected completion ACK, got: %s", compRes.Reply)
	}
}

func TestBotService_StaffGroupWeighing_OverloadTabung(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockOrderRepository)
	orderSvc := service.NewOrderService(mockRepo)

	mockRepo.On("GetByOrderNumber", ctx, "ORD261118000003").Return(model.OrderTrackingDetail{
		ID:           "order-ss-id",
		OrderNumber:  "ORD261118000003",
		CustomerID:   "cust-3",
		CustomerName: "Ahmad",
		ServiceName:  "Cuci Kering (Self Service)",
		DurationDay:  0,
		Status:       model.OrderStatusWaitingPickup,
		Notes:        "Self Service",
	}, nil)

	svc := service.NewBotService(service.BotServiceConfig{
		OrderService: orderSvc,
		StaffPhone:   "120363416194985610@g.us",
	})

	// TIMBANG #261118000003 9.5 (Exceeds 8.5 kg tolerance!)
	weighReq := model.IncomingMessageRequest{
		From:     "120363416194985610@g.us",
		Text:     "TIMBANG #261118000003 9.5",
		PushName: "Admin Outlet",
		IsGroup:  true,
		GroupJID: "120363416194985610@g.us",
	}
	weighRes, err := svc.ProcessMessage(ctx, weighReq)
	if err != nil {
		t.Fatalf("overload weigh error: %v", err)
	}
	if weighRes.Category != "GROUP_STAFF_WEIGH_OVERLOAD" {
		t.Errorf("expected GROUP_STAFF_WEIGH_OVERLOAD, got: %s", weighRes.Category)
	}
	if !strings.Contains(weighRes.Reply, "MENUNGGU KONFIRMASI PELANGGAN") {
		t.Errorf("expected overload prompt ACK, got: %s", weighRes.Reply)
	}
}

func TestBotService_CheckStatusExplicitOrderNumbers(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockOrderRepository)
	orderSvc := service.NewOrderService(mockRepo)

	targetOrder := model.OrderTrackingDetail{
		ID:           "target-order-id",
		OrderNumber:  "ORD261121000001",
		CustomerID:   "cust-nebula",
		CustomerName: "Nebula",
		ServiceName:  "Cuci Boneka",
		DurationDay:  3,
		Status:       model.OrderStatusCompleted,
		TotalPrice:   50000,
	}

	mockRepo.On("GetByOrderNumber", ctx, "ORD261121000001").Return(targetOrder, nil)
	mockRepo.On("GetByOrderNumber", ctx, "261121000001").Return(model.OrderTrackingDetail{}, nil)
	mockRepo.On("GetByOrderNumber", ctx, "ORD999999999999").Return(model.OrderTrackingDetail{}, nil)
	mockRepo.On("GetByOrderNumber", ctx, "999999999999").Return(model.OrderTrackingDetail{}, nil)

	svc := service.NewBotService(service.BotServiceConfig{
		OrderService: orderSvc,
	})

	testCases := []struct {
		name       string
		inputText  string
		wantInBody string
		wantCat    string
	}{
		{
			name:       "Status with hash prefix #261121000001",
			inputText:  "Status #261121000001",
			wantInBody: "#261121000001",
			wantCat:    bot.CategoryStatus,
		},
		{
			name:       "Status with pure numeric 261121000001",
			inputText:  "Status 261121000001",
			wantInBody: "#261121000001",
			wantCat:    bot.CategoryStatus,
		},
		{
			name:       "Typo in status keyword Stataus #261121000001",
			inputText:  "Stataus #261121000001",
			wantInBody: "#261121000001",
			wantCat:    bot.CategoryStatus,
		},
		{
			name:       "Non-existent order number does not fallback to random order",
			inputText:  "Status #999999999999",
			wantInBody: "tidak ditemukan di sistem kami",
			wantCat:    bot.CategoryStatus,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := svc.ProcessMessage(ctx, model.IncomingMessageRequest{
				From:     "6285777374040",
				Text:     tc.inputText,
				PushName: "Nuryanto",
			})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if res.Category != tc.wantCat {
				t.Errorf("expected category %s, got: %s", tc.wantCat, res.Category)
			}
			if !strings.Contains(res.Reply, tc.wantInBody) {
				t.Errorf("expected reply to contain %q, but got:\n%s", tc.wantInBody, res.Reply)
			}
		})
	}
}
