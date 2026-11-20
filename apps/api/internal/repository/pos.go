package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/nuryanto2121/laundryos/apps/api/internal/db/generated"
	"github.com/nuryanto2121/laundryos/apps/api/internal/model"
)

// POSRepository handles all Point of Sale transactional and reporting persistence.
type POSRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

// NewPOSRepository creates a new POSRepository.
func NewPOSRepository(pool *pgxpool.Pool) *POSRepository {
	return &POSRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

// POSOrderDetail holds full aggregate details for an order.
type POSOrderDetail struct {
	Order   db.GetPOSOrderDetailRow `json:"order"`
	Items   []db.OrderItem          `json:"items"`
	Photos  []db.OrderPhoto         `json:"photos"`
	History []db.OrderStatusHistory `json:"history"`
}

// CreatePOSOrder creates an order, its items, and an initial audit history entry.
func (r *POSRepository) CreatePOSOrder(
	ctx context.Context,
	customerID uuid.UUID,
	orderNumber string,
	req model.CreatePOSOrderRequest,
) (db.Order, []db.OrderItem, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return db.Order{}, nil, fmt.Errorf("pos repo: begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	orderID := uuid.New()

	var outletUUID pgtype.UUID
	if req.OutletID != "" {
		if u, err := uuid.Parse(req.OutletID); err == nil {
			outletUUID = pgtype.UUID{Bytes: u, Valid: true}
		}
	}

	var perfumeUUID pgtype.UUID
	if req.PerfumeID != "" {
		if u, err := uuid.Parse(req.PerfumeID); err == nil {
			perfumeUUID = pgtype.UUID{Bytes: u, Valid: true}
		}
	}

	var discountUUID pgtype.UUID
	if req.DiscountID != "" {
		if u, err := uuid.Parse(req.DiscountID); err == nil {
			discountUUID = pgtype.UUID{Bytes: u, Valid: true}
		}
	}

	var paymentMethodUUID pgtype.UUID
	if req.PaymentMethodID != "" {
		if u, err := uuid.Parse(req.PaymentMethodID); err == nil {
			paymentMethodUUID = pgtype.UUID{Bytes: u, Valid: true}
		}
	}

	var primaryServiceUUID pgtype.UUID
	if len(req.Items) > 0 && req.Items[0].ServiceID != "" {
		if u, err := uuid.Parse(req.Items[0].ServiceID); err == nil {
			primaryServiceUUID = pgtype.UUID{Bytes: u, Valid: true}
		}
	}

	paymentStatus := "UNPAID"
	if req.PaidAmount >= req.TotalPrice && req.TotalPrice > 0 {
		paymentStatus = "PAID"
	} else if req.PaidAmount > 0 {
		paymentStatus = "PARTIAL"
	}

	var paidAt pgtype.Timestamp
	if req.PaidAmount > 0 {
		paidAt = pgtype.Timestamp{Time: time.Now(), Valid: true}
	}

	var estCompletion pgtype.Timestamp
	if req.EstimatedCompletionAt != nil {
		estCompletion = pgtype.Timestamp{Time: *req.EstimatedCompletionAt, Valid: true}
	}

	totalWeightNum := toOrderNumeric(req.TotalWeight)
	totalPriceNum := toOrderNumeric(req.TotalPrice)
	discountAmountNum := toOrderNumeric(req.DiscountAmount)
	paidAmountNum := toOrderNumeric(req.PaidAmount)

	order, err := qtx.CreatePOSOrder(ctx, db.CreatePOSOrderParams{
		ID:                    pgtype.UUID{Bytes: orderID, Valid: true},
		OrderNumber:           orderNumber,
		CustomerID:            pgtype.UUID{Bytes: customerID, Valid: true},
		ServiceID:             primaryServiceUUID,
		Status:                pgtype.Text{String: "ANTRIAN", Valid: true},
		TotalWeight:           totalWeightNum,
		TotalPrice:            totalPriceNum,
		PickupRequired:        pgtype.Bool{Bool: req.PickupRequired, Valid: true},
		Notes:                 pgtype.Text{String: req.Notes, Valid: req.Notes != ""},
		OutletID:              outletUUID,
		RackNumber:            pgtype.Text{String: req.RackNumber, Valid: req.RackNumber != ""},
		PerfumeID:             perfumeUUID,
		PerfumeName:           pgtype.Text{String: req.PerfumeName, Valid: req.PerfumeName != ""},
		DiscountID:            discountUUID,
		DiscountAmount:        discountAmountNum,
		PaidAmount:            paidAmountNum,
		PaymentStatus:         paymentStatus,
		PaymentMethodID:       paymentMethodUUID,
		PaymentMethodName:     pgtype.Text{String: req.PaymentMethodName, Valid: req.PaymentMethodName != ""},
		PaidAt:                paidAt,
		EstimatedCompletionAt: estCompletion,
	})
	if err != nil {
		return db.Order{}, nil, fmt.Errorf("pos repo: create order: %w", err)
	}

	// Insert order items
	var createdItems []db.OrderItem
	for _, it := range req.Items {
		var itemServiceUUID pgtype.UUID
		if it.ServiceID != "" {
			if u, err := uuid.Parse(it.ServiceID); err == nil {
				itemServiceUUID = pgtype.UUID{Bytes: u, Valid: true}
			}
		}

		item, err := qtx.CreateOrderItem(ctx, db.CreateOrderItemParams{
			ID:           pgtype.UUID{Bytes: uuid.New(), Valid: true},
			OrderID:      pgtype.UUID{Bytes: orderID, Valid: true},
			ServiceID:    itemServiceUUID,
			ServiceName:  it.ServiceName,
			Category:     it.Category,
			Quantity:     toOrderNumeric(it.Quantity),
			Unit:         it.Unit,
			PricePerUnit: toOrderNumeric(it.PricePerUnit),
			Subtotal:     toOrderNumeric(it.Subtotal),
			Notes:        pgtype.Text{String: it.Notes, Valid: it.Notes != ""},
		})
		if err != nil {
			return db.Order{}, nil, fmt.Errorf("pos repo: create order item: %w", err)
		}
		createdItems = append(createdItems, item)
	}

	// Initial Audit Trail
	createdBy := req.CreatedBy
	if createdBy == "" {
		createdBy = "Kasir"
	}
	_, err = qtx.CreateOrderStatusHistory(ctx, db.CreateOrderStatusHistoryParams{
		ID:             pgtype.UUID{Bytes: uuid.New(), Valid: true},
		OrderID:        pgtype.UUID{Bytes: orderID, Valid: true},
		PreviousStatus: pgtype.Text{Valid: false},
		NewStatus:      "ANTRIAN",
		ChangedBy:      createdBy,
		Notes:          pgtype.Text{String: "Order dibuat dari Walk-in POS", Valid: true},
	})
	if err != nil {
		return db.Order{}, nil, fmt.Errorf("pos repo: audit trail: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return db.Order{}, nil, fmt.Errorf("pos repo: commit: %w", err)
	}

	return order, createdItems, nil
}

// ListPOSOrders lists orders with status and search query filter.
func (r *POSRepository) ListPOSOrders(ctx context.Context, status, search string) ([]db.ListPOSOrdersRow, error) {
	var statusParam pgtype.Text
	if status != "" && status != "SEMUA" {
		statusParam = pgtype.Text{String: status, Valid: true}
	}
	var searchParam pgtype.Text
	if search != "" {
		searchParam = pgtype.Text{String: search, Valid: true}
	}

	return r.queries.ListPOSOrders(ctx, db.ListPOSOrdersParams{
		Status: statusParam,
		Search: searchParam,
	})
}

// GetPOSOrderDetail retrieves detailed order with items, photos, and history.
func (r *POSRepository) GetPOSOrderDetail(ctx context.Context, idStr string) (POSOrderDetail, error) {
	orderUUID, err := uuid.Parse(idStr)
	if err != nil {
		return POSOrderDetail{}, fmt.Errorf("invalid order id: %w", err)
	}

	orderRow, err := r.queries.GetPOSOrderDetail(ctx, pgtype.UUID{Bytes: orderUUID, Valid: true})
	if err != nil {
		return POSOrderDetail{}, fmt.Errorf("pos repo: get order: %w", err)
	}

	items, err := r.queries.ListOrderItemsByOrderId(ctx, pgtype.UUID{Bytes: orderUUID, Valid: true})
	if err != nil {
		items = []db.OrderItem{}
	}

	photos, err := r.queries.ListOrderPhotosByOrderId(ctx, pgtype.UUID{Bytes: orderUUID, Valid: true})
	if err != nil {
		photos = []db.OrderPhoto{}
	}

	history, err := r.queries.ListOrderStatusHistoryByOrderId(ctx, pgtype.UUID{Bytes: orderUUID, Valid: true})
	if err != nil {
		history = []db.OrderStatusHistory{}
	}

	return POSOrderDetail{
		Order:   orderRow,
		Items:   items,
		Photos:  photos,
		History: history,
	}, nil
}

// UpdatePayment settles or updates the payment for an order.
func (r *POSRepository) UpdatePayment(ctx context.Context, idStr string, req model.UpdatePaymentRequest) (db.Order, error) {
	orderUUID, err := uuid.Parse(idStr)
	if err != nil {
		return db.Order{}, fmt.Errorf("invalid order id: %w", err)
	}

	order, err := r.queries.GetOrder(ctx, pgtype.UUID{Bytes: orderUUID, Valid: true})
	if err != nil {
		return db.Order{}, fmt.Errorf("get order: %w", err)
	}

	var paymentMethodUUID pgtype.UUID
	if req.PaymentMethodID != "" {
		if u, err := uuid.Parse(req.PaymentMethodID); err == nil {
			paymentMethodUUID = pgtype.UUID{Bytes: u, Valid: true}
		}
	}

	totalPrice, _ := order.TotalPrice.Float64Value()
	newPaidAmount := req.PaidAmount
	paymentStatus := "UNPAID"
	if newPaidAmount >= totalPrice.Float64 && totalPrice.Float64 > 0 {
		paymentStatus = "PAID"
	} else if newPaidAmount > 0 {
		paymentStatus = "PARTIAL"
	}

	updated, err := r.queries.UpdatePOSOrderPayment(ctx, db.UpdatePOSOrderPaymentParams{
		ID:                pgtype.UUID{Bytes: orderUUID, Valid: true},
		PaidAmount:        toOrderNumeric(newPaidAmount),
		PaymentStatus:     paymentStatus,
		PaymentMethodID:   paymentMethodUUID,
		PaymentMethodName: pgtype.Text{String: req.PaymentMethodName, Valid: req.PaymentMethodName != ""},
		PaidAt:            pgtype.Timestamp{Time: time.Now(), Valid: true},
	})
	if err != nil {
		return db.Order{}, err
	}

	staff := req.StaffName
	if staff == "" {
		staff = "Kasir"
	}
	r.queries.CreateOrderStatusHistory(ctx, db.CreateOrderStatusHistoryParams{
		ID:             pgtype.UUID{Bytes: uuid.New(), Valid: true},
		OrderID:        pgtype.UUID{Bytes: orderUUID, Valid: true},
		PreviousStatus: order.Status,
		NewStatus:      order.Status.String,
		ChangedBy:      staff,
		Notes:          pgtype.Text{String: fmt.Sprintf("Pembayaran dicatat: Rp %.0f (%s)", newPaidAmount, paymentStatus), Valid: true},
	})

	return updated, nil
}

// UpdateStatus changes the order status and records the audit history.
func (r *POSRepository) UpdateStatus(ctx context.Context, idStr string, req model.UpdateStatusRequest) (db.Order, error) {
	orderUUID, err := uuid.Parse(idStr)
	if err != nil {
		return db.Order{}, fmt.Errorf("invalid order id: %w", err)
	}

	order, err := r.queries.GetOrder(ctx, pgtype.UUID{Bytes: orderUUID, Valid: true})
	if err != nil {
		return db.Order{}, fmt.Errorf("get order: %w", err)
	}

	updated, err := r.queries.UpdatePOSOrderStatus(ctx, db.UpdatePOSOrderStatusParams{
		ID:     pgtype.UUID{Bytes: orderUUID, Valid: true},
		Status: req.Status,
	})
	if err != nil {
		return db.Order{}, err
	}

	changedBy := req.ChangedBy
	if changedBy == "" {
		changedBy = "Kasir"
	}

	r.queries.CreateOrderStatusHistory(ctx, db.CreateOrderStatusHistoryParams{
		ID:             pgtype.UUID{Bytes: uuid.New(), Valid: true},
		OrderID:        pgtype.UUID{Bytes: orderUUID, Valid: true},
		PreviousStatus: order.Status,
		NewStatus:      req.Status,
		ChangedBy:      changedBy,
		Notes:          pgtype.Text{String: req.Notes, Valid: req.Notes != ""},
	})

	return updated, nil
}

// CancelOrder marks an order as canceled.
func (r *POSRepository) CancelOrder(ctx context.Context, idStr string, req model.CancelOrderRequest) (db.Order, error) {
	orderUUID, err := uuid.Parse(idStr)
	if err != nil {
		return db.Order{}, fmt.Errorf("invalid order id: %w", err)
	}

	order, err := r.queries.GetOrder(ctx, pgtype.UUID{Bytes: orderUUID, Valid: true})
	if err != nil {
		return db.Order{}, fmt.Errorf("get order: %w", err)
	}

	updated, err := r.queries.CancelPOSOrder(ctx, db.CancelPOSOrderParams{
		ID:             pgtype.UUID{Bytes: orderUUID, Valid: true},
		CanceledReason: pgtype.Text{String: req.Reason, Valid: req.Reason != ""},
	})
	if err != nil {
		return db.Order{}, err
	}

	changedBy := req.ChangedBy
	if changedBy == "" {
		changedBy = "Kasir"
	}

	r.queries.CreateOrderStatusHistory(ctx, db.CreateOrderStatusHistoryParams{
		ID:             pgtype.UUID{Bytes: uuid.New(), Valid: true},
		OrderID:        pgtype.UUID{Bytes: orderUUID, Valid: true},
		PreviousStatus: order.Status,
		NewStatus:      "BATAL",
		ChangedBy:      changedBy,
		Notes:          pgtype.Text{String: "Pembatalan: " + req.Reason, Valid: true},
	})

	return updated, nil
}

// UpdateRack sets the storage rack for an order.
func (r *POSRepository) UpdateRack(ctx context.Context, idStr, rack string) (db.Order, error) {
	orderUUID, err := uuid.Parse(idStr)
	if err != nil {
		return db.Order{}, fmt.Errorf("invalid order id: %w", err)
	}

	return r.queries.UpdatePOSOrderRack(ctx, db.UpdatePOSOrderRackParams{
		ID:         pgtype.UUID{Bytes: orderUUID, Valid: true},
		RackNumber: pgtype.Text{String: rack, Valid: rack != ""},
	})
}

// AddPhoto attaches a photo documentation record.
func (r *POSRepository) AddPhoto(ctx context.Context, idStr string, req model.AddOrderPhotoRequest) (db.OrderPhoto, error) {
	orderUUID, err := uuid.Parse(idStr)
	if err != nil {
		return db.OrderPhoto{}, fmt.Errorf("invalid order id: %w", err)
	}

	return r.queries.CreateOrderPhoto(ctx, db.CreateOrderPhotoParams{
		ID:        pgtype.UUID{Bytes: uuid.New(), Valid: true},
		OrderID:   pgtype.UUID{Bytes: orderUUID, Valid: true},
		PhotoType: req.PhotoType,
		PhotoUrl:  req.PhotoURL,
		Notes:     pgtype.Text{String: req.Notes, Valid: req.Notes != ""},
	})
}
