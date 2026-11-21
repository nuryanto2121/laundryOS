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
)

// AttendanceHandler handles HTTP endpoints for shifts, attendance, and employee leaves.
type AttendanceHandler struct {
	svc service.AttendanceService
}

// NewAttendanceHandler creates a new AttendanceHandler.
func NewAttendanceHandler(svc service.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{svc: svc}
}

// -----------------------------------------------------------------------------
// Shifts
// -----------------------------------------------------------------------------

// ListShifts handles GET /master/shifts
func (h *AttendanceHandler) ListShifts(w http.ResponseWriter, r *http.Request) {
	shifts, err := h.svc.ListShifts(r.Context())
	if err != nil {
		slog.Error("attendance: list shifts failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal memuat daftar shift")
		return
	}
	writeJSON(w, http.StatusOK, shifts)
}

// UpdateShift handles PUT /master/shifts/{id}
func (h *AttendanceHandler) UpdateShift(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "shift id is required")
		return
	}

	var req model.UpdateShiftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	shift, err := h.svc.UpdateShift(r.Context(), id, req)
	if err != nil {
		slog.Error("attendance: update shift failed", "error", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, shift)
}

// -----------------------------------------------------------------------------
// Attendances (Clock-in, Clock-out, History)
// -----------------------------------------------------------------------------

// CheckIn handles POST /attendances/check-in
func (h *AttendanceHandler) CheckIn(w http.ResponseWriter, r *http.Request) {
	var req model.CheckInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	att, err := h.svc.CheckIn(r.Context(), req)
	if err != nil {
		slog.Warn("attendance: check-in failed", "error", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, att)
}

// CheckOut handles POST /attendances/check-out
func (h *AttendanceHandler) CheckOut(w http.ResponseWriter, r *http.Request) {
	var req model.CheckOutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	att, err := h.svc.CheckOut(r.Context(), req)
	if err != nil {
		slog.Warn("attendance: check-out failed", "error", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, att)
}

// GetActiveAttendance handles GET /attendances/today?employee=...
func (h *AttendanceHandler) GetActiveAttendance(w http.ResponseWriter, r *http.Request) {
	emp := strings.TrimSpace(r.URL.Query().Get("employee"))
	if emp == "" {
		writeError(w, http.StatusBadRequest, "parameter employee wajib diisi")
		return
	}

	att, err := h.svc.GetActiveAttendance(r.Context(), emp)
	if err != nil {
		slog.Error("attendance: get active failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal memuat status shift")
		return
	}

	if att.ID == "" {
		writeJSON(w, http.StatusOK, nil)
		return
	}
	writeJSON(w, http.StatusOK, att)
}

// ListActiveAttendances handles GET /attendances/active
func (h *AttendanceHandler) ListActiveAttendances(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.ListActiveAttendances(r.Context())
	if err != nil {
		slog.Error("attendance: list active failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal memuat staf bertugas")
		return
	}
	writeJSON(w, http.StatusOK, list)
}

// ListAttendances handles GET /attendances/history
func (h *AttendanceHandler) ListAttendances(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("start")
	endDate := r.URL.Query().Get("end")
	employee := r.URL.Query().Get("employee")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if limit <= 0 {
		limit = 15
	}
	if page <= 0 {
		page = 1
	}

	list, total, err := h.svc.ListAttendances(r.Context(), startDate, endDate, employee, page, limit)
	if err != nil {
		slog.Error("attendance: list history failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal memuat riwayat absensi")
		return
	}

	totalPages := (total + limit - 1) / limit
	res := map[string]interface{}{
		"data": list,
		"meta": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total_items": total,
			"total_pages": totalPages,
		},
	}
	writeJSON(w, http.StatusOK, res)
}

// -----------------------------------------------------------------------------
// Leave Requests (Pengajuan Libur: Max 2x Sebulan)
// -----------------------------------------------------------------------------

// GetLeaveQuota handles GET /leaves/quota?employee=...&date=...
func (h *AttendanceHandler) GetLeaveQuota(w http.ResponseWriter, r *http.Request) {
	emp := strings.TrimSpace(r.URL.Query().Get("employee"))
	if emp == "" {
		writeError(w, http.StatusBadRequest, "parameter employee wajib diisi")
		return
	}
	dateStr := r.URL.Query().Get("date")

	quota, err := h.svc.GetLeaveQuota(r.Context(), emp, dateStr)
	if err != nil {
		slog.Error("attendance: get quota failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal menghitung sisa kuota libur")
		return
	}

	writeJSON(w, http.StatusOK, quota)
}

// CreateLeaveRequest handles POST /leaves/request
func (h *AttendanceHandler) CreateLeaveRequest(w http.ResponseWriter, r *http.Request) {
	var req model.CreateLeaveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	created, err := h.svc.CreateLeaveRequest(r.Context(), req)
	if err != nil {
		slog.Warn("attendance: create leave failed", "error", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

// ListLeaveRequests handles GET /leaves
func (h *AttendanceHandler) ListLeaveRequests(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	employee := r.URL.Query().Get("employee")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if limit <= 0 {
		limit = 15
	}
	if page <= 0 {
		page = 1
	}

	list, total, err := h.svc.ListLeaveRequests(r.Context(), status, employee, page, limit)
	if err != nil {
		slog.Error("attendance: list leaves failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal memuat daftar pengajuan libur")
		return
	}

	totalPages := (total + limit - 1) / limit
	res := map[string]interface{}{
		"data": list,
		"meta": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total_items": total,
			"total_pages": totalPages,
		},
	}
	writeJSON(w, http.StatusOK, res)
}

// ReviewLeaveRequest handles PUT /leaves/{id}/review
func (h *AttendanceHandler) ReviewLeaveRequest(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "leave id is required")
		return
	}

	var req model.ReviewLeaveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	reviewed, err := h.svc.ReviewLeaveRequest(r.Context(), id, req)
	if err != nil {
		slog.Warn("attendance: review leave failed", "error", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, reviewed)
}
