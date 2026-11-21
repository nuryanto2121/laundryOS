package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/nuryanto2121/laundryos/apps/api/internal/model"
	"github.com/nuryanto2121/laundryos/apps/api/internal/repository"
	"github.com/nuryanto2121/laundryos/apps/api/internal/util"
)

// AttendanceService defines business operations for shifts, attendances, and leaves.
type AttendanceService interface {
	ListShifts(ctx context.Context) ([]model.Shift, error)
	UpdateShift(ctx context.Context, id string, req model.UpdateShiftRequest) (model.Shift, error)

	CheckIn(ctx context.Context, req model.CheckInRequest) (model.EmployeeAttendance, error)
	CheckOut(ctx context.Context, req model.CheckOutRequest) (model.EmployeeAttendance, error)
	GetActiveAttendance(ctx context.Context, employeeName string) (model.EmployeeAttendance, error)
	ListActiveAttendances(ctx context.Context) ([]model.EmployeeAttendance, error)
	ListAttendances(ctx context.Context, startDate, endDate, employeeName string, page, limit int) ([]model.EmployeeAttendance, int, error)

	GetLeaveQuota(ctx context.Context, employeeName, dateStr string) (model.LeaveQuotaResponse, error)
	CreateLeaveRequest(ctx context.Context, req model.CreateLeaveRequest) (model.LeaveRequest, error)
	ListLeaveRequests(ctx context.Context, status, employeeName string, page, limit int) ([]model.LeaveRequest, int, error)
	ReviewLeaveRequest(ctx context.Context, id string, req model.ReviewLeaveRequest) (model.LeaveRequest, error)
}

type attendanceService struct {
	repo *repository.AttendanceRepository
}

// NewAttendanceService creates a new AttendanceService.
func NewAttendanceService(repo *repository.AttendanceRepository) AttendanceService {
	return &attendanceService{repo: repo}
}

// -----------------------------------------------------------------------------
// Shift Management
// -----------------------------------------------------------------------------

func (s *attendanceService) ListShifts(ctx context.Context) ([]model.Shift, error) {
	return s.repo.ListShifts(ctx)
}

func (s *attendanceService) UpdateShift(ctx context.Context, id string, req model.UpdateShiftRequest) (model.Shift, error) {
	if strings.TrimSpace(req.ShiftName) == "" {
		return model.Shift{}, errors.New("shift_name is required")
	}
	if strings.TrimSpace(req.StartTime) == "" || strings.TrimSpace(req.EndTime) == "" {
		return model.Shift{}, errors.New("start_time and end_time are required")
	}
	return s.repo.UpdateShift(ctx, id, req)
}

// -----------------------------------------------------------------------------
// Check-In (Clock-In)
// -----------------------------------------------------------------------------

