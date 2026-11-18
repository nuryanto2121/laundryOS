package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nuryanto2121/laundryos/apps/api/internal/model"
	"github.com/nuryanto2121/laundryos/apps/api/internal/service"
)

// OrderHandler exposes HTTP endpoints for the order resource.
type OrderHandler struct {
	svc *service.OrderService
}

// NewOrderHandler creates a new OrderHandler.
func NewOrderHandler(svc *service.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

// Create handles POST /orders.
func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	order, err := h.svc.Create(r.Context(), req)
	if err != nil {
		slog.Error("create order failed", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to create order")
		return
	}
	writeJSON(w, http.StatusCreated, order)
}

// Get handles GET /orders/{id}.
func (h *OrderHandler) Get(w http.ResponseWriter, r *http.Request) {
	order, err := h.svc.Get(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		slog.Error("get order failed", "error", err)
		writeError(w, http.StatusNotFound, "order not found")
		return
	}
	writeJSON(w, http.StatusOK, order)
}

// List handles GET /orders.
func (h *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
	orders, err := h.svc.List(r.Context())
	if err != nil {
		slog.Error("list orders failed", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to list orders")
		return
	}
	writeJSON(w, http.StatusOK, orders)
}

// Update handles PUT /orders/{id}.
func (h *OrderHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req model.UpdateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	order, err := h.svc.Update(r.Context(), chi.URLParam(r, "id"), req)
	if err != nil {
		slog.Error("update order failed", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to update order")
		return
	}
	writeJSON(w, http.StatusOK, order)
}

// Delete handles DELETE /orders/{id}.
func (h *OrderHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.Delete(r.Context(), chi.URLParam(r, "id")); err != nil {
		slog.Error("delete order failed", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to delete order")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
