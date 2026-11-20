package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/nuryanto2121/laundryos/apps/api/internal/model"
	"github.com/nuryanto2121/laundryos/apps/api/internal/service"
	"github.com/nuryanto2121/laundryos/apps/api/internal/util"
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
// @Summary Create Customer
// @Description Register a new customer or reactivate an existing inactive customer
// @Tags Customers
// @Accept json
// @Produce json
// @Param request body model.CreateCustomerRequest true "Customer creation payload"
// @Success 201 {object} model.Customer
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /customers [post]
func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := util.ValidateStruct(req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	customer, err := h.svc.Create(r.Context(), req)
	if err != nil {
		slog.Error("create customer failed", "error", err)
		if err.Error() == "customer already registered" || err.Error() == "customer is blocked" {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to create customer")
		return
	}

	writeJSON(w, http.StatusCreated, customer)
}

// Get handles GET /customers/{id}
// @Summary Get Customer Detail
// @Description Retrieve a customer by their unique UUID
// @Tags Customers
// @Produce json
// @Param id path string true "Customer UUID"
// @Success 200 {object} model.Customer
// @Failure 404 {object} model.ErrorResponse
// @Router /customers/{id} [get]
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

// List handles GET /customers (supporting optional ?page=1&limit=10&search=...).
func (h *CustomerHandler) List(w http.ResponseWriter, r *http.Request) {
	search := strings.TrimSpace(r.URL.Query().Get("search"))
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	if pageStr != "" || limitStr != "" {
		page, _ := strconv.Atoi(pageStr)
		limit, _ := strconv.Atoi(limitStr)
		paginated, err := h.svc.ListPaginated(r.Context(), search, page, limit)
		if err != nil {
			slog.Error("list customers paginated failed", "error", err)
			writeError(w, http.StatusInternalServerError, "failed to list customers")
			return
		}
		writeJSON(w, http.StatusOK, paginated)
		return
	}

	customers, err := h.svc.List(r.Context())
	if err != nil {
		slog.Error("list customers failed", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to list customers")
		return
	}

	if search != "" {
		searchLower := strings.ToLower(search)
		filtered := make([]model.Customer, 0)
		for _, c := range customers {
			if strings.Contains(strings.ToLower(c.FullName), searchLower) ||
				strings.Contains(c.Phone, search) ||
				strings.Contains(strings.ToLower(c.CustomerCode), searchLower) ||
				strings.Contains(strings.ToLower(c.Address), searchLower) {
				filtered = append(filtered, c)
			}
		}
		customers = filtered
	}

	writeJSON(w, http.StatusOK, customers)
}

// Update handles PUT /customers/{id}
// @Summary Update Customer
// @Description Update mutable customer details
// @Tags Customers
// @Accept json
// @Produce json
// @Param id path string true "Customer UUID"
// @Param request body model.UpdateCustomerRequest true "Customer update payload"
// @Success 200 {object} model.Customer
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /customers/{id} [put]
func (h *CustomerHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req model.UpdateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := util.ValidateStruct(req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
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

// UpdateStatus handles PATCH /customers/{id}/status
// @Summary Update Customer Status
// @Description Transition customer lifecycle status (ACTIVE, INACTIVE, BLOCKED)
// @Tags Customers
// @Accept json
// @Produce json
// @Param id path string true "Customer UUID"
// @Param request body model.UpdateCustomerStatusRequest true "Status update payload"
// @Success 200 {object} model.Customer
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /customers/{id}/status [patch]
func (h *CustomerHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req model.UpdateCustomerStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := util.ValidateStruct(req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	customer, err := h.svc.UpdateStatus(r.Context(), id, req.Status)
	if err != nil {
		slog.Error("update customer status failed", "id", id, "error", err)
		if err.Error() == "invalid target status" || strings.HasPrefix(err.Error(), "customer already") {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update customer status")
		return
	}

	writeJSON(w, http.StatusOK, customer)
}

// Delete handles DELETE /customers/{id}
// @Summary Delete Customer
// @Description Soft-delete a customer by setting deleted_at timestamp
// @Tags Customers
// @Param id path string true "Customer UUID"
// @Success 204 "No Content"
// @Failure 500 {object} model.ErrorResponse
// @Router /customers/{id} [delete]
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
	writeJSON(w, status, model.ErrorResponse{Error: msg})
}
