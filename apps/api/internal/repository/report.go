package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/nuryanto2121/laundryos/apps/api/internal/db/generated"
)

// ReportRepository handles reporting, metrics, and CRM analytics queries.
type ReportRepository struct {
	queries *db.Queries
}

// NewReportRepository creates a new ReportRepository.
func NewReportRepository(pool *pgxpool.Pool) *ReportRepository {
	return &ReportRepository{queries: db.New(pool)}
}

// GetHomeMetrics retrieves real-time KPI metrics for today.
func (r *ReportRepository) GetHomeMetrics(ctx context.Context) (db.GetHomeMetricsTodayRow, error) {
	row, err := r.queries.GetHomeMetricsToday(ctx)
	if err != nil {
		return db.GetHomeMetricsTodayRow{}, fmt.Errorf("report repo: get home metrics: %w", err)
	}
	return row, nil
}

// GetDailyRevenueChart retrieves daily revenue breakdown within a date range.
func (r *ReportRepository) GetDailyRevenueChart(ctx context.Context, start, end time.Time) ([]db.GetDailyRevenueChartRow, error) {
	rows, err := r.queries.GetDailyRevenueChart(ctx, db.GetDailyRevenueChartParams{
		CreatedAt:   pgtype.Timestamp{Time: start, Valid: true},
		CreatedAt_2: pgtype.Timestamp{Time: end, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("report repo: daily revenue: %w", err)
	}
	return rows, nil
}

// ListUnpaidOrders retrieves outstanding receivables (piutang).
func (r *ReportRepository) ListUnpaidOrders(ctx context.Context) ([]db.ListUnpaidOrdersRow, error) {
	rows, err := r.queries.ListUnpaidOrders(ctx)
	if err != nil {
		return nil, fmt.Errorf("report repo: unpaid orders: %w", err)
	}
	return rows, nil
}

// GetServiceVolumeBreakdown retrieves volume breakdown by service category and unit.
func (r *ReportRepository) GetServiceVolumeBreakdown(ctx context.Context, start, end time.Time) ([]db.GetServiceVolumeBreakdownRow, error) {
	rows, err := r.queries.GetServiceVolumeBreakdown(ctx, db.GetServiceVolumeBreakdownParams{
		CreatedAt:   pgtype.Timestamp{Time: start, Valid: true},
		CreatedAt_2: pgtype.Timestamp{Time: end, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("report repo: service volume: %w", err)
	}
	return rows, nil
}

// GetTopCustomersLeaderboard retrieves top customers by total spend.
func (r *ReportRepository) GetTopCustomersLeaderboard(ctx context.Context, limit int32) ([]db.GetTopCustomersLeaderboardRow, error) {
	if limit <= 0 {
		limit = 10
	}
	rows, err := r.queries.GetTopCustomersLeaderboard(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("report repo: top customers: %w", err)
	}
	return rows, nil
}

// GetCustomerCRMStats retrieves CRM stats and lifetime spend for a customer.
func (r *ReportRepository) GetCustomerCRMStats(ctx context.Context, customerIDStr string) (db.GetCustomerCRMStatsRow, error) {
	parsed, err := uuid.Parse(customerIDStr)
	if err != nil {
		return db.GetCustomerCRMStatsRow{}, fmt.Errorf("invalid customer id: %w", err)
	}
	return r.queries.GetCustomerCRMStats(ctx, pgtype.UUID{Bytes: parsed, Valid: true})
}
