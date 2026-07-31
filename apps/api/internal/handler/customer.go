package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nuryanto2121/laundryos/apps/api/internal/model"
	"github.com/nuryanto2121/laundryos/apps/api/internal/service"
)

// CustomerHandler exposes HTTP endpoints for the customer resource.
type CustomerHandler struct {
	svc *service.CustomerService
}

// NewCustomerHandler creates a new CustomerHandler.
func NewCustomerHandler(svc *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{svc: svc}
}

// Create handles POST /customers
func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	customer, err := h.svc.Create(r.Context(), req)
	if err != nil {
		slog.Error("create customer failed", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to create customer")
		return
	}

	writeJSON(w, http.StatusCreated, customer)
}

// Get handles GET /customers/{id}
func (h *CustomerHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	customer, err := h.svc.Get(r.Context(), id)
	if err != nil {
		slog.Error("get customer failed", "id", id, "error", err)
		writeError(w, http.StatusNotFound, "customer not found")
		return
	}

	writeJSON(w, http.StatusOK, customer)
}

// List handles GET /customers
func (h *CustomerHandler) List(w http.ResponseWriter, r *http.Request) {
	customers, err := h.svc.List(r.Context())
	if err != nil {
		slog.Error("list customers failed", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to list customers")
		return
	}

	writeJSON(w, http.StatusOK, customers)
}

// Update handles PUT /customers/{id}
func (h *CustomerHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req model.UpdateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	customer, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		slog.Error("update customer failed", "id", id, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to update customer")
		return
	}

	writeJSON(w, http.StatusOK, customer)
}

// Delete handles DELETE /customers/{id}
func (h *CustomerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.svc.Delete(r.Context(), id); err != nil {
		slog.Error("delete customer failed", "id", id, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to delete customer")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// writeJSON encodes v as JSON and writes it with the given status code.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("failed to encode response", "error", err)
	}
}

// writeError writes a JSON error response.
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
