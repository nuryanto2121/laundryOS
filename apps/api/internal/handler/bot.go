package handler

import (
	"encoding/json"
	"net/http"

	"github.com/nuryanto2121/laundryos/apps/api/internal/model"
	"github.com/nuryanto2121/laundryos/apps/api/internal/service"
	"github.com/nuryanto2121/laundryos/apps/api/internal/util"
)

// BotHandler handles incoming chat messages forwarded by WhatsApp Gateway.
type BotHandler struct {
	svc service.BotService
}

// NewBotHandler creates a new BotHandler.
func NewBotHandler(svc service.BotService) *BotHandler {
	return &BotHandler{svc: svc}
}

// HandleMessage handles POST /bot/messages
// @Summary Process Incoming WhatsApp Message
// @Description Matches incoming message text with FAQ knowledge base and responds with formatted text
// @Tags Bot
// @Accept json
// @Produce json
// @Param request body model.IncomingMessageRequest true "Incoming WhatsApp message payload"
// @Success 200 {object} model.BotResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /bot/messages [post]
func (h *BotHandler) HandleMessage(w http.ResponseWriter, r *http.Request) {
	var req model.IncomingMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := util.ValidateStruct(req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.svc.ProcessMessage(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to process message")
		return
	}

	writeJSON(w, http.StatusOK, res)
}
