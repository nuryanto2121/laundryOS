package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	db "github.com/nuryanto2121/laundryos/apps/api/internal/db/generated"
	"github.com/nuryanto2121/laundryos/apps/api/internal/model"
	"github.com/nuryanto2121/laundryos/apps/api/internal/repository"
	"github.com/nuryanto2121/laundryos/apps/api/internal/util"
)

// POSService provides business logic for the Web POS cashier.
type POSService struct {
	posRepo      *repository.POSRepository
	orderRepo    *repository.OrderRepository
	customerRepo *repository.CustomerRepository
	customerSvc  *CustomerService
	notifSvc     NotificationService
}

// NewPOSService creates a new POSService.
func NewPOSService(
	posRepo *repository.POSRepository,
	orderRepo *repository.OrderRepository,
	customerRepo *repository.CustomerRepository,
	customerSvc *CustomerService,
	notifSvc NotificationService,
) *POSService {
	return &POSService{
		posRepo:      posRepo,
		orderRepo:    orderRepo,
		customerRepo: customerRepo,
		customerSvc:  customerSvc,
		notifSvc:     notifSvc,
	}
}

// CreateOrder processes a new walk-in POS order, resolving or creating customer,
// generating unique order number, saving items, and optionally triggering a WhatsApp invoice/receipt.
func (s *POSService) CreateOrder(ctx context.Context, req model.CreatePOSOrderRequest) (db.Order, []db.OrderItem, error) {
	var customerID uuid.UUID
	var customerPhone string
	var customerName string

	// 1. Resolve Customer
	if req.CustomerID != "" {
		if u, err := uuid.Parse(req.CustomerID); err == nil {
			customerID = u
			if c, err := s.customerRepo.Get(ctx, req.CustomerID); err == nil {
				customerPhone = c.Phone
				customerName = c.FullName
			}
		}
	}

	if customerID == uuid.Nil {
		cleanPhone := strings.TrimSpace(req.CustomerPhone)
		if cleanPhone == "" {
			return db.Order{}, nil, fmt.Errorf("nomor telepon atau customer_id wajib diisi")
		}

		existing, err := s.customerRepo.GetByPhone(ctx, cleanPhone)
		if err == nil && existing.ID != "" {
			customerID = uuid.MustParse(existing.ID)
			customerPhone = existing.Phone
			customerName = existing.FullName
		} else {
			// Register new customer
			name := strings.TrimSpace(req.CustomerName)
			if name == "" {
				name = "Pelanggan Walk-in"
			}
			newCust, err := s.customerSvc.Create(ctx, model.CreateCustomerRequest{
				FullName: name,
				Phone:    cleanPhone,
				Address:  req.CustomerAddress,
			})
			if err != nil {
				return db.Order{}, nil, fmt.Errorf("failed to create walk-in customer: %w", err)
			}
			customerID = uuid.MustParse(newCust.ID)
			customerPhone = newCust.Phone
			customerName = newCust.FullName
		}
	}

	// 2. Generate Order Number
	latestNumber, err := s.orderRepo.GetLatestOrderNumberToday(ctx)
	if err != nil {
		return db.Order{}, nil, fmt.Errorf("get latest order number: %w", err)
	}
	sequence, err := nextOrderSequence(latestNumber, time.Now())
	if err != nil {
		return db.Order{}, nil, err
	}
	orderNumber := util.GenerateOrderNumber(sequence)

	// 3. Persist Order in POS Repo
	order, items, err := s.posRepo.CreatePOSOrder(ctx, customerID, orderNumber, req)
	if err != nil {
		return db.Order{}, nil, err
	}

	// 4. Trigger Customer WA Notification if phone exists (async, fail-soft)
	if customerPhone != "" && s.notifSvc != nil {
		go func(phone, name, orderNum string, its []model.POSOrderItemRequest, price, paid float64, est *time.Time) {
			bgCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			var itemNames []string
			for _, it := range its {
				itemNames = append(itemNames, fmt.Sprintf("%s (%.1f %s)", it.ServiceName, it.Quantity, it.Unit))
			}
			summary := strings.Join(itemNames, ", ")

			estStr := "-"
			if est != nil {
				estStr = est.Format("02-01-2006 15:04 WIB")
			}

			_ = s.notifSvc.NotifyCustomerOrderCreated(bgCtx, phone, name, orderNum, summary, price, paid, estStr)
		}(customerPhone, customerName, orderNumber, req.Items, req.TotalPrice, req.PaidAmount, req.EstimatedCompletionAt)
	}

	return order, items, nil
}

