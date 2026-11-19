package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/nuryanto2121/laundryos/apps/api/internal/model"
	"github.com/nuryanto2121/laundryos/apps/api/internal/service"
)

// ExpenseHandler handles HTTP endpoints for operational expenses.
type ExpenseHandler struct {
	svc *service.ExpenseService
}

// NewExpenseHandler creates a new ExpenseHandler.
func NewExpenseHandler(svc *service.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{svc: svc}
}

// Create handles POST /expenses.
func (h *ExpenseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Category == "" || req.Amount <= 0 || req.Description == "" {
		writeError(w, http.StatusBadRequest, "kategori, nominal (>0), dan keterangan wajib diisi")
		return
	}

	expense, err := h.svc.Create(r.Context(), req)
	if err != nil {
		slog.Error("expense: create failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal menyimpan pengeluaran")
		return
	}

	writeJSON(w, http.StatusCreated, expense)
}

// List handles GET /expenses.
func (h *ExpenseHandler) List(w http.ResponseWriter, r *http.Request) {
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	start := time.Now().AddDate(0, 0, -30)
	end := time.Now()

	if t, err := time.Parse("2006-01-02", startStr); err == nil {
		start = t
	}
	if t, err := time.Parse("2006-01-02", endStr); err == nil {
		end = t
	}

	expenses, err := h.svc.List(r.Context(), start, end)
	if err != nil {
		slog.Error("expense: list failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal memuat daftar pengeluaran")
		return
	}

	writeJSON(w, http.StatusOK, expenses)
}

// GetSummary handles GET /expenses/summary.
func (h *ExpenseHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	start := time.Now().AddDate(0, 0, -30)
	end := time.Now()

	if t, err := time.Parse("2006-01-02", startStr); err == nil {
		start = t
	}
	if t, err := time.Parse("2006-01-02", endStr); err == nil {
		end = t
	}

	summary, err := h.svc.GetSummary(r.Context(), start, end)
	if err != nil {
		slog.Error("expense: summary failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal memuat rekap pengeluaran")
		return
	}

	writeJSON(w, http.StatusOK, summary)
}

// Delete handles DELETE /expenses/{id}.
func (h *ExpenseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "expense id is required")
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		slog.Error("expense: delete failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal menghapus pengeluaran")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "pengeluaran berhasil dihapus"})
}
