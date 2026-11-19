package service

import (
	"context"

	"github.com/nuryanto2121/laundryos/apps/api/internal/model"
)

// OrderRepository defines persistence operations for orders.
type OrderRepository interface {
	Create(ctx context.Context, req model.CreateOrderRequest, orderNumber string, status model.OrderStatus) (model.Order, error)
	Get(ctx context.Context, id string) (model.Order, error)
	GetLatestOrderNumberToday(ctx context.Context) (string, error)
	GetLatestByCustomerPhone(ctx context.Context, phone string) (model.OrderTrackingDetail, error)
	GetByOrderNumber(ctx context.Context, orderNumber string) (model.OrderTrackingDetail, error)
	List(ctx context.Context) ([]model.Order, error)
	Update(ctx context.Context, id string, req model.UpdateOrderRequest) (model.Order, error)
	UpdateStatus(ctx context.Context, id string, status model.OrderStatus) (model.Order, error)
	Delete(ctx context.Context, id string) error
}
