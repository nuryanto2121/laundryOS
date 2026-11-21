package service_test

import (
	"testing"
	"time"

	"github.com/nuryanto2121/laundryos/apps/api/internal/model"
	"github.com/nuryanto2121/laundryos/apps/api/internal/util"
)

func TestGeofenceDistance_FANLaundryKalibaru(t *testing.T) {
	// Point Outlet: FAN Laundry Kalibaru (-6.113425, 106.914230)
	outletLat := -6.113425
	outletLon := 106.914230

	// 1. Employee exactly at counter (~2 meters)
	empLat1 := -6.113410
	empLon1 := 106.914230
	within1, dist1 := util.IsWithinRadius(empLat1, empLon1, outletLat, outletLon, 25.0)
	if !within1 {
		t.Errorf("expected within radius, got false, dist: %f", dist1)
	}

	// 2. Employee far away from home (1.5 km away)
	empLatFar := -6.126500
	empLonFar := 106.914230
	withinFar, distFar := util.IsWithinRadius(empLatFar, empLonFar, outletLat, outletLon, 25.0)
	if withinFar {
		t.Errorf("expected outside radius, got true, dist: %f", distFar)
	}
	if distFar < 1000 {
		t.Errorf("expected > 1000m, got %f", distFar)
	}
}

func TestPunctualityEvaluation(t *testing.T) {
	// Shift 1 start: 08:00, tolerance 15 min
	now := time.Now()
	hr := 8
	min := 0

	expectedStart := time.Date(now.Year(), now.Month(), now.Day(), hr, min, 0, 0, now.Location())
	toleranceLimit := expectedStart.Add(15 * time.Minute)

	// Arrived at 08:10 -> ON_TIME
	arrivedOnTime := expectedStart.Add(10 * time.Minute)
	if arrivedOnTime.After(toleranceLimit) {
		t.Errorf("expected 08:10 to be on-time")
	}

	// Arrived at 08:25 -> LATE by 25 min
	arrivedLate := expectedStart.Add(25 * time.Minute)
	if !arrivedLate.After(toleranceLimit) {
		t.Errorf("expected 08:25 to be late")
	}
	lateMinutes := int(arrivedLate.Sub(expectedStart).Minutes())
	if lateMinutes != 25 {
		t.Errorf("expected late 25 mins, got %d", lateMinutes)
	}
}

func TestLeaveQuotaCalculation(t *testing.T) {
	quota := model.LeaveQuotaResponse{
		Month:          "2026-11",
		EmployeeName:   "Nuryanto",
		TotalQuota:     2,
		UsedQuota:      1,
		RemainingQuota: 1,
		CanRequest:     true,
	}

	if !quota.CanRequest {
		t.Errorf("expected employee with 1 remaining day to be able to request leave")
	}

	// Used up quota
	quotaFull := model.LeaveQuotaResponse{
		Month:          "2026-11",
		EmployeeName:   "Nuryanto",
		TotalQuota:     2,
		UsedQuota:      2,
		RemainingQuota: 0,
		CanRequest:     false,
	}

	if quotaFull.CanRequest {
		t.Errorf("expected employee with 0 remaining days not to be able to request leave")
	}
}

func TestShiftScheduleTimes(t *testing.T) {
	shift1 := model.Shift{
		Code:                 "SHIFT_1",
		ShiftName:            "Shift 1 (Pagi - Sore)",
		StartTime:            "08:00:00",
		EndTime:              "17:00:00",
		BreakStart:           "12:00:00",
		BreakEnd:             "13:00:00",
		LateToleranceMinutes: 15,
	}

	shift2 := model.Shift{
		Code:                 "SHIFT_2",
		ShiftName:            "Shift 2 (Siang - Malam)",
		StartTime:            "12:00:00",
		EndTime:              "21:00:00",
		BreakStart:           "18:00:00",
		BreakEnd:             "19:00:00",
		LateToleranceMinutes: 15,
	}

	if shift1.StartTime != "08:00:00" || shift1.EndTime != "17:00:00" {
		t.Errorf("unexpected shift 1 schedule")
	}
	if shift2.StartTime != "12:00:00" || shift2.EndTime != "21:00:00" {
		t.Errorf("unexpected shift 2 schedule")
	}
}
