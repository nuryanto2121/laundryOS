package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/nuryanto2121/laundryos/apps/api/internal/model"
	"github.com/nuryanto2121/laundryos/apps/api/internal/service"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	svc *service.MasterService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(svc *service.MasterService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Login handles POST /auth/login.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "format body request tidak valid")
		return
	}

	if req.Username == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "username dan password wajib diisi")
		return
	}

	user, err := h.svc.Authenticate(r.Context(), req.Username, req.Password)
	if err != nil {
		slog.Warn("auth: login failed", "username", req.Username, "error", err)
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}

	token := "los_token_" + user.Role + "_" + user.ID

	writeJSON(w, http.StatusOK, model.LoginResponse{
		Token: token,
		User:  user,
	})
}
