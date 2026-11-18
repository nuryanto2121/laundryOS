package service_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/nuryanto2121/laundryos/apps/api/internal/model"
	"github.com/nuryanto2121/laundryos/apps/api/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Create(ctx context.Context, req model.CreateOrderRequest, orderNumber string, status model.OrderStatus) (model.Order, error) {
	args := m.Called(ctx, req, orderNumber, status)
	return args.Get(0).(model.Order), args.Error(1)
}

func (m *MockOrderRepository) Get(ctx context.Context, id string) (model.Order, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Order), args.Error(1)
}

func (m *MockOrderRepository) GetLatestOrderNumberToday(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m *MockOrderRepository) GetLatestByCustomerPhone(ctx context.Context, phone string) (model.OrderTrackingDetail, error) {
	args := m.Called(ctx, phone)
	return args.Get(0).(model.OrderTrackingDetail), args.Error(1)
}

func (m *MockOrderRepository) GetByOrderNumber(ctx context.Context, orderNumber string) (model.OrderTrackingDetail, error) {
	args := m.Called(ctx, orderNumber)
	return args.Get(0).(model.OrderTrackingDetail), args.Error(1)
}

func (m *MockOrderRepository) List(ctx context.Context) ([]model.Order, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Order), args.Error(1)
}

func (m *MockOrderRepository) Update(ctx context.Context, id string, req model.UpdateOrderRequest) (model.Order, error) {
	args := m.Called(ctx, id, req)
	return args.Get(0).(model.Order), args.Error(1)
}

func (m *MockOrderRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestOrderService_CreateDeterministicNumber(t *testing.T) {
	req := model.CreateOrderRequest{CustomerID: "customer-id", ServiceID: "service-id"}
	dateStr := time.Now().Format("060102")
	latest := fmt.Sprintf("ORD%s000042", dateStr)
	expected := fmt.Sprintf("ORD%s000043", dateStr)
	created := model.Order{OrderNumber: expected, Status: model.OrderStatusNew}

	repo := new(MockOrderRepository)
	repo.On("GetLatestOrderNumberToday", mock.Anything).Return(latest, nil)
	repo.On("Create", mock.Anything, req, expected, model.OrderStatusNew).Return(created, nil)

	result, err := service.NewOrderService(repo).Create(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, created, result)
	repo.AssertExpectations(t)
}

func TestOrderService_CreateStartsAtOneWhenNoPreviousOrder(t *testing.T) {
	req := model.CreateOrderRequest{}
	dateStr := time.Now().Format("060102")
	expected := fmt.Sprintf("ORD%s000001", dateStr)

	repo := new(MockOrderRepository)
	repo.On("GetLatestOrderNumberToday", mock.Anything).Return("", nil)
	repo.On("Create", mock.Anything, req, expected, model.OrderStatusNew).Return(model.Order{OrderNumber: expected}, nil)

	_, err := service.NewOrderService(repo).Create(context.Background(), req)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestOrderService_CreateRejectsMalformedNumber(t *testing.T) {
	repo := new(MockOrderRepository)
	repo.On("GetLatestOrderNumberToday", mock.Anything).Return("INVALID_CODE", nil)

	_, err := service.NewOrderService(repo).Create(context.Background(), model.CreateOrderRequest{})

	assert.EqualError(t, err, "service: malformed latest order number")
	repo.AssertNotCalled(t, "Create")
}

func TestOrderService_UpdateDelegatesMutableFields(t *testing.T) {
	req := model.UpdateOrderRequest{CustomerID: "customer-id", ServiceID: "service-id", TotalWeight: 4}
	expected := model.Order{ID: "order-id", OrderNumber: "ORD260810000001", Status: model.OrderStatusNew}

	repo := new(MockOrderRepository)
	repo.On("Update", mock.Anything, "order-id", req).Return(expected, nil)

	result, err := service.NewOrderService(repo).Update(context.Background(), "order-id", req)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestOrderService_DeleteDelegatesToRepository(t *testing.T) {
	repo := new(MockOrderRepository)
	repo.On("Delete", mock.Anything, "order-id").Return(nil)

	err := service.NewOrderService(repo).Delete(context.Background(), "order-id")

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}
