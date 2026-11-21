package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nuryanto2121/laundryos/apps/api/internal/model"
)

// AttendanceRepository handles database persistence for shifts, attendance logs, and leave requests.
type AttendanceRepository struct {
	pool *pgxpool.Pool
}

// NewAttendanceRepository creates a new AttendanceRepository.
func NewAttendanceRepository(pool *pgxpool.Pool) *AttendanceRepository {
	return &AttendanceRepository{pool: pool}
}

// -----------------------------------------------------------------------------
// Master Shifts
// -----------------------------------------------------------------------------

// ListShifts retrieves all active shifts sorted by start time.
func (r *AttendanceRepository) ListShifts(ctx context.Context) ([]model.Shift, error) {
	query := `
		SELECT id, code, shift_name, 
		       to_char(start_time, 'HH24:MI:SS'), 
		       to_char(end_time, 'HH24:MI:SS'), 
		       COALESCE(to_char(break_start, 'HH24:MI:SS'), ''), 
		       COALESCE(to_char(break_end, 'HH24:MI:SS'), ''), 
		       late_tolerance_minutes, early_leave_tolerance_minutes, 
		       is_active, created_at, updated_at
		FROM shifts
		ORDER BY start_time ASC
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("repository: list shifts: %w", err)
	}
	defer rows.Close()

	var list []model.Shift
	for rows.Next() {
		var s model.Shift
		var id uuid.UUID
		err := rows.Scan(
			&id, &s.Code, &s.ShiftName, &s.StartTime, &s.EndTime,
			&s.BreakStart, &s.BreakEnd, &s.LateToleranceMinutes,
			&s.EarlyLeaveToleranceMinutes, &s.IsActive, &s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("repository: scan shift: %w", err)
		}
		s.ID = id.String()
		list = append(list, s)
	}
	return list, nil
}

// GetShiftByID retrieves a single shift by ID.
func (r *AttendanceRepository) GetShiftByID(ctx context.Context, id string) (model.Shift, error) {
	query := `
		SELECT id, code, shift_name, 
		       to_char(start_time, 'HH24:MI:SS'), 
		       to_char(end_time, 'HH24:MI:SS'), 
		       COALESCE(to_char(break_start, 'HH24:MI:SS'), ''), 
		       COALESCE(to_char(break_end, 'HH24:MI:SS'), ''), 
		       late_tolerance_minutes, early_leave_tolerance_minutes, 
		       is_active, created_at, updated_at
		FROM shifts
		WHERE id = $1
	`
	var s model.Shift
	var uid uuid.UUID
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&uid, &s.Code, &s.ShiftName, &s.StartTime, &s.EndTime,
		&s.BreakStart, &s.BreakEnd, &s.LateToleranceMinutes,
		&s.EarlyLeaveToleranceMinutes, &s.IsActive, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Shift{}, nil
		}
		return model.Shift{}, fmt.Errorf("repository: get shift by id: %w", err)
	}
	s.ID = uid.String()
	return s, nil
}

// UpdateShift updates shift timings and tolerances.
func (r *AttendanceRepository) UpdateShift(ctx context.Context, id string, req model.UpdateShiftRequest) (model.Shift, error) {
	query := `
		UPDATE shifts
		SET shift_name = $2,
		    start_time = $3::time,
		    end_time = $4::time,
		    break_start = NULLIF($5, '')::time,
		    break_end = NULLIF($6, '')::time,
		    late_tolerance_minutes = $7,
		    early_leave_tolerance_minutes = $8,
		    is_active = $9,
		    updated_at = NOW()
		WHERE id = $1
		RETURNING id, code, shift_name, 
		          to_char(start_time, 'HH24:MI:SS'), 
		          to_char(end_time, 'HH24:MI:SS'), 
		          COALESCE(to_char(break_start, 'HH24:MI:SS'), ''), 
		          COALESCE(to_char(break_end, 'HH24:MI:SS'), ''), 
		          late_tolerance_minutes, early_leave_tolerance_minutes, 
		          is_active, created_at, updated_at
	`
	var s model.Shift
	var uid uuid.UUID
	err := r.pool.QueryRow(ctx, query,
		id, req.ShiftName, req.StartTime, req.EndTime,
		req.BreakStart, req.BreakEnd, req.LateToleranceMinutes,
		req.EarlyLeaveToleranceMinutes, req.IsActive,
	).Scan(
		&uid, &s.Code, &s.ShiftName, &s.StartTime, &s.EndTime,
		&s.BreakStart, &s.BreakEnd, &s.LateToleranceMinutes,
		&s.EarlyLeaveToleranceMinutes, &s.IsActive, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return model.Shift{}, fmt.Errorf("repository: update shift: %w", err)
	}
	s.ID = uid.String()
	return s, nil
}

// GetOutletGeofence returns latitude, longitude, and geofence radius for an outlet.
func (r *AttendanceRepository) GetOutletGeofence(ctx context.Context, outletID string) (float64, float64, float64, error) {
	query := `
		SELECT COALESCE(latitude, -6.113425), 
		       COALESCE(longitude, 106.914230), 
		       COALESCE(geofence_radius_meters, 25)
		FROM outlets
	`
	var lat, lon float64
	var radius int
	var err error

	if outletID != "" {
		err = r.pool.QueryRow(ctx, query+" WHERE id = $1 LIMIT 1", outletID).Scan(&lat, &lon, &radius)
	} else {
		err = r.pool.QueryRow(ctx, query+" LIMIT 1").Scan(&lat, &lon, &radius)
	}

	if err != nil {
		// Fallback to default coordinates of FAN Laundry Kalibaru
		return -6.113425, 106.914230, 25.0, nil
	}
	return lat, lon, float64(radius), nil
}

// -----------------------------------------------------------------------------
// Employee Attendances
// -----------------------------------------------------------------------------

// CheckIn inserts a new attendance record marking employee as ON_DUTY.
func (r *AttendanceRepository) CheckIn(ctx context.Context, att model.EmployeeAttendance) (model.EmployeeAttendance, error) {
	id := uuid.New()
	query := `
		INSERT INTO employee_attendances (
			id, employee_name, employee_phone, shift_id, shift_name, outlet_id,
			clock_in_at, clock_in_photo_url, clock_in_lat, clock_in_lng, clock_in_distance_meters,
			starting_cash, status, punctuality, late_minutes, notes, created_at, updated_at
		) VALUES (
			$1, $2, $3, NULLIF($4, '')::uuid, $5, NULLIF($6, '')::uuid,
			$7, $8, $9, $10, $11,
			$12, 'ON_DUTY', $13, $14, $15, NOW(), NOW()
		)
		RETURNING id, created_at, updated_at
	`
	var uid uuid.UUID
	var createdAt, updatedAt time.Time
	err := r.pool.QueryRow(ctx, query,
		id, att.EmployeeName, att.EmployeePhone, att.ShiftID, att.ShiftName, att.OutletID,
		att.ClockInAt, att.ClockInPhotoURL, att.ClockInLat, att.ClockInLng, att.ClockInDistanceMeters,
		att.StartingCash, string(att.Punctuality), att.LateMinutes, att.Notes,
	).Scan(&uid, &createdAt, &updatedAt)
	if err != nil {
		return model.EmployeeAttendance{}, fmt.Errorf("repository: check-in attendance: %w", err)
	}

	att.ID = uid.String()
	att.Status = model.AttendanceStatusOnDuty
	att.CreatedAt = createdAt
	att.UpdatedAt = updatedAt
	return att, nil
}

// CheckOut updates an existing ON_DUTY attendance record marking it COMPLETED.
func (r *AttendanceRepository) CheckOut(
	ctx context.Context,
	id string,
	clockOutAt time.Time,
	photoURL string,
	lat, lon, dist *float64,
	endingCash float64,
	punctuality model.PunctualityStatus,
	earlyLeaveMinutes int,
	notes string,
) (model.EmployeeAttendance, error) {
	query := `
		UPDATE employee_attendances
		SET clock_out_at = $2,
		    clock_out_photo_url = $3,
		    clock_out_lat = $4,
		    clock_out_lng = $5,
		    clock_out_distance_meters = $6,
		    ending_cash = $7,
		    status = 'COMPLETED',
		    punctuality = $8,
		    early_leave_minutes = $9,
		    notes = CASE WHEN $10 <> '' THEN COALESCE(notes, '') || ' | ' || $10 ELSE notes END,
		    updated_at = NOW()
		WHERE id = $1 AND status = 'ON_DUTY'
		RETURNING id, employee_name, COALESCE(employee_phone, ''), 
		          COALESCE(shift_id::text, ''), shift_name, COALESCE(outlet_id::text, ''),
		          clock_in_at, COALESCE(clock_in_photo_url, ''), clock_in_lat, clock_in_lng, clock_in_distance_meters,
		          clock_out_at, COALESCE(clock_out_photo_url, ''), clock_out_lat, clock_out_lng, clock_out_distance_meters,
		          starting_cash, ending_cash, status, punctuality, late_minutes, early_leave_minutes,
		          COALESCE(notes, ''), created_at, updated_at
	`
	var a model.EmployeeAttendance
	var uid uuid.UUID
	err := r.pool.QueryRow(ctx, query,
		id, clockOutAt, photoURL, lat, lon, dist, endingCash, string(punctuality), earlyLeaveMinutes, notes,
	).Scan(
		&uid, &a.EmployeeName, &a.EmployeePhone,
		&a.ShiftID, &a.ShiftName, &a.OutletID,
		&a.ClockInAt, &a.ClockInPhotoURL, &a.ClockInLat, &a.ClockInLng, &a.ClockInDistanceMeters,
		&a.ClockOutAt, &a.ClockOutPhotoURL, &a.ClockOutLat, &a.ClockOutLng, &a.ClockOutDistanceMeters,
		&a.StartingCash, &a.EndingCash, &a.Status, &a.Punctuality, &a.LateMinutes, &a.EarlyLeaveMinutes,
		&a.Notes, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.EmployeeAttendance{}, errors.New("attendance record not found or already completed")
		}
		return model.EmployeeAttendance{}, fmt.Errorf("repository: check-out attendance: %w", err)
	}
	a.ID = uid.String()
	return a, nil
}

// GetActiveAttendanceByEmployee returns active ON_DUTY attendance for an employee.
func (r *AttendanceRepository) GetActiveAttendanceByEmployee(ctx context.Context, employeeName string) (model.EmployeeAttendance, error) {
	query := `
		SELECT id, employee_name, COALESCE(employee_phone, ''), 
		       COALESCE(shift_id::text, ''), shift_name, COALESCE(outlet_id::text, ''),
		       clock_in_at, COALESCE(clock_in_photo_url, ''), clock_in_lat, clock_in_lng, clock_in_distance_meters,
		       clock_out_at, COALESCE(clock_out_photo_url, ''), clock_out_lat, clock_out_lng, clock_out_distance_meters,
		       starting_cash, ending_cash, status, punctuality, late_minutes, early_leave_minutes,
		       COALESCE(notes, ''), created_at, updated_at
		FROM employee_attendances
		WHERE employee_name = $1 AND status = 'ON_DUTY'
		ORDER BY clock_in_at DESC
		LIMIT 1
	`
	var a model.EmployeeAttendance
	var uid uuid.UUID
	err := r.pool.QueryRow(ctx, query, employeeName).Scan(
		&uid, &a.EmployeeName, &a.EmployeePhone,
		&a.ShiftID, &a.ShiftName, &a.OutletID,
		&a.ClockInAt, &a.ClockInPhotoURL, &a.ClockInLat, &a.ClockInLng, &a.ClockInDistanceMeters,
		&a.ClockOutAt, &a.ClockOutPhotoURL, &a.ClockOutLat, &a.ClockOutLng, &a.ClockOutDistanceMeters,
		&a.StartingCash, &a.EndingCash, &a.Status, &a.Punctuality, &a.LateMinutes, &a.EarlyLeaveMinutes,
		&a.Notes, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.EmployeeAttendance{}, nil
		}
		return model.EmployeeAttendance{}, fmt.Errorf("repository: get active attendance: %w", err)
	}
	a.ID = uid.String()
	return a, nil
}

// ListActiveAttendances returns all staff currently ON_DUTY.
func (r *AttendanceRepository) ListActiveAttendances(ctx context.Context) ([]model.EmployeeAttendance, error) {
	query := `
		SELECT id, employee_name, COALESCE(employee_phone, ''), 
		       COALESCE(shift_id::text, ''), shift_name, COALESCE(outlet_id::text, ''),
		       clock_in_at, COALESCE(clock_in_photo_url, ''), clock_in_lat, clock_in_lng, clock_in_distance_meters,
		       clock_out_at, COALESCE(clock_out_photo_url, ''), clock_out_lat, clock_out_lng, clock_out_distance_meters,
		       starting_cash, ending_cash, status, punctuality, late_minutes, early_leave_minutes,
		       COALESCE(notes, ''), created_at, updated_at
		FROM employee_attendances
		WHERE status = 'ON_DUTY'
		ORDER BY clock_in_at DESC
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("repository: list active attendances: %w", err)
	}
	defer rows.Close()

	var list []model.EmployeeAttendance
	for rows.Next() {
		var a model.EmployeeAttendance
		var uid uuid.UUID
		err := rows.Scan(
			&uid, &a.EmployeeName, &a.EmployeePhone,
			&a.ShiftID, &a.ShiftName, &a.OutletID,
			&a.ClockInAt, &a.ClockInPhotoURL, &a.ClockInLat, &a.ClockInLng, &a.ClockInDistanceMeters,
			&a.ClockOutAt, &a.ClockOutPhotoURL, &a.ClockOutLat, &a.ClockOutLng, &a.ClockOutDistanceMeters,
			&a.StartingCash, &a.EndingCash, &a.Status, &a.Punctuality, &a.LateMinutes, &a.EarlyLeaveMinutes,
			&a.Notes, &a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("repository: scan active attendance: %w", err)
		}
		a.ID = uid.String()
		list = append(list, a)
	}
	return list, nil
}

