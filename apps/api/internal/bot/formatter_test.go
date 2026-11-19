package bot

import (
	"strings"
	"testing"
	"time"

	"github.com/nuryanto2121/laundryos/apps/api/internal/model"
)

func TestFormatOrderStatusResponse_DisplayOrderNumber(t *testing.T) {
	order := &model.OrderTrackingDetail{
		OrderNumber: "ORD261118000001",
		ServiceName: "Cuci Setrika Regular",
		TotalWeight: 5.0,
		TotalPrice:  35000,
		Status:      model.OrderStatusInProgress,
		DurationDay: 3,
		CreatedAt:   "2026-11-19T11:10:00Z",
	}

	res := FormatOrderStatusResponse(order, "Budi")
	if !strings.Contains(res, "• *No. Order:* #261118000001") {
		t.Errorf("expected hash display order number, got: %s", res)
	}
	if !strings.Contains(res, "PROSES PENGERJAAN") {
		t.Errorf("expected humanized status, got: %s", res)
	}
	if !strings.Contains(res, "• *Estimasi Selesai:*") {
		t.Errorf("expected estimated completion time, got: %s", res)
	}
}

func TestCalculateCompletionEstimate(t *testing.T) {
	wib := time.FixedZone("WIB", 7*3600)
	baseTime := time.Date(2026, 11, 19, 11, 10, 0, 0, wib)

	// Test 1: Daily SLA (Cuci Setrika 3 days) -> 19-11-2026 11:10 + 3 days = 22-11-2026 11:10 WIB
	slaDaily := CalculateCompletionEstimate("Cuci Setrika Regular", 3, baseTime)
	if !strings.Contains(slaDaily, "22-11-2026 11:10 WIB") {
		t.Errorf("expected 22-11-2026 11:10 WIB, got: %s", slaDaily)
	}

	// Test 2: Self Service (Cuci Kering) -> 3 Jam SLA same day (11:10 + 3h = 14:10 WIB)
	slaSS := CalculateCompletionEstimate("Cuci Kering (Self Service)", 0, baseTime)
	if !strings.Contains(slaSS, "14:10 WIB") {
		t.Errorf("expected 14:10 WIB, got: %s", slaSS)
	}

	// Test 3: Self Service received at night (19:30 WIB) -> 19:30 + 3h = 22:30 (>21:00) -> Besok pagi
	nightTime := time.Date(2026, 11, 19, 19, 30, 0, 0, wib)
	slaNight := CalculateCompletionEstimate("Cuci Basah (Self Service)", 0, nightTime)
	if !strings.Contains(slaNight, "Besok pagi") {
		t.Errorf("expected Besok pagi for night order, got: %s", slaNight)
	}
}

func TestFormatDisplayOrderNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"ORD261119000012", "#261119000012"},
		{"#261119000012", "#261119000012"},
		{"261119000012", "#261119000012"},
		{"", ""},
	}

	for _, tt := range tests {
		got := FormatDisplayOrderNumber(tt.input)
		if got != tt.expected {
			t.Errorf("FormatDisplayOrderNumber(%q) = %q; want %q", tt.input, got, tt.expected)
		}
	}
}

func TestFormatIDR(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{0, "0"},
		{500, "500"},
		{1000, "1.000"},
		{25000, "25.000"},
		{1500000, "1.500.000"},
	}

	for _, tt := range tests {
		got := FormatIDR(tt.input)
		if got != tt.expected {
			t.Errorf("FormatIDR(%f) = %q; want %q", tt.input, got, tt.expected)
		}
	}
}
