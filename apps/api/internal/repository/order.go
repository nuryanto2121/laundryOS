package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/nuryanto2121/laundryos/apps/api/internal/db/generated"
	"github.com/nuryanto2121/laundryos/apps/api/internal/model"
)

// OrderRepository handles persistence operations for orders.
type OrderRepository struct {
	queries *db.Queries
}

// NewOrderRepository creates a new OrderRepository.
func NewOrderRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{queries: db.New(pool)}
}

// Create inserts a new order and returns the created record.
func (r *OrderRepository) Create(ctx context.Context, req model.CreateOrderRequest, orderNumber string, status model.OrderStatus) (model.Order, error) {
	id := uuid.New()
	customerID, err := parseOrderUUID(req.CustomerID, "customer")
	if err != nil {
		return model.Order{}, err
	}
	serviceID, err := parseOrderUUID(req.ServiceID, "service")
	if err != nil {
		return model.Order{}, err
	}

	order, err := r.queries.CreateOrder(ctx, db.CreateOrderParams{
		ID:             pgtype.UUID{Bytes: id, Valid: true},
		OrderNumber:    orderNumber,
		CustomerID:     customerID,
		ServiceID:      serviceID,
		Status:         pgtype.Text{String: string(status), Valid: true},
		TotalWeight:    toOrderNumeric(req.TotalWeight),
		TotalPrice:     toOrderNumeric(req.TotalPrice),
		PickupRequired: pgtype.Bool{Bool: req.PickupRequired, Valid: true},
		Notes:          pgtype.Text{String: req.Notes, Valid: req.Notes != ""},
	})
	if err != nil {
		return model.Order{}, fmt.Errorf("repository: create order: %w", err)
	}
	return orderToModel(order), nil
}

// Get retrieves an order by string ID.
func (r *OrderRepository) Get(ctx context.Context, id string) (model.Order, error) {
	orderID, err := parseOrderUUID(id, "order")
	if err != nil {
		return model.Order{}, err
	}
	order, err := r.queries.GetOrder(ctx, orderID)
	if err != nil {
		return model.Order{}, fmt.Errorf("repository: get order: %w", err)
	}
	return orderToModel(order), nil
}

// GetLatestOrderNumberToday retrieves today's latest order number.
func (r *OrderRepository) GetLatestOrderNumberToday(ctx context.Context) (string, error) {
	orderNumber, err := r.queries.GetLatestOrderNumberToday(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", fmt.Errorf("repository: get latest order number today: %w", err)
	}
	return orderNumber, nil
}

// List returns all orders ordered by creation time.
func (r *OrderRepository) List(ctx context.Context) ([]model.Order, error) {
	orders, err := r.queries.ListOrders(ctx)
	if err != nil {
		return nil, fmt.Errorf("repository: list orders: %w", err)
	}
	result := make([]model.Order, len(orders))
	for i, order := range orders {
		result[i] = orderToModel(order)
	}
	return result, nil
}

// Update updates mutable order fields.
func (r *OrderRepository) Update(ctx context.Context, id string, req model.UpdateOrderRequest) (model.Order, error) {
	orderID, err := parseOrderUUID(id, "order")
	if err != nil {
		return model.Order{}, err
	}
	customerID, err := parseOrderUUID(req.CustomerID, "customer")
	if err != nil {
		return model.Order{}, err
	}
	serviceID, err := parseOrderUUID(req.ServiceID, "service")
	if err != nil {
		return model.Order{}, err
	}

	order, err := r.queries.UpdateOrder(ctx, db.UpdateOrderParams{
		ID:             orderID,
		CustomerID:     customerID,
		ServiceID:      serviceID,
		TotalWeight:    toOrderNumeric(req.TotalWeight),
		TotalPrice:     toOrderNumeric(req.TotalPrice),
		PickupRequired: pgtype.Bool{Bool: req.PickupRequired, Valid: true},
		Notes:          pgtype.Text{String: req.Notes, Valid: req.Notes != ""},
	})
	if err != nil {
		return model.Order{}, fmt.Errorf("repository: update order: %w", err)
	}
	return orderToModel(order), nil
}

// Delete permanently removes an order because orders has no deleted_at column.
func (r *OrderRepository) Delete(ctx context.Context, id string) error {
	orderID, err := parseOrderUUID(id, "order")
	if err != nil {
		return err
	}
	if err := r.queries.DeleteOrder(ctx, orderID); err != nil {
		return fmt.Errorf("repository: delete order: %w", err)
	}
	return nil
}

func parseOrderUUID(value string, name string) (pgtype.UUID, error) {
	parsed, err := uuid.Parse(value)
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("invalid %s id: %w", name, err)
	}
	return pgtype.UUID{Bytes: parsed, Valid: true}, nil
}

func toOrderNumeric(value float64) pgtype.Numeric {
	var numeric pgtype.Numeric
	_ = numeric.Scan(strconv.FormatFloat(value, 'f', -1, 64))
	return numeric
}

func orderNumericToFloat64(value pgtype.Numeric) float64 {
	if !value.Valid {
		return 0
	}
	result, err := value.Float64Value()
	if err != nil || !result.Valid {
		return 0
	}
	return result.Float64
}

func orderToModel(order db.Order) model.Order {
	result := model.Order{
		OrderNumber:    order.OrderNumber,
		Status:         model.OrderStatus(order.Status.String),
		TotalWeight:    orderNumericToFloat64(order.TotalWeight),
		TotalPrice:     orderNumericToFloat64(order.TotalPrice),
		PickupRequired: order.PickupRequired.Bool,
		Notes:          order.Notes.String,
	}
	if order.ID.Valid {
		result.ID = uuid.UUID(order.ID.Bytes).String()
	}
	if order.CustomerID.Valid {
		result.CustomerID = uuid.UUID(order.CustomerID.Bytes).String()
	}
	if order.ServiceID.Valid {
		result.ServiceID = uuid.UUID(order.ServiceID.Bytes).String()
	}
	if order.CreatedAt.Valid {
		result.CreatedAt = order.CreatedAt.Time.UTC().Format(time.RFC3339)
	}
	if order.UpdatedAt.Valid {
		result.UpdatedAt = order.UpdatedAt.Time.UTC().Format(time.RFC3339)
	}
	return result
}