// ListAttendances retrieves history log with optional filtering.
func (r *AttendanceRepository) ListAttendances(
	ctx context.Context,
	startDate, endDate, employeeName string,
	limit, offset int,
) ([]model.EmployeeAttendance, int, error) {
	countQuery := `
		SELECT COUNT(*)
		FROM employee_attendances
		WHERE ($1 = '' OR clock_in_at::date >= $1::date)
		  AND ($2 = '' OR clock_in_at::date <= $2::date)
		  AND ($3 = '' OR employee_name ILIKE '%' || $3 || '%')
	`
	var total int
	err := r.pool.QueryRow(ctx, countQuery, startDate, endDate, employeeName).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("repository: count attendances: %w", err)
	}

	query := `
		SELECT id, employee_name, COALESCE(employee_phone, ''), 
		       COALESCE(shift_id::text, ''), shift_name, COALESCE(outlet_id::text, ''),
		       clock_in_at, COALESCE(clock_in_photo_url, ''), clock_in_lat, clock_in_lng, clock_in_distance_meters,
		       clock_out_at, COALESCE(clock_out_photo_url, ''), clock_out_lat, clock_out_lng, clock_out_distance_meters,
		       starting_cash, ending_cash, status, punctuality, late_minutes, early_leave_minutes,
		       COALESCE(notes, ''), created_at, updated_at
		FROM employee_attendances
		WHERE ($1 = '' OR clock_in_at::date >= $1::date)
		  AND ($2 = '' OR clock_in_at::date <= $2::date)
		  AND ($3 = '' OR employee_name ILIKE '%' || $3 || '%')
		ORDER BY clock_in_at DESC
		LIMIT $4 OFFSET $5
	`
	rows, err := r.pool.Query(ctx, query, startDate, endDate, employeeName, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("repository: list attendances: %w", err)
	}
	defer rows.Close()

	var list []model.EmployeeAttendance
	for rows.Next() {
		var a model.EmployeeAttendance
		var uid uuid.UUID
		err := rows.Scan(
			&uid, &a.EmployeeName, &a.EmployeePhone,
			&a.ShiftID, &a.ShiftName, &a.OutletID,
			&a.ClockInAt, &a.ClockInPhotoURL, &a.ClockInLat, &a.ClockInLng, &a.ClockInDistanceMeters,
			&a.ClockOutAt, &a.ClockOutPhotoURL, &a.ClockOutLat, &a.ClockOutLng, &a.ClockOutDistanceMeters,
			&a.StartingCash, &a.EndingCash, &a.Status, &a.Punctuality, &a.LateMinutes, &a.EarlyLeaveMinutes,
			&a.Notes, &a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("repository: scan attendance: %w", err)
		}
		a.ID = uid.String()
		list = append(list, a)
	}
	return list, total, nil
}