func (s *attendanceService) CheckIn(ctx context.Context, req model.CheckInRequest) (model.EmployeeAttendance, error) {
	empName := strings.TrimSpace(req.EmployeeName)
	if empName == "" {
		return model.EmployeeAttendance{}, errors.New("nama karyawan wajib diisi")
	}
	if strings.TrimSpace(req.ShiftID) == "" {
		return model.EmployeeAttendance{}, errors.New("pilihan shift wajib diisi")
	}
	if strings.TrimSpace(req.PhotoBase64) == "" {
		return model.EmployeeAttendance{}, errors.New("foto selfie di tempat laundry wajib diambil")
	}

	// 1. Prevent double check-in
	existing, err := s.repo.GetActiveAttendanceByEmployee(ctx, empName)
	if err == nil && existing.ID != "" {
		return model.EmployeeAttendance{}, fmt.Errorf("karyawan %s saat ini sudah dalam status Sedang Bertugas (Shift: %s)", empName, existing.ShiftName)
	}

	// 2. Fetch shift details for schedule and tolerance
	shift, err := s.repo.GetShiftByID(ctx, req.ShiftID)
	if err != nil || shift.ID == "" {
		return model.EmployeeAttendance{}, errors.New("shift tidak ditemukan")
	}

	// 3. Geofencing GPS Verification
	outletLat, outletLon, maxRadius, _ := s.repo.GetOutletGeofence(ctx, req.OutletID)
	var distMeters *float64
	var userLat, userLng *float64

	if req.Latitude != 0 && req.Longitude != 0 {
		dist := util.CalculateDistanceMeters(req.Latitude, req.Longitude, outletLat, outletLon)
		distRounded := util.CalculateDistanceMeters(req.Latitude, req.Longitude, outletLat, outletLon)
		distMeters = &distRounded
		userLat = &req.Latitude
		userLng = &req.Longitude

		if dist > maxRadius {
			return model.EmployeeAttendance{}, fmt.Errorf("lokasi Anda berjarak %.1f meter dari outlet (maksimal radius %.0f meter). Mohon lakukan absensi di area outlet FAN Laundry", dist, maxRadius)
		}
	}

	now := time.Now()

	// 4. Punctuality Evaluation (Late Calculation)
	// Parse shift start time on today's date
	punctuality := model.PunctualityOnTime
	lateMinutes := 0

	shiftStartParts := strings.Split(shift.StartTime, ":")
	if len(shiftStartParts) >= 2 {
		var hr, min int
		fmt.Sscanf(shiftStartParts[0], "%d", &hr)
		fmt.Sscanf(shiftStartParts[1], "%d", &min)

		expectedStart := time.Date(now.Year(), now.Month(), now.Day(), hr, min, 0, 0, now.Location())
		toleranceLimit := expectedStart.Add(time.Duration(shift.LateToleranceMinutes) * time.Minute)

		if now.After(toleranceLimit) {
			punctuality = model.PunctualityLate
			lateMinutes = int(now.Sub(expectedStart).Minutes())
		}
	}

	att := model.EmployeeAttendance{
		EmployeeName:          empName,
		EmployeePhone:         strings.TrimSpace(req.EmployeePhone),
		ShiftID:               shift.ID,
		ShiftName:             shift.ShiftName,
		OutletID:              req.OutletID,
		ClockInAt:             now,
		ClockInPhotoURL:       req.PhotoBase64,
		ClockInLat:            userLat,
		ClockInLng:            userLng,
		ClockInDistanceMeters: distMeters,
		StartingCash:          req.StartingCash,
		Punctuality:           punctuality,
		LateMinutes:           lateMinutes,
		Notes:                 strings.TrimSpace(req.Notes),
	}

	return s.repo.CheckIn(ctx, att)
}

// -----------------------------------------------------------------------------
// Check-Out (Clock-Out)
// -----------------------------------------------------------------------------

func (s *attendanceService) CheckOut(ctx context.Context, req model.CheckOutRequest) (model.EmployeeAttendance, error) {
	if strings.TrimSpace(req.AttendanceID) == "" {
		return model.EmployeeAttendance{}, errors.New("attendance_id wajib diisi")
	}
	if strings.TrimSpace(req.PhotoBase64) == "" {
		return model.EmployeeAttendance{}, errors.New("foto penutupan shift di tempat laundry wajib diambil")
	}

	// 1. Geofencing check
	outletLat, outletLon, maxRadius, _ := s.repo.GetOutletGeofence(ctx, "")
	var distMeters *float64
	var userLat, userLng *float64

	if req.Latitude != 0 && req.Longitude != 0 {
		dist := util.CalculateDistanceMeters(req.Latitude, req.Longitude, outletLat, outletLon)
		distMeters = &dist
		userLat = &req.Latitude
		userLng = &req.Longitude

		if dist > maxRadius {
			return model.EmployeeAttendance{}, fmt.Errorf("lokasi Anda berjarak %.1f meter dari outlet (maksimal radius %.0f meter). Mohon lakukan penutupan shift di area outlet", dist, maxRadius)
		}
	}

	now := time.Now()
	punctuality := model.PunctualityOnTime
	earlyMinutes := 0

	return s.repo.CheckOut(
		ctx, req.AttendanceID, now, req.PhotoBase64,
		userLat, userLng, distMeters, req.EndingCash,
		punctuality, earlyMinutes, strings.TrimSpace(req.Notes),
	)
}

func (s *attendanceService) GetActiveAttendance(ctx context.Context, employeeName string) (model.EmployeeAttendance, error) {
	return s.repo.GetActiveAttendanceByEmployee(ctx, employeeName)
}

