package service

import (
	"context"
	"time"

	db "github.com/nuryanto2121/laundryos/apps/api/internal/db/generated"
	"github.com/nuryanto2121/laundryos/apps/api/internal/repository"
)

// FinancialReport aggregates revenue, expenses, and net profit.
type FinancialReport struct {
	TotalOmset      float64                      `json:"total_omset"`
	TotalPendapatan float64                      `json:"total_pendapatan"`
	TotalPengeluaran float64                     `json:"total_pengeluaran"`
	LabaBersih      float64                      `json:"laba_bersih"`
	DailyChart      []db.GetDailyRevenueChartRow `json:"daily_chart"`
}

// ReportService handles business calculations for dashboards and financial statements.
type ReportService struct {
	reportRepo  *repository.ReportRepository
	expenseRepo *repository.ExpenseRepository
}

// NewReportService creates a new ReportService.
func NewReportService(reportRepo *repository.ReportRepository, expenseRepo *repository.ExpenseRepository) *ReportService {
	return &ReportService{
		reportRepo:  reportRepo,
		expenseRepo: expenseRepo,
	}
}

// GetHomeMetrics returns the 4 primary KPI cards for the Home page.
func (s *ReportService) GetHomeMetrics(ctx context.Context) (db.GetHomeMetricsTodayRow, error) {
	return s.reportRepo.GetHomeMetrics(ctx)
}

// GetFinancialReport compiles Omset, Cash Inflow, Expenses, Net Profit, and Daily Chart.
func (s *ReportService) GetFinancialReport(ctx context.Context, start, end time.Time) (FinancialReport, error) {
	chart, err := s.reportRepo.GetDailyRevenueChart(ctx, start, end)
	if err != nil {
		return FinancialReport{}, err
	}

	var totalOmset, totalPendapatan float64
	for _, row := range chart {
		omset, _ := row.Omset.Float64Value()
		pendapatan, _ := row.PendapatanKas.Float64Value()
		totalOmset += omset.Float64
		totalPendapatan += pendapatan.Float64
	}

	totalPengeluaran, err := s.expenseRepo.GetTotal(ctx, start, end)
	if err != nil {
		totalPengeluaran = 0
	}

	labaBersih := totalPendapatan - totalPengeluaran

	return FinancialReport{
		TotalOmset:       totalOmset,
		TotalPendapatan:  totalPendapatan,
		TotalPengeluaran: totalPengeluaran,
		LabaBersih:       labaBersih,
		DailyChart:       chart,
	}, nil
}

// ListUnpaidOrders returns orders with outstanding balance.
func (s *ReportService) ListUnpaidOrders(ctx context.Context) ([]db.ListUnpaidOrdersRow, error) {
	return s.reportRepo.ListUnpaidOrders(ctx)
}

// GetServiceVolumeReport returns breakdown per service category (Kg, Pcs, Tabung, Meter).
func (s *ReportService) GetServiceVolumeReport(ctx context.Context, start, end time.Time) ([]db.GetServiceVolumeBreakdownRow, error) {
	return s.reportRepo.GetServiceVolumeBreakdown(ctx, start, end)
}

// GetTopCustomers returns top customers leaderboard.
func (s *ReportService) GetTopCustomers(ctx context.Context, limit int32) ([]db.GetTopCustomersLeaderboardRow, error) {
	return s.reportRepo.GetTopCustomersLeaderboard(ctx, limit)
}

// GetCustomerCRMProfile returns lifetime stats for an individual customer.
func (s *ReportService) GetCustomerCRMProfile(ctx context.Context, customerID string) (db.GetCustomerCRMStatsRow, error) {
	return s.reportRepo.GetCustomerCRMStats(ctx, customerID)
}