// -----------------------------------------------------------------------------
// Leave Requests (Pengajuan Libur: Max 2x Sebulan)
// -----------------------------------------------------------------------------

// GetMonthlyLeaveCount counts approved or pending leaves in a given month.
func (r *AttendanceRepository) GetMonthlyLeaveCount(ctx context.Context, employeeName string, year, month int) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM leave_requests
		WHERE employee_name = $1
		  AND EXTRACT(YEAR FROM leave_date) = $2
		  AND EXTRACT(MONTH FROM leave_date) = $3
		  AND status IN ('PENDING', 'APPROVED')
	`
	var count int
	err := r.pool.QueryRow(ctx, query, employeeName, year, month).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("repository: get monthly leave count: %w", err)
	}
	return count, nil
}

// HasConflictLeave checks if another employee is already APPROVED to leave on the same date.
func (r *AttendanceRepository) HasConflictLeave(ctx context.Context, leaveDate string, excludeEmployee string) (bool, string, error) {
	query := `
		SELECT employee_name
		FROM leave_requests
		WHERE leave_date = $1::date
		  AND status = 'APPROVED'
		  AND ($2 = '' OR employee_name <> $2)
		LIMIT 1
	`
	var emp string
	err := r.pool.QueryRow(ctx, query, leaveDate, excludeEmployee).Scan(&emp)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, "", nil
		}
		return false, "", fmt.Errorf("repository: check leave conflict: %w", err)
	}
	return true, emp, nil
}

// CreateLeaveRequest inserts a new leave request.
func (r *AttendanceRepository) CreateLeaveRequest(ctx context.Context, req model.CreateLeaveRequest) (model.LeaveRequest, error) {
	id := uuid.New()
	query := `
		INSERT INTO leave_requests (
			id, employee_name, employee_phone, leave_date, reason, status, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4::date, $5, 'PENDING', NOW(), NOW()
		)
		RETURNING id, created_at, updated_at
	`
	var uid uuid.UUID
	var createdAt, updatedAt time.Time
	err := r.pool.QueryRow(ctx, query,
		id, req.EmployeeName, req.EmployeePhone, req.LeaveDate, req.Reason,
	).Scan(&uid, &createdAt, &updatedAt)
	if err != nil {
		return model.LeaveRequest{}, fmt.Errorf("repository: create leave request: %w", err)
	}

	return model.LeaveRequest{
		ID:            uid.String(),
		EmployeeName:  req.EmployeeName,
		EmployeePhone: req.EmployeePhone,
		LeaveDate:     req.LeaveDate,
		Reason:        req.Reason,
		Status:        model.LeaveStatusPending,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}, nil
}

// ListLeaveRequests retrieves all leave applications.
func (r *AttendanceRepository) ListLeaveRequests(
	ctx context.Context,
	status, employeeName string,
	limit, offset int,
) ([]model.LeaveRequest, int, error) {
	countQuery := `
		SELECT COUNT(*)
		FROM leave_requests
		WHERE ($1 = '' OR status = $1)
		  AND ($2 = '' OR employee_name ILIKE '%' || $2 || '%')
	`
	var total int
	err := r.pool.QueryRow(ctx, countQuery, status, employeeName).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("repository: count leave requests: %w", err)
	}

	query := `
		SELECT id, employee_name, COALESCE(employee_phone, ''), 
		       to_char(leave_date, 'YYYY-MM-DD'), reason, status, 
		       COALESCE(reviewed_by, ''), reviewed_at, COALESCE(review_notes, ''),
		       created_at, updated_at
		FROM leave_requests
		WHERE ($1 = '' OR status = $1)
		  AND ($2 = '' OR employee_name ILIKE '%' || $2 || '%')
		ORDER BY leave_date DESC, created_at DESC
		LIMIT $3 OFFSET $4
	`
	rows, err := r.pool.Query(ctx, query, status, employeeName, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("repository: list leave requests: %w", err)
	}
	defer rows.Close()

	var list []model.LeaveRequest
	for rows.Next() {
		var lr model.LeaveRequest
		var uid uuid.UUID
		err := rows.Scan(
			&uid, &lr.EmployeeName, &lr.EmployeePhone,
			&lr.LeaveDate, &lr.Reason, &lr.Status,
			&lr.ReviewedBy, &lr.ReviewedAt, &lr.ReviewNotes,
			&lr.CreatedAt, &lr.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("repository: scan leave request: %w", err)
		}
		lr.ID = uid.String()
		list = append(list, lr)
	}
	return list, total, nil
}

// ReviewLeaveRequest updates status (APPROVED or REJECTED) with reviewer notes.
func (r *AttendanceRepository) ReviewLeaveRequest(
	ctx context.Context,
	id string,
	status model.LeaveStatus,
	reviewedBy, reviewNotes string,
) (model.LeaveRequest, error) {
	query := `
		UPDATE leave_requests
		SET status = $2,
		    reviewed_by = $3,
		    review_notes = $4,
		    reviewed_at = NOW(),
		    updated_at = NOW()
		WHERE id = $1
		RETURNING id, employee_name, COALESCE(employee_phone, ''), 
		          to_char(leave_date, 'YYYY-MM-DD'), reason, status, 
		          COALESCE(reviewed_by, ''), reviewed_at, COALESCE(review_notes, ''),
		          created_at, updated_at
	`
	var lr model.LeaveRequest
	var uid uuid.UUID
	err := r.pool.QueryRow(ctx, query, id, string(status), reviewedBy, reviewNotes).Scan(
		&uid, &lr.EmployeeName, &lr.EmployeePhone,
		&lr.LeaveDate, &lr.Reason, &lr.Status,
		&lr.ReviewedBy, &lr.ReviewedAt, &lr.ReviewNotes,
		&lr.CreatedAt, &lr.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.LeaveRequest{}, errors.New("leave request not found")
		}
		return model.LeaveRequest{}, fmt.Errorf("repository: review leave request: %w", err)
	}
	lr.ID = uid.String()
	return lr, nil
}
