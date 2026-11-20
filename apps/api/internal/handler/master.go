package handler

import (
	"log/slog"
	"net/http"

	"github.com/nuryanto2121/laundryos/apps/api/internal/service"
)

// MasterHandler handles requests for lookup master tables.
type MasterHandler struct {
	svc *service.MasterService
}

// NewMasterHandler creates a new MasterHandler.
func NewMasterHandler(svc *service.MasterService) *MasterHandler {
	return &MasterHandler{svc: svc}
}

// ListPerfumes handles GET /master/perfumes.
func (h *MasterHandler) ListPerfumes(w http.ResponseWriter, r *http.Request) {
	perfumes, err := h.svc.ListPerfumes(r.Context())
	if err != nil {
		slog.Error("master: list perfumes failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal memuat data parfum")
		return
	}
	writeJSON(w, http.StatusOK, perfumes)
}

// ListDiscounts handles GET /master/discounts.
func (h *MasterHandler) ListDiscounts(w http.ResponseWriter, r *http.Request) {
	discounts, err := h.svc.ListDiscounts(r.Context())
	if err != nil {
		slog.Error("master: list discounts failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal memuat data diskon")
		return
	}
	writeJSON(w, http.StatusOK, discounts)
}

// ListPaymentMethods handles GET /master/payment-methods.
func (h *MasterHandler) ListPaymentMethods(w http.ResponseWriter, r *http.Request) {
	methods, err := h.svc.ListPaymentMethods(r.Context())
	if err != nil {
		slog.Error("master: list payment methods failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal memuat metode pembayaran")
		return
	}
	writeJSON(w, http.StatusOK, methods)
}

// ListOutlets handles GET /master/outlets.
func (h *MasterHandler) ListOutlets(w http.ResponseWriter, r *http.Request) {
	outlets, err := h.svc.ListOutlets(r.Context())
	if err != nil {
		slog.Error("master: list outlets failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal memuat data outlet")
		return
	}
	writeJSON(w, http.StatusOK, outlets)
}

// ListServices handles GET /master/services.
func (h *MasterHandler) ListServices(w http.ResponseWriter, r *http.Request) {
	services, err := h.svc.ListServices(r.Context())
	if err != nil {
		slog.Error("master: list services failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal memuat data layanan")
		return
	}
	writeJSON(w, http.StatusOK, services)
}