func (s *attendanceService) ListActiveAttendances(ctx context.Context) ([]model.EmployeeAttendance, error) {
	return s.repo.ListActiveAttendances(ctx)
}

func (s *attendanceService) ListAttendances(
	ctx context.Context,
	startDate, endDate, employeeName string,
	page, limit int,
) ([]model.EmployeeAttendance, int, error) {
	if limit <= 0 {
		limit = 15
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	return s.repo.ListAttendances(ctx, startDate, endDate, employeeName, limit, offset)
}

// -----------------------------------------------------------------------------
// Leave Management (Jatah Libur Max 2x Sebulan)
// -----------------------------------------------------------------------------

func (s *attendanceService) GetLeaveQuota(ctx context.Context, employeeName, dateStr string) (model.LeaveQuotaResponse, error) {
	refTime := time.Now()
	if dateStr != "" {
		if t, err := time.Parse("2006-01-02", dateStr); err == nil {
			refTime = t
		}
	}

	year := refTime.Year()
	month := int(refTime.Month())
	monthStr := fmt.Sprintf("%04d-%02d", year, month)

	used, err := s.repo.GetMonthlyLeaveCount(ctx, employeeName, year, month)
	if err != nil {
		return model.LeaveQuotaResponse{}, err
	}

	total := 2
	remaining := total - used
	if remaining < 0 {
		remaining = 0
	}

	return model.LeaveQuotaResponse{
		Month:          monthStr,
		EmployeeName:   employeeName,
		TotalQuota:     total,
		UsedQuota:      used,
		RemainingQuota: remaining,
		CanRequest:     remaining > 0,
	}, nil
}

func (s *attendanceService) CreateLeaveRequest(ctx context.Context, req model.CreateLeaveRequest) (model.LeaveRequest, error) {
	if strings.TrimSpace(req.LeaveDate) == "" {
		req.LeaveDate = strings.TrimSpace(req.Date)
	}
	empName := strings.TrimSpace(req.EmployeeName)
	if empName == "" {
		return model.LeaveRequest{}, errors.New("nama karyawan wajib diisi")
	}
	if strings.TrimSpace(req.LeaveDate) == "" {
		return model.LeaveRequest{}, errors.New("tanggal libur wajib dipilih")
	}
	if strings.TrimSpace(req.Reason) == "" {
		return model.LeaveRequest{}, errors.New("alasan libur wajib diisi")
	}

	// 1. Validate max 2x quota per month
	quota, err := s.GetLeaveQuota(ctx, empName, req.LeaveDate)
	if err != nil {
		return model.LeaveRequest{}, err
	}
	if !quota.CanRequest {
		return model.LeaveRequest{}, fmt.Errorf("jatah libur untuk bulan %s sudah habis (maksimal %d hari per bulan)", quota.Month, quota.TotalQuota)
	}

	// 2. Check conflict with other approved staff
	conflict, otherStaff, _ := s.repo.HasConflictLeave(ctx, req.LeaveDate, empName)
	if conflict {
		req.Reason = fmt.Sprintf("%s [Catatan: Staf lain (%s) juga libur di tanggal ini]", req.Reason, otherStaff)
	}

	return s.repo.CreateLeaveRequest(ctx, req)
}

func (s *attendanceService) ListLeaveRequests(
	ctx context.Context,
	status, employeeName string,
	page, limit int,
) ([]model.LeaveRequest, int, error) {
	if limit <= 0 {
		limit = 15
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	return s.repo.ListLeaveRequests(ctx, status, employeeName, limit, offset)
}

func (s *attendanceService) ReviewLeaveRequest(ctx context.Context, id string, req model.ReviewLeaveRequest) (model.LeaveRequest, error) {
	if req.Status != model.LeaveStatusApproved && req.Status != model.LeaveStatusRejected {
		return model.LeaveRequest{}, errors.New("status review harus APPROVED atau REJECTED")
	}
	reviewer := strings.TrimSpace(req.ReviewedBy)
	if reviewer == "" {
		reviewer = "Owner"
	}
	notes := strings.TrimSpace(req.ReviewNotes)
	if notes == "" {
		notes = strings.TrimSpace(req.ReviewerNotes)
	}
	return s.repo.ReviewLeaveRequest(ctx, id, req.Status, reviewer, notes)
}
