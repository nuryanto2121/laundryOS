package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
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

type createPerfumeRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

// CreatePerfume handles POST /master/perfumes.
func (h *MasterHandler) CreatePerfume(w http.ResponseWriter, r *http.Request) {
	var req createPerfumeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "nama parfum wajib diisi")
		return
	}
	perfume, err := h.svc.CreatePerfume(r.Context(), req.Name, req.Description, true)
	if err != nil {
		slog.Error("master: create perfume failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal membuat parfum")
		return
	}
	writeJSON(w, http.StatusCreated, perfume)
}

// UpdatePerfume handles PUT /master/perfumes/{id}.
func (h *MasterHandler) UpdatePerfume(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req createPerfumeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	perfume, err := h.svc.UpdatePerfume(r.Context(), id, req.Name, req.Description, req.IsActive)
	if err != nil {
		slog.Error("master: update perfume failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal memperbarui parfum")
		return
	}
	writeJSON(w, http.StatusOK, perfume)
}

// DeletePerfume handles DELETE /master/perfumes/{id}.
func (h *MasterHandler) DeletePerfume(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.DeletePerfume(r.Context(), id); err != nil {
		slog.Error("master: delete perfume failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal menghapus parfum")
		return
	}
	w.WriteHeader(http.StatusNoContent)
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

type createDiscountRequest struct {
	Name            string  `json:"name"`
	DiscountType    string  `json:"discount_type"` // 'PERCENTAGE', 'NOMINAL'
	Value           float64 `json:"value"`
	DiscountPercent float64 `json:"discount_percent"`
	FixedAmount     float64 `json:"fixed_amount"`
	MinOrderAmount  float64 `json:"min_order_amount"`
	IsActive        bool    `json:"is_active"`
}

// CreateDiscount handles POST /master/discounts.
func (h *MasterHandler) CreateDiscount(w http.ResponseWriter, r *http.Request) {
	var req createDiscountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Value == 0 {
		if req.DiscountPercent > 0 {
			req.DiscountType = "PERCENTAGE"
			req.Value = req.DiscountPercent
		} else if req.FixedAmount > 0 {
			req.DiscountType = "NOMINAL"
			req.Value = req.FixedAmount
		}
	}
	if req.DiscountType == "" {
		req.DiscountType = "NOMINAL"
	}
	if req.Name == "" || req.Value <= 0 {
		writeError(w, http.StatusBadRequest, "nama dan nilai diskon (>0) wajib diisi")
		return
	}
	discount, err := h.svc.CreateDiscount(r.Context(), req.Name, req.DiscountType, req.Value, req.MinOrderAmount, true)
	if err != nil {
		slog.Error("master: create discount failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal membuat diskon")
		return
	}
	writeJSON(w, http.StatusCreated, discount)
}

// UpdateDiscount handles PUT /master/discounts/{id}.
func (h *MasterHandler) UpdateDiscount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req createDiscountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Value == 0 {
		if req.DiscountPercent > 0 {
			req.DiscountType = "PERCENTAGE"
			req.Value = req.DiscountPercent
		} else if req.FixedAmount > 0 {
			req.DiscountType = "NOMINAL"
			req.Value = req.FixedAmount
		}
	}
	if req.DiscountType == "" {
		req.DiscountType = "NOMINAL"
	}
	discount, err := h.svc.UpdateDiscount(r.Context(), id, req.Name, req.DiscountType, req.Value, req.MinOrderAmount, req.IsActive)
	if err != nil {
		slog.Error("master: update discount failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal memperbarui diskon")
		return
	}
	writeJSON(w, http.StatusOK, discount)
}

// DeleteDiscount handles DELETE /master/discounts/{id}.
func (h *MasterHandler) DeleteDiscount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.DeleteDiscount(r.Context(), id); err != nil {
		slog.Error("master: delete discount failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal menghapus diskon")
		return
	}
	w.WriteHeader(http.StatusNoContent)
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
	sortBy := strings.TrimSpace(r.URL.Query().Get("sort_by"))
	sortDir := strings.TrimSpace(r.URL.Query().Get("sort_dir"))
	services, err := h.svc.ListServices(r.Context(), sortBy, sortDir)
	if err != nil {
		slog.Error("master: list services failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal memuat data layanan")
		return
	}
	writeJSON(w, http.StatusOK, services)
}