// ListOrders returns filtered orders for the POS order list table.
func (s *POSService) ListOrders(ctx context.Context, status, search string) ([]repository.POSOrderListItem, error) {
	return s.posRepo.ListPOSOrders(ctx, status, search)
}

// GetOrderDetail returns detailed order with items, photos, and audit trail.
func (s *POSService) GetOrderDetail(ctx context.Context, id string) (repository.POSOrderDetail, error) {
	return s.posRepo.GetPOSOrderDetail(ctx, id)
}

// UpdatePayment settles or updates partial/full payment on an order.
func (s *POSService) UpdatePayment(ctx context.Context, id string, req model.UpdatePaymentRequest) (db.Order, error) {
	return s.posRepo.UpdatePayment(ctx, id, req)
}

// UpdateStatus moves the order to the next stage (e.g. ANTRIAN -> PROSES -> SIAP AMBIL -> SELESAI).
func (s *POSService) UpdateStatus(ctx context.Context, id string, req model.UpdateStatusRequest) (db.Order, error) {
	order, err := s.posRepo.UpdateStatus(ctx, id, req)
	if err != nil {
		return db.Order{}, err
	}

	// Notify customer via WhatsApp
	if s.notifSvc != nil {
		go func(orderID string, newStatus string) {
			bgCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			detail, err := s.posRepo.GetPOSOrderDetail(bgCtx, orderID)
			if err != nil {
				return
			}
			phone := detail.Order.CustomerPhone
			name := detail.Order.CustomerName
			orderNum := detail.Order.OrderNumber
			totalPrice, _ := detail.Order.TotalPrice.Float64Value()
			paidAmount, _ := detail.Order.PaidAmount.Float64Value()
			rackNum := detail.Order.RackNumber.String

			_ = s.notifSvc.NotifyCustomerOrderStatus(bgCtx, phone, name, orderNum, newStatus, rackNum, totalPrice.Float64, paidAmount.Float64)
		}(id, req.Status)
	}

	return order, nil
}

// CancelOrder marks an order as canceled with reason.
func (s *POSService) CancelOrder(ctx context.Context, id string, req model.CancelOrderRequest) (db.Order, error) {
	return s.posRepo.CancelOrder(ctx, id, req)
}

// UpdateRack sets the storage rack number.
func (s *POSService) UpdateRack(ctx context.Context, id string, rack string) (db.Order, error) {
	return s.posRepo.UpdateRack(ctx, id, rack)
}

// AddPhoto attaches documentation photos.
func (s *POSService) AddPhoto(ctx context.Context, id string, req model.AddOrderPhotoRequest) (db.OrderPhoto, error) {
	return s.posRepo.AddPhoto(ctx, id, req)
}

// ResendCustomerInvoice triggers manual WhatsApp invoice dispatch.
func (s *POSService) ResendCustomerInvoice(ctx context.Context, id string) error {
	if s.notifSvc == nil {
		return fmt.Errorf("notification service unavailable")
	}

	detail, err := s.posRepo.GetPOSOrderDetail(ctx, id)
	if err != nil {
		return err
	}

	phone := detail.Order.CustomerPhone
	name := detail.Order.CustomerName
	orderNum := detail.Order.OrderNumber
	totalPrice, _ := detail.Order.TotalPrice.Float64Value()
	paidAmount, _ := detail.Order.PaidAmount.Float64Value()

	var itemSummaries []string
	for _, it := range detail.Items {
		q, _ := it.Quantity.Float64Value()
		itemSummaries = append(itemSummaries, fmt.Sprintf("%s (%.1f %s)", it.ServiceName, q.Float64, it.Unit))
	}

	estStr := "-"
	if detail.Order.EstimatedCompletionAt.Valid {
		estStr = detail.Order.EstimatedCompletionAt.Time.Format("02-01-2006 15:04 WIB")
	}

	slog.Info("pos: resending customer invoice via whatsapp", "phone", phone, "order", orderNum)
	return s.notifSvc.NotifyCustomerOrderCreated(ctx, phone, name, orderNum, strings.Join(itemSummaries, ", "), totalPrice.Float64, paidAmount.Float64, estStr)
}
