package service

import (
	"context"

	"github.com/nuryanto2121/laundryos/apps/api/internal/model"
)

// CustomerRepository defines the persistence operations for customers.
// It decouples the service layer from any concrete database or driver implementation.
type CustomerRepository interface {
	Create(ctx context.Context, req model.CreateCustomerRequest, customerCode string, status model.CustomerStatus) (model.Customer, error)
	Get(ctx context.Context, id string) (model.Customer, error)
	GetByPhone(ctx context.Context, phone string) (model.Customer, error)
	GetLatestCustomerCodeToday(ctx context.Context) (string, error)
	List(ctx context.Context) ([]model.Customer, error)
	ListPaginated(ctx context.Context, search, sortBy, sortDir string, limit, offset int) ([]model.Customer, int64, error)
	Update(ctx context.Context, id string, req model.UpdateCustomerRequest, status model.CustomerStatus, deletedAt *string) (model.Customer, error)
	UpdateStatus(ctx context.Context, id string, status model.CustomerStatus) (model.Customer, error)
	Delete(ctx context.Context, id string) error
	GetCustomerDefaultOutletID(ctx context.Context, customerID string) string
	SetCustomerOutletID(ctx context.Context, customerID string, outletID string) error
}
