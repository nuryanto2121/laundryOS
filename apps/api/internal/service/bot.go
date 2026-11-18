package service

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/nuryanto2121/laundryos/apps/api/internal/bot"
	"github.com/nuryanto2121/laundryos/apps/api/internal/model"
	"github.com/nuryanto2121/laundryos/apps/api/internal/repository"
)

// BotService defines the contract for processing customer chat messages.
type BotService interface {
	ProcessMessage(ctx context.Context, req model.IncomingMessageRequest) (*model.BotResponse, error)
}

type botService struct {
	serviceRepo *repository.ServiceRepository
	matcher     *bot.Matcher
	gatewayURL  string
	httpClient  *http.Client
}

// NewBotService creates a new BotService instance.
func NewBotService(serviceRepo *repository.ServiceRepository, gatewayURL string) BotService {
	matcher := bot.NewMatcher(bot.GetKnowledgeBase())
	return &botService{
		serviceRepo: serviceRepo,
		matcher:     matcher,
		gatewayURL:  gatewayURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// ProcessMessage matches the customer message against the FAQ knowledge base, formats the reply,
// logs unknown queries, and asynchronously dispatches the response to the WhatsApp Gateway.
func (s *botService) ProcessMessage(ctx context.Context, req model.IncomingMessageRequest) (*model.BotResponse, error) {
	category, score := s.matcher.Match(req.Text)

	// Audit log for unrecognized queries
	if category == bot.CategoryUnknown {
		slog.Warn("bot: unrecognized customer query",
			"from", req.From,
			"push_name", req.PushName,
			"text", req.Text,
			"score", score,
		)
	} else {
		slog.Info("bot: matched FAQ category",
			"from", req.From,
			"category", category,
			"score", score,
		)
	}

	var services []model.ServiceItem
	if category == bot.CategoryPrice && s.serviceRepo != nil {
		activeServices, err := s.serviceRepo.ListActiveServices(ctx)
		if err != nil {
			slog.Warn("bot: failed to load active services from DB, falling back to default", "error", err)
		} else {
			services = activeServices
		}
	}

	reply := bot.FormatResponseWithServices(category, req.PushName, services)

	// Asynchronously dispatch reply to WhatsApp Gateway if gatewayURL and phone are present
	if s.gatewayURL != "" && req.From != "" {
		go s.dispatchReply(req.From, reply)
	}

	return &model.BotResponse{
		Category: category,
		Reply:    reply,
	}, nil
}

func (s *botService) dispatchReply(phone, message string) {
	payload := map[string]string{
		"phone":   phone,
		"message": message,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		slog.Error("bot: failed to marshal gateway dispatch payload", "error", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, s.gatewayURL+"/send-message", bytes.NewReader(body))
	if err != nil {
		slog.Error("bot: failed to construct gateway dispatch request", "error", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		slog.Error("bot: gateway dispatch failed", "error", err, "gateway_url", s.gatewayURL)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		slog.Warn("bot: gateway responded with non-2xx status", "status", resp.StatusCode)
	} else {
		slog.Info("bot: reply dispatched to gateway successfully", "phone", phone)
	}
}