type createServiceRequest struct {
	ServiceCode  string  `json:"service_code"`
	ServiceName  string  `json:"service_name"`
	DurationDay  int     `json:"duration_day"`
	Price        float64 `json:"price"`
	Active       bool    `json:"active"`
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Category     string  `json:"category"`
	Unit         string  `json:"unit"`
}

// CreateService handles POST /master/services.
func (h *MasterHandler) CreateService(w http.ResponseWriter, r *http.Request) {
	var req createServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.ServiceCode == "" || req.ServiceName == "" || req.Price <= 0 {
		writeError(w, http.StatusBadRequest, "kode, nama layanan, dan harga (>0) wajib diisi")
		return
	}
	catName := req.CategoryName
	if catName == "" {
		catName = req.Category
	}
	unit := req.Unit
	if unit == "" {
		unit = "Kg"
	}
	srv, err := h.svc.CreateService(r.Context(), req.ServiceCode, req.ServiceName, req.DurationDay, req.Price, true, req.CategoryID, catName, unit)
	if err != nil {
		slog.Error("master: create service failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal membuat layanan")
		return
	}
	writeJSON(w, http.StatusCreated, srv)
}

// UpdateService handles PUT /master/services/{id}.
func (h *MasterHandler) UpdateService(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req createServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	catName := req.CategoryName
	if catName == "" {
		catName = req.Category
	}
	unit := req.Unit
	if unit == "" {
		unit = "Kg"
	}
	srv, err := h.svc.UpdateService(r.Context(), id, req.ServiceName, req.DurationDay, req.Price, req.Active, req.CategoryID, catName, unit)
	if err != nil {
		slog.Error("master: update service failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal memperbarui layanan")
		return
	}
	writeJSON(w, http.StatusOK, srv)
}

// DeleteService handles DELETE /master/services/{id}.
func (h *MasterHandler) DeleteService(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.DeleteService(r.Context(), id); err != nil {
		slog.Error("master: delete service failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal menghapus layanan")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListServiceCategories handles GET /master/service-categories.
func (h *MasterHandler) ListServiceCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.svc.ListAllServiceCategories(r.Context())
	if err != nil {
		slog.Error("master: list service categories failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal memuat data kategori layanan")
		return
	}
	writeJSON(w, http.StatusOK, categories)
}

type createServiceCategoryRequest struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	DefaultUnit  string `json:"default_unit"`
	DisplayOrder int    `json:"display_order"`
	IsActive     bool   `json:"is_active"`
}

// CreateServiceCategory handles POST /master/service-categories.
func (h *MasterHandler) CreateServiceCategory(w http.ResponseWriter, r *http.Request) {
	var req createServiceCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Code == "" || req.Name == "" {
		writeError(w, http.StatusBadRequest, "kode dan nama kategori wajib diisi")
		return
	}
	if req.DefaultUnit == "" {
		req.DefaultUnit = "Kg"
	}
	cat, err := h.svc.CreateServiceCategory(r.Context(), req.Code, req.Name, req.DefaultUnit, req.DisplayOrder, true)
	if err != nil {
		slog.Error("master: create service category failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal membuat kategori layanan")
		return
	}
	writeJSON(w, http.StatusCreated, cat)
}

// UpdateServiceCategory handles PUT /master/service-categories/{id}.
func (h *MasterHandler) UpdateServiceCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req createServiceCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	cat, err := h.svc.UpdateServiceCategory(r.Context(), id, req.Name, req.DefaultUnit, req.DisplayOrder, req.IsActive)
	if err != nil {
		slog.Error("master: update service category failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal memperbarui kategori layanan")
		return
	}
	writeJSON(w, http.StatusOK, cat)
}

