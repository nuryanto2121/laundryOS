package model

import "time"

// PunctualityStatus represents employee attendance status
type PunctualityStatus string

const (
	PunctualityOnTime     PunctualityStatus = "ON_TIME"
	PunctualityLate       PunctualityStatus = "LATE"
	PunctualityEarlyLeave PunctualityStatus = "EARLY_LEAVE"
	PunctualityOvertime   PunctualityStatus = "OVERTIME"
)

// AttendanceStatus represents active lifecycle of a shift
type AttendanceStatus string

const (
	AttendanceStatusOnDuty    AttendanceStatus = "ON_DUTY"
	AttendanceStatusCompleted AttendanceStatus = "COMPLETED"
)

// LeaveStatus represents employee leave request review state
type LeaveStatus string

const (
	LeaveStatusPending   LeaveStatus = "PENDING"
	LeaveStatusApproved  LeaveStatus = "APPROVED"
	LeaveStatusRejected  LeaveStatus = "REJECTED"
	LeaveStatusCancelled LeaveStatus = "CANCELLED"
)

// Shift represents master shift schedule
type Shift struct {
	ID                         string    `json:"id"`
	Code                       string    `json:"code"`
	ShiftName                  string    `json:"shift_name"`
	StartTime                  string    `json:"start_time"` // "08:00:00"
	EndTime                    string    `json:"end_time"`   // "17:00:00"
	BreakStart                 string    `json:"break_start"`
	BreakEnd                   string    `json:"break_end"`
	LateToleranceMinutes       int       `json:"late_tolerance_minutes"`
	EarlyLeaveToleranceMinutes int       `json:"early_leave_tolerance_minutes"`
	IsActive                   bool      `json:"is_active"`
	CreatedAt                  time.Time `json:"created_at"`
	UpdatedAt                  time.Time `json:"updated_at"`
}

// UpdateShiftRequest represents payload to configure shift timings and tolerances
type UpdateShiftRequest struct {
	ShiftName                  string `json:"shift_name" validate:"required"`
	StartTime                  string `json:"start_time" validate:"required"`
	EndTime                    string `json:"end_time" validate:"required"`
	BreakStart                 string `json:"break_start"`
	BreakEnd                   string `json:"break_end"`
	LateToleranceMinutes       int    `json:"late_tolerance_minutes"`
	EarlyLeaveToleranceMinutes int    `json:"early_leave_tolerance_minutes"`
	IsActive                   bool   `json:"is_active"`
}

// EmployeeAttendance represents attendance and cashier shift log
type EmployeeAttendance struct {
	ID                      string            `json:"id"`
	EmployeeName            string            `json:"employee_name"`
	EmployeePhone           string            `json:"employee_phone"`
	ShiftID                 string            `json:"shift_id"`
	ShiftName               string            `json:"shift_name"`
	OutletID                string            `json:"outlet_id"`
	ClockInAt               time.Time         `json:"clock_in_at"`
	ClockInPhotoURL         string            `json:"clock_in_photo_url"`
	ClockInLat              *float64          `json:"clock_in_lat"`
	ClockInLng              *float64          `json:"clock_in_lng"`
	ClockInDistanceMeters   *float64          `json:"clock_in_distance_meters"`
	ClockOutAt              *time.Time        `json:"clock_out_at"`
	ClockOutPhotoURL        string            `json:"clock_out_photo_url"`
	ClockOutLat             *float64          `json:"clock_out_lat"`
	ClockOutLng             *float64          `json:"clock_out_lng"`
	ClockOutDistanceMeters  *float64          `json:"clock_out_distance_meters"`
	StartingCash            float64           `json:"starting_cash"`
	EndingCash              float64           `json:"ending_cash"`
	Status                  AttendanceStatus  `json:"status"`
	Punctuality             PunctualityStatus `json:"punctuality"`
	LateMinutes             int               `json:"late_minutes"`
	EarlyLeaveMinutes       int               `json:"early_leave_minutes"`
	Notes                   string            `json:"notes"`
	CreatedAt               time.Time         `json:"created_at"`
	UpdatedAt               time.Time         `json:"updated_at"`
}

// CheckInRequest represents payload submitted when employee starts shift
type CheckInRequest struct {
	EmployeeName string  `json:"employee_name" validate:"required"`
	EmployeePhone string `json:"employee_phone"`
	ShiftID      string  `json:"shift_id" validate:"required"`
	OutletID     string  `json:"outlet_id"`
	StartingCash float64 `json:"starting_cash"`
	PhotoBase64  string  `json:"photo_base64" validate:"required"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Notes        string  `json:"notes"`
}

// CheckOutRequest represents payload submitted when employee ends shift
type CheckOutRequest struct {
	AttendanceID string  `json:"attendance_id" validate:"required"`
	EndingCash   float64 `json:"ending_cash"`
	PhotoBase64  string  `json:"photo_base64" validate:"required"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Notes        string  `json:"notes"`
}

// LeaveRequest represents employee leave application
type LeaveRequest struct {
	ID            string      `json:"id"`
	EmployeeName  string      `json:"employee_name"`
	EmployeePhone string      `json:"employee_phone"`
	LeaveDate     string      `json:"leave_date"` // YYYY-MM-DD
	Reason        string      `json:"reason"`
	Status        LeaveStatus `json:"status"`
	ReviewedBy    string      `json:"reviewed_by"`
	ReviewedAt    *time.Time  `json:"reviewed_at"`
	ReviewNotes   string      `json:"review_notes"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

// CreateLeaveRequest represents payload when employee applies for off-day
type CreateLeaveRequest struct {
	EmployeeName  string `json:"employee_name"`
	EmployeePhone string `json:"employee_phone"`
	LeaveDate     string `json:"leave_date"` // YYYY-MM-DD
	Date          string `json:"date"`       // Alias support
	Reason        string `json:"reason"`
}

// ReviewLeaveRequest represents Owner approval / rejection payload
type ReviewLeaveRequest struct {
	Status        LeaveStatus `json:"status"` // APPROVED or REJECTED
	ReviewNotes   string      `json:"review_notes"`
	ReviewerNotes string      `json:"reviewer_notes"` // Alias support
	ReviewedBy    string      `json:"reviewed_by"`
}

// LeaveQuotaResponse represents monthly leave quota calculation (2x per month)
type LeaveQuotaResponse struct {
	Month          string `json:"month"` // YYYY-MM
	EmployeeName   string `json:"employee_name"`
	TotalQuota     int    `json:"total_quota"`     // 2
	UsedQuota      int    `json:"used_quota"`      // count of approved or pending
	RemainingQuota int    `json:"remaining_quota"` // TotalQuota - UsedQuota
	CanRequest     bool   `json:"can_request"`
}
