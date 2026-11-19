package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/nuryanto2121/laundryos/apps/api/internal/model"
	"github.com/nuryanto2121/laundryos/apps/api/internal/util"
)

// OrderService handles business operations for orders.
type OrderService struct {
	repo OrderRepository
}

// NewOrderService creates a new OrderService.
func NewOrderService(repo OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

// Create creates a new order with a deterministic order number and NEW status.
func (s *OrderService) Create(ctx context.Context, req model.CreateOrderRequest) (model.Order, error) {
	return s.createWithStatus(ctx, req, model.OrderStatusNew)
}

// CreatePickupOrder creates a new order with WAITING_PICKUP status.
func (s *OrderService) CreatePickupOrder(ctx context.Context, req model.CreateOrderRequest) (model.Order, error) {
	return s.createWithStatus(ctx, req, model.OrderStatusWaitingPickup)
}

func (s *OrderService) createWithStatus(ctx context.Context, req model.CreateOrderRequest, status model.OrderStatus) (model.Order, error) {
	latestNumber, err := s.repo.GetLatestOrderNumberToday(ctx)
	if err != nil {
		return model.Order{}, fmt.Errorf("service: get latest order number: %w", err)
	}

	sequence, err := nextOrderSequence(latestNumber, time.Now())
	if err != nil {
		return model.Order{}, err
	}

	return s.repo.Create(ctx, req, util.GenerateOrderNumber(sequence), status)
}

// Get retrieves an order by ID.
func (s *OrderService) Get(ctx context.Context, id string) (model.Order, error) {
	return s.repo.Get(ctx, id)
}

// GetLatestByCustomerPhone retrieves the latest order for a customer by their phone number.
func (s *OrderService) GetLatestByCustomerPhone(ctx context.Context, phone string) (model.OrderTrackingDetail, error) {
	return s.repo.GetLatestByCustomerPhone(ctx, phone)
}

// GetByOrderNumber retrieves tracking details for an order by its unique order number.
func (s *OrderService) GetByOrderNumber(ctx context.Context, orderNumber string) (model.OrderTrackingDetail, error) {
	return s.repo.GetByOrderNumber(ctx, orderNumber)
}

// List returns all orders.
func (s *OrderService) List(ctx context.Context) ([]model.Order, error) {
	return s.repo.List(ctx)
}

// Update updates mutable order fields while preserving order number and status.
func (s *OrderService) Update(ctx context.Context, id string, req model.UpdateOrderRequest) (model.Order, error) {
	return s.repo.Update(ctx, id, req)
}

// UpdateStatus transitions an order to a new lifecycle status.
func (s *OrderService) UpdateStatus(ctx context.Context, id string, status model.OrderStatus) (model.Order, error) {
	return s.repo.UpdateStatus(ctx, id, status)
}

// Delete permanently removes an order. The current schema has no deleted_at column.
func (s *OrderService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func nextOrderSequence(latest string, now time.Time) (int64, error) {
	if latest == "" {
		return 1, nil
	}

	prefix := "ORD" + now.Format("060102")
	if len(latest) != len(prefix)+6 || !strings.HasPrefix(latest, prefix) {
		return 0, errors.New("service: malformed latest order number")
	}

	sequence, err := strconv.ParseInt(latest[len(prefix):], 10, 64)
	if err != nil || sequence < 1 {
		return 0, errors.New("service: malformed latest order number")
	}
	return sequence + 1, nil
}