// DeleteServiceCategory handles DELETE /master/service-categories/{id}.
func (h *MasterHandler) DeleteServiceCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.DeleteServiceCategory(r.Context(), id); err != nil {
		slog.Error("master: delete service category failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal menghapus kategori layanan")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListExpenseCategories handles GET /master/expense-categories.
func (h *MasterHandler) ListExpenseCategories(w http.ResponseWriter, r *http.Request) {
	sortBy := strings.TrimSpace(r.URL.Query().Get("sort_by"))
	sortDir := strings.TrimSpace(r.URL.Query().Get("sort_dir"))
	cats, err := h.svc.ListAllExpenseCategories(r.Context(), sortBy, sortDir)
	if err != nil {
		slog.Error("master: list expense categories failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal memuat kategori pengeluaran")
		return
	}
	writeJSON(w, http.StatusOK, cats)
}

type createExpenseCategoryRequest struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

// CreateExpenseCategory handles POST /master/expense-categories.
func (h *MasterHandler) CreateExpenseCategory(w http.ResponseWriter, r *http.Request) {
	var req createExpenseCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Code == "" || req.Name == "" {
		writeError(w, http.StatusBadRequest, "kode dan nama kategori wajib diisi")
		return
	}
	cat, err := h.svc.CreateExpenseCategory(r.Context(), req.Code, req.Name, true)
	if err != nil {
		slog.Error("master: create expense category failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal membuat kategori pengeluaran")
		return
	}
	writeJSON(w, http.StatusCreated, cat)
}

// UpdateExpenseCategory handles PUT /master/expense-categories/{id}.
func (h *MasterHandler) UpdateExpenseCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req createExpenseCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	cat, err := h.svc.UpdateExpenseCategory(r.Context(), id, req.Name, req.IsActive)
	if err != nil {
		slog.Error("master: update expense category failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal memperbarui kategori")
		return
	}
	writeJSON(w, http.StatusOK, cat)
}

// DeleteExpenseCategory handles DELETE /master/expense-categories/{id}.
func (h *MasterHandler) DeleteExpenseCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.DeleteExpenseCategory(r.Context(), id); err != nil {
		slog.Error("master: delete expense category failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal menghapus kategori")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListEmployees handles GET /master/employees.
func (h *MasterHandler) ListEmployees(w http.ResponseWriter, r *http.Request) {
	emps, err := h.svc.ListEmployees(r.Context())
	if err != nil {
		slog.Error("master: list employees failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal memuat data pegawai")
		return
	}
	writeJSON(w, http.StatusOK, emps)
}

type createEmployeeRequest struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
	Pin      string `json:"pin"`
	IsActive bool   `json:"is_active"`
}

// CreateEmployee handles POST /master/employees.
func (h *MasterHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var req createEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Username == "" || req.FullName == "" {
		writeError(w, http.StatusBadRequest, "username dan nama lengkap wajib diisi")
		return
	}
	if req.Role == "" {
		req.Role = "KASIR"
	}
	if req.Pin == "" {
		req.Pin = "1234"
	}
	emp, err := h.svc.CreateEmployee(r.Context(), req.Username, req.FullName, req.Phone, req.Role, req.Pin, true)
	if err != nil {
		slog.Error("master: create employee failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal membuat pegawai")
		return
	}
	writeJSON(w, http.StatusCreated, emp)
}

// UpdateEmployee handles PUT /master/employees/{id}.
func (h *MasterHandler) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req createEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	emp, err := h.svc.UpdateEmployee(r.Context(), id, req.FullName, req.Phone, req.Role, req.Pin, req.IsActive)
	if err != nil {
		slog.Error("master: update employee failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal memperbarui pegawai")
		return
	}
	writeJSON(w, http.StatusOK, emp)
}

// DeleteEmployee handles DELETE /master/employees/{id}.
func (h *MasterHandler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.DeleteEmployee(r.Context(), id); err != nil {
		slog.Error("master: delete employee failed", "error", err)
		writeError(w, http.StatusInternalServerError, "gagal menghapus pegawai")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
