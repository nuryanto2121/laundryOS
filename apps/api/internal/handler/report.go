package handler

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/nuryanto2121/laundryos/apps/api/internal/service"
)

// ReportHandler handles HTTP endpoints for reporting, charts, and CRM analytics.
type ReportHandler struct {
	svc *service.ReportService
}

// NewReportHandler creates a new ReportHandler.
func NewReportHandler(svc *service.ReportService) *ReportHandler {
	return &ReportHandler{svc: svc}
}

// GetHomeMetrics handles GET /reports/home-metrics.
func (h *ReportHandler) GetHomeMetrics(w http.ResponseWriter, r *http.Request) {
	metrics, err := h.svc.GetHomeMetrics(r.Context())
	if err != nil {
		slog.Error("report: get home metrics failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal memuat metrik beranda")
		return
	}
	writeJSON(w, http.StatusOK, metrics)
}

// GetFinancial handles GET /reports/financial.
func (h *ReportHandler) GetFinancial(w http.ResponseWriter, r *http.Request) {
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

	report, err := h.svc.GetFinancialReport(r.Context(), start, end)
	if err != nil {
		slog.Error("report: get financial failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal memuat laporan keuangan")
		return
	}

	writeJSON(w, http.StatusOK, report)
}

// ListUnpaid handles GET /reports/unpaid.
func (h *ReportHandler) ListUnpaid(w http.ResponseWriter, r *http.Request) {
	unpaid, err := h.svc.ListUnpaidOrders(r.Context())
	if err != nil {
		slog.Error("report: list unpaid failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal memuat laporan piutang")
		return
	}
	writeJSON(w, http.StatusOK, unpaid)
}

// GetServices handles GET /reports/services.
func (h *ReportHandler) GetServices(w http.ResponseWriter, r *http.Request) {
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

	services, err := h.svc.GetServiceVolumeReport(r.Context(), start, end)
	if err != nil {
		slog.Error("report: get services failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal memuat rekap volume layanan")
		return
	}

	writeJSON(w, http.StatusOK, services)
}

// GetTopCustomers handles GET /reports/top-customers.
func (h *ReportHandler) GetTopCustomers(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	customers, err := h.svc.GetTopCustomers(r.Context(), int32(limit))
	if err != nil {
		slog.Error("report: get top customers failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal memuat peringkat pelanggan")
		return
	}

	writeJSON(w, http.StatusOK, customers)
}

// GetCustomerCRM handles GET /reports/customers/{id}/crm.
func (h *ReportHandler) GetCustomerCRM(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "customer id is required")
		return
	}

	profile, err := h.svc.GetCustomerCRMProfile(r.Context(), id)
	if err != nil {
		slog.Error("report: get customer crm failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal memuat profil CRM pelanggan")
		return
	}

	writeJSON(w, http.StatusOK, profile)
}
