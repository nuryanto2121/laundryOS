package service

import (
	"context"

	db "github.com/nuryanto2121/laundryos/apps/api/internal/db/generated"
	"github.com/nuryanto2121/laundryos/apps/api/internal/repository"
)

// MasterService handles business logic for master lookup data.
type MasterService struct {
	repo *repository.MasterRepository
}

// NewMasterService creates a new MasterService.
func NewMasterService(repo *repository.MasterRepository) *MasterService {
	return &MasterService{repo: repo}
}

// ListPerfumes returns all active perfumes.
func (s *MasterService) ListPerfumes(ctx context.Context) ([]db.Perfume, error) {
	return s.repo.ListPerfumes(ctx)
}

// ListDiscounts returns all active discounts.
func (s *MasterService) ListDiscounts(ctx context.Context) ([]db.Discount, error) {
	return s.repo.ListDiscounts(ctx)
}

// ListPaymentMethods returns all active payment methods.
func (s *MasterService) ListPaymentMethods(ctx context.Context) ([]db.PaymentMethod, error) {
	return s.repo.ListPaymentMethods(ctx)
}

// ListOutlets returns all active outlets.
func (s *MasterService) ListOutlets(ctx context.Context) ([]db.Outlet, error) {
	return s.repo.ListOutlets(ctx)
}

// GetDefaultOutlet returns the primary outlet.
func (s *MasterService) GetDefaultOutlet(ctx context.Context) (db.Outlet, error) {
	return s.repo.GetDefaultOutlet(ctx)
}
