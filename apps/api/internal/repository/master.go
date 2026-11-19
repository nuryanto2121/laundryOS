package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/nuryanto2121/laundryos/apps/api/internal/db/generated"
)

// MasterRepository provides access to master data (perfumes, discounts, payment methods, outlets).
type MasterRepository struct {
	queries *db.Queries
}

// NewMasterRepository creates a new MasterRepository.
func NewMasterRepository(pool *pgxpool.Pool) *MasterRepository {
	return &MasterRepository{queries: db.New(pool)}
}

// ListPerfumes retrieves all active perfumes.
func (r *MasterRepository) ListPerfumes(ctx context.Context) ([]db.Perfume, error) {
	perfumes, err := r.queries.ListActivePerfumes(ctx)
	if err != nil {
		return nil, fmt.Errorf("master repo: list perfumes: %w", err)
	}
	return perfumes, nil
}

// ListDiscounts retrieves all active discounts.
func (r *MasterRepository) ListDiscounts(ctx context.Context) ([]db.Discount, error) {
	discounts, err := r.queries.ListActiveDiscounts(ctx)
	if err != nil {
		return nil, fmt.Errorf("master repo: list discounts: %w", err)
	}
	return discounts, nil
}

// ListPaymentMethods retrieves all active payment methods.
func (r *MasterRepository) ListPaymentMethods(ctx context.Context) ([]db.PaymentMethod, error) {
	methods, err := r.queries.ListActivePaymentMethods(ctx)
	if err != nil {
		return nil, fmt.Errorf("master repo: list payment methods: %w", err)
	}
	return methods, nil
}

// ListOutlets retrieves all active outlets.
func (r *MasterRepository) ListOutlets(ctx context.Context) ([]db.Outlet, error) {
	outlets, err := r.queries.ListOutlets(ctx)
	if err != nil {
		return nil, fmt.Errorf("master repo: list outlets: %w", err)
	}
	return outlets, nil
}

// GetDefaultOutlet retrieves the primary default outlet.
func (r *MasterRepository) GetDefaultOutlet(ctx context.Context) (db.Outlet, error) {
	outlet, err := r.queries.GetDefaultOutlet(ctx)
	if err != nil {
		return db.Outlet{}, fmt.Errorf("master repo: get default outlet: %w", err)
	}
	return outlet, nil
}
