package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nuryanto2121/laundryos/apps/api/internal/model"
	"github.com/nuryanto2121/laundryos/apps/api/internal/service"
	"github.com/nuryanto2121/laundryos/apps/api/internal/util"
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
// @Summary Create Order
// @Description Place a new laundry order with auto-generated order sequence
// @Tags Orders
// @Accept json
// @Produce json
// @Param request body model.CreateOrderRequest true "Order creation payload"
// @Success 201 {object} model.Order
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /orders [post]
func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := util.ValidateStruct(req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
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
// @Summary Get Order Detail
// @Description Retrieve order details by order UUID
// @Tags Orders
// @Produce json
// @Param id path string true "Order UUID"
// @Success 200 {object} model.Order
// @Failure 404 {object} model.ErrorResponse
// @Router /orders/{id} [get]
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
// @Summary List Orders
// @Description Retrieve a list of all laundry orders
// @Tags Orders
// @Produce json
// @Success 200 {array} model.Order
// @Failure 500 {object} model.ErrorResponse
// @Router /orders [get]
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
// @Summary Update Order
// @Description Update mutable order information (weight, price, notes)
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order UUID"
// @Param request body model.UpdateOrderRequest true "Order update payload"
// @Success 200 {object} model.Order
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /orders/{id} [put]
func (h *OrderHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req model.UpdateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := util.ValidateStruct(req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
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
// @Summary Delete Order
// @Description Remove an order record by ID
// @Tags Orders
// @Param id path string true "Order UUID"
// @Success 204 "No Content"
// @Failure 500 {object} model.ErrorResponse
// @Router /orders/{id} [delete]
func (h *OrderHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.Delete(r.Context(), chi.URLParam(r, "id")); err != nil {
		slog.Error("delete order failed", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to delete order")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
