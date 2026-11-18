package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/nuryanto2121/laundryos/apps/api/internal/model"
	"github.com/nuryanto2121/laundryos/apps/api/internal/util"
)

// CustomerService handles business operations for customers.
// It depends only on domain models — no database-specific types.
type CustomerService struct {
	repo CustomerRepository
}

// NewCustomerService creates a new CustomerService.
func NewCustomerService(repo CustomerRepository) *CustomerService {
	return &CustomerService{repo: repo}
}

// Create registers a new customer or reactivates an inactive existing customer,
// strictly enforcing business rules for Cases A, B, C, and D.
func (s *CustomerService) Create(ctx context.Context, req model.CreateCustomerRequest) (model.Customer, error) {
	existing, err := s.repo.GetByPhone(ctx, req.Phone)
	if err != nil {
		return model.Customer{}, fmt.Errorf("service: check duplicate phone: %w", err)
	}

	if existing.ID != "" {
		switch existing.Status {
		case model.CustomerStatusActive:
			// Case B: Customer Exists and is ACTIVE
			return model.Customer{}, errors.New("customer already registered")
		case model.CustomerStatusBlocked:
			// Case D: Customer Exists and is BLOCKED
			return model.Customer{}, errors.New("customer is blocked")
		case model.CustomerStatusInactive:
			// Case C: Customer Exists and is INACTIVE (Reactivation Flow)
			updateReq := model.UpdateCustomerRequest(req)
			return s.repo.Update(ctx, existing.ID, updateReq, model.CustomerStatusActive, nil)
		}
	}

	// Case A: Customer Not Found (New Registration)
	latestCode, err := s.repo.GetLatestCustomerCodeToday(ctx)
	if err != nil {
		// If an error is returned that wraps "no rows", handle it as normal condition (sequence = 1).
		// Otherwise, return persistence failure.
		if !strings.Contains(strings.ToLower(err.Error()), "no rows") && !strings.Contains(strings.ToLower(err.Error()), "not found") {
			return model.Customer{}, fmt.Errorf("service: get latest customer code: %w", err)
		}
		latestCode = ""
	}

	sequence := int64(1)
	if len(latestCode) >= 6 {
		tail := latestCode[len(latestCode)-6:]
		if parsed, err := strconv.ParseInt(tail, 10, 64); err == nil && parsed >= 1 {
			sequence = parsed + 1
		}
	}

	newCode := util.GenerateCustomerCode(sequence)
	return s.repo.Create(ctx, req, newCode, model.CustomerStatusActive)
}

// Get retrieves a customer by string ID.
func (s *CustomerService) Get(ctx context.Context, id string) (model.Customer, error) {
	return s.repo.Get(ctx, id)
}

// List returns all active customers.
func (s *CustomerService) List(ctx context.Context) ([]model.Customer, error) {
	return s.repo.List(ctx)
}

// Update updates a customer's editable profile information while strictly
// guaranteeing immutability of customer_code, status, created_at, and deleted_at.
func (s *CustomerService) Update(ctx context.Context, id string, req model.UpdateCustomerRequest) (model.Customer, error) {
	existing, err := s.repo.Get(ctx, id)
	if err != nil {
		return model.Customer{}, err
	}
	return s.repo.Update(ctx, id, req, existing.Status, existing.DeletedAt)
}

// Delete soft-deletes an active customer resource.
func (s *CustomerService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// UpdateStatus changes a customer's lifecycle status, strictly enforcing the Customer Status Transition Matrix.
func (s *CustomerService) UpdateStatus(ctx context.Context, id string, status model.CustomerStatus) (model.Customer, error) {
	if status != model.CustomerStatusActive && status != model.CustomerStatusInactive && status != model.CustomerStatusBlocked {
		return model.Customer{}, errors.New("invalid target status")
	}

	existing, err := s.repo.Get(ctx, id)
	if err != nil {
		return model.Customer{}, err
	}

	if existing.Status == status {
		return model.Customer{}, fmt.Errorf("customer already %s", status)
	}

	return s.repo.UpdateStatus(ctx, id, status)
}

// GetByPhone retrieves a customer by phone number.
func (s *CustomerService) GetByPhone(ctx context.Context, phone string) (model.Customer, error) {
	return s.repo.GetByPhone(ctx, phone)
}

// GetOrCreateByPhone retrieves an existing customer or automatically creates a new customer.
func (s *CustomerService) GetOrCreateByPhone(ctx context.Context, phone, fullName, address string) (model.Customer, error) {
	existing, err := s.repo.GetByPhone(ctx, phone)
	if err != nil {
		return model.Customer{}, fmt.Errorf("service: get customer by phone: %w", err)
	}

	if existing.ID != "" {
		if address != "" && existing.Address != address {
			updateReq := model.UpdateCustomerRequest{
				FullName: existing.FullName,
				Phone:    existing.Phone,
				Email:    existing.Email,
				Address:  address,
			}
			return s.repo.Update(ctx, existing.ID, updateReq, existing.Status, existing.DeletedAt)
		}
		return existing, nil
	}

	name := strings.TrimSpace(fullName)
	if name == "" {
		name = "Pelanggan " + phone
	}
	return s.Create(ctx, model.CreateCustomerRequest{
		FullName: name,
		Phone:    phone,
		Address:  address,
	})
}
