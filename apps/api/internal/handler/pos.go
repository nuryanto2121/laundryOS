package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nuryanto2121/laundryos/apps/api/internal/model"
	"github.com/nuryanto2121/laundryos/apps/api/internal/service"
)

// POSHandler handles HTTP endpoints for Web POS cashier and order progression.
type POSHandler struct {
	svc *service.POSService
}

// NewPOSHandler creates a new POSHandler.
func NewPOSHandler(svc *service.POSService) *POSHandler {
	return &POSHandler{svc: svc}
}

// Create handles POST /pos/orders.
func (h *POSHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreatePOSOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	order, items, err := h.svc.CreateOrder(r.Context(), req)
	if err != nil {
		slog.Error("pos: create order failed", "error", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"order": order,
		"items": items,
	})
}

// List handles GET /pos/orders.
func (h *POSHandler) List(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	search := r.URL.Query().Get("search")

	orders, err := h.svc.ListOrders(r.Context(), status, search)
	if err != nil {
		slog.Error("pos: list orders failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal memuat daftar antrian order")
		return
	}

	writeJSON(w, http.StatusOK, orders)
}

// Get handles GET /pos/orders/{id}.
func (h *POSHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "order id is required")
		return
	}

	detail, err := h.svc.GetOrderDetail(r.Context(), id)
	if err != nil {
		slog.Error("pos: get order detail failed", "error", err)
		writeError(w, http.StatusNotFound, "order tidak ditemukan")
		return
	}

	writeJSON(w, http.StatusOK, detail)
}

// UpdatePayment handles PATCH /pos/orders/{id}/payment.
func (h *POSHandler) UpdatePayment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "order id is required")
		return
	}

	var req model.UpdatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	order, err := h.svc.UpdatePayment(r.Context(), id, req)
	if err != nil {
		slog.Error("pos: update payment failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal mencatat pembayaran")
		return
	}

	writeJSON(w, http.StatusOK, order)
}

// UpdateStatus handles PATCH /pos/orders/{id}/status.
func (h *POSHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "order id is required")
		return
	}

	var req model.UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Status == "" {
		writeError(w, http.StatusBadRequest, "status is required")
		return
	}

	order, err := h.svc.UpdateStatus(r.Context(), id, req)
	if err != nil {
		slog.Error("pos: update status failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal memperbarui status order")
		return
	}

	writeJSON(w, http.StatusOK, order)
}

// Cancel handles PATCH /pos/orders/{id}/cancel.
func (h *POSHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "order id is required")
		return
	}

	var req model.CancelOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	order, err := h.svc.CancelOrder(r.Context(), id, req)
	if err != nil {
		slog.Error("pos: cancel order failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal membatalkan order")
		return
	}

	writeJSON(w, http.StatusOK, order)
}

// UpdateRack handles PATCH /pos/orders/{id}/rack.
func (h *POSHandler) UpdateRack(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "order id is required")
		return
	}

	var req model.UpdateRackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	order, err := h.svc.UpdateRack(r.Context(), id, req.RackNumber)
	if err != nil {
		slog.Error("pos: update rack failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal memperbarui nomor rak")
		return
	}

	writeJSON(w, http.StatusOK, order)
}

// AddPhoto handles POST /pos/orders/{id}/photos.
func (h *POSHandler) AddPhoto(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "order id is required")
		return
	}

	var req model.AddOrderPhotoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.PhotoType == "" || req.PhotoURL == "" {
		writeError(w, http.StatusBadRequest, "photo_type and photo_url are required")
		return
	}

	photo, err := h.svc.AddPhoto(r.Context(), id, req)
	if err != nil {
		slog.Error("pos: add photo failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal menyimpan foto dokumentasi")
		return
	}

	writeJSON(w, http.StatusCreated, photo)
}

// ResendInvoice handles POST /pos/orders/{id}/resend-invoice.
func (h *POSHandler) ResendInvoice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "order id is required")
		return
	}

	if err := h.svc.ResendCustomerInvoice(r.Context(), id); err != nil {
		slog.Error("pos: resend invoice failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal mengirim ulang invoice WA")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "invoice WhatsApp berhasil dikirim ke pelanggan"})
}
