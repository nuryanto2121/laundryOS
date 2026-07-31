package service

import (
	"context"

	"github.com/nuryanto2121/laundryos/apps/api/internal/model"
	"github.com/nuryanto2121/laundryos/apps/api/internal/repository"
)

// CustomerService handles business operations for customers.
// It depends only on domain models — no database-specific types.
type CustomerService struct {
	repo *repository.CustomerRepository
}

// NewCustomerService creates a new CustomerService.
func NewCustomerService(repo *repository.CustomerRepository) *CustomerService {
	return &CustomerService{repo: repo}
}

// Create creates a new customer.
func (s *CustomerService) Create(ctx context.Context, req model.CreateCustomerRequest) (model.Customer, error) {
	return s.repo.Create(ctx, req)
}

// Get retrieves a customer by string ID.
func (s *CustomerService) Get(ctx context.Context, id string) (model.Customer, error) {
	return s.repo.Get(ctx, id)
}

// List returns all active customers.
func (s *CustomerService) List(ctx context.Context) ([]model.Customer, error) {
	return s.repo.List(ctx)
}

// Update updates a customer.
func (s *CustomerService) Update(ctx context.Context, id string, req model.UpdateCustomerRequest) (model.Customer, error) {
	return s.repo.Update(ctx, id, req)
}

// Delete soft-deletes a customer.
func (s *CustomerService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
