package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nuryanto2121/laundryos/apps/api/internal/handler"
	"github.com/nuryanto2121/laundryos/apps/api/internal/model"
)

type mockBotService struct {
	processFunc func(ctx context.Context, req model.IncomingMessageRequest) (*model.BotResponse, error)
}

func (m *mockBotService) ProcessMessage(ctx context.Context, req model.IncomingMessageRequest) (*model.BotResponse, error) {
	if m.processFunc != nil {
		return m.processFunc(ctx, req)
	}
	return &model.BotResponse{Category: "HELP", Reply: "Mock reply"}, nil
}

func TestBotHandler_HandleMessage(t *testing.T) {
	t.Run("Valid request returns 200", func(t *testing.T) {
		mockSvc := &mockBotService{
			processFunc: func(ctx context.Context, req model.IncomingMessageRequest) (*model.BotResponse, error) {
				return &model.BotResponse{
					Category: "PRICE",
					Reply:    "Daftar Harga",
				}, nil
			},
		}

		h := handler.NewBotHandler(mockSvc)
		body, _ := json.Marshal(model.IncomingMessageRequest{
			From:     "628123456789",
			Text:     "Berapa harga kiloan?",
			PushName: "Budi",
		})

		req := httptest.NewRequest(http.MethodPost, "/bot/messages", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		h.HandleMessage(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d. Body: %s", rec.Code, rec.Body.String())
		}

		var res model.BotResponse
		if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if res.Category != "PRICE" {
			t.Errorf("expected category PRICE, got %q", res.Category)
		}
	})

	t.Run("Missing required fields returns 400", func(t *testing.T) {
		mockSvc := &mockBotService{}
		h := handler.NewBotHandler(mockSvc)

		body, _ := json.Marshal(model.IncomingMessageRequest{
			From: "", // missing from
			Text: "", // missing text
		})

		req := httptest.NewRequest(http.MethodPost, "/bot/messages", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		h.HandleMessage(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", rec.Code)
		}
	})

	t.Run("Invalid JSON body returns 400", func(t *testing.T) {
		mockSvc := &mockBotService{}
		h := handler.NewBotHandler(mockSvc)

		req := httptest.NewRequest(http.MethodPost, "/bot/messages", bytes.NewReader([]byte("{invalid-json")))
		rec := httptest.NewRecorder()

		h.HandleMessage(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", rec.Code)
		}
	})
}
