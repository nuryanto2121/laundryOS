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

// MockCustomerRepository is a mock implementation of service.CustomerRepository.
type MockCustomerRepository struct {
	mock.Mock
}

func (m *MockCustomerRepository) Create(ctx context.Context, req model.CreateCustomerRequest, customerCode string, status model.CustomerStatus) (model.Customer, error) {
	args := m.Called(ctx, req, customerCode, status)
	return args.Get(0).(model.Customer), args.Error(1)
}

func (m *MockCustomerRepository) Get(ctx context.Context, id string) (model.Customer, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Customer), args.Error(1)
}

func (m *MockCustomerRepository) GetByPhone(ctx context.Context, phone string) (model.Customer, error) {
	args := m.Called(ctx, phone)
	return args.Get(0).(model.Customer), args.Error(1)
}

func (m *MockCustomerRepository) GetLatestCustomerCodeToday(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m *MockCustomerRepository) List(ctx context.Context) ([]model.Customer, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Customer), args.Error(1)
}

func (m *MockCustomerRepository) ListPaginated(ctx context.Context, search, sortBy, sortDir string, limit, offset int) ([]model.Customer, int64, error) {
	args := m.Called(ctx, search, sortBy, sortDir, limit, offset)
	return args.Get(0).([]model.Customer), args.Get(1).(int64), args.Error(2)
}

func (m *MockCustomerRepository) Update(ctx context.Context, id string, req model.UpdateCustomerRequest, status model.CustomerStatus, deletedAt *string) (model.Customer, error) {
	args := m.Called(ctx, id, req, status, deletedAt)
	return args.Get(0).(model.Customer), args.Error(1)
}

func (m *MockCustomerRepository) UpdateStatus(ctx context.Context, id string, status model.CustomerStatus) (model.Customer, error) {
	args := m.Called(ctx, id, status)
	return args.Get(0).(model.Customer), args.Error(1)
}

func (m *MockCustomerRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCustomerService_Create_CaseA_NewCustomer(t *testing.T) {
	req := model.CreateCustomerRequest{
		FullName: "John Doe",
		Phone:    "08123456789",
		Email:    "john@example.com",
		Address:  "Jl. Merad",
	}
	dateStr := time.Now().Format("060102")
	expectedCode := fmt.Sprintf("CST%s000001", dateStr)

	createdCust := model.Customer{
		ID:           "uuid-123",
		CustomerCode: expectedCode,
		FullName:     req.FullName,
		Phone:        req.Phone,
		Status:       model.CustomerStatusActive,
	}

	mockRepo := new(MockCustomerRepository)
	mockRepo.On("GetByPhone", mock.Anything, req.Phone).Return(model.Customer{}, nil)
	// Normal condition: first registration of the day returns empty string without error
	mockRepo.On("GetLatestCustomerCodeToday", mock.Anything).Return("", nil)
	mockRepo.On("Create", mock.Anything, req, expectedCode, model.CustomerStatusActive).Return(createdCust, nil)

	svc := service.NewCustomerService(mockRepo)
	res, err := svc.Create(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, createdCust, res)
	mockRepo.AssertExpectations(t)
}

func TestCustomerService_Create_CaseB_DuplicateActive(t *testing.T) {
	req := model.CreateCustomerRequest{
		FullName: "Jane Doe",
		Phone:    "08123456789",
	}
	existing := model.Customer{
		ID:     "existing-uuid",
		Phone:  req.Phone,
		Status: model.CustomerStatusActive,
	}

	mockRepo := new(MockCustomerRepository)
	mockRepo.On("GetByPhone", mock.Anything, req.Phone).Return(existing, nil)

	svc := service.NewCustomerService(mockRepo)
	_, err := svc.Create(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, "customer already registered", err.Error())
	mockRepo.AssertNotCalled(t, "Create")
	mockRepo.AssertNotCalled(t, "Update")
}

func TestCustomerService_Create_CaseC_Reactivation(t *testing.T) {
	req := model.CreateCustomerRequest{
		FullName: "Reactivated User",
		Phone:    "08123456789",
		Address:  "New Address",
	}
	existing := model.Customer{
		ID:           "existing-uuid",
		CustomerCode: "CST260101000001",
		FullName:     "Old Name",
		Phone:        req.Phone,
		Status:       model.CustomerStatusInactive,
		CreatedAt:    "2026-01-01T00:00:00Z",
	}

	updateReq := model.UpdateCustomerRequest{
		FullName: req.FullName,
		Phone:    req.Phone,
		Address:  req.Address,
	}
	reactivated := model.Customer{
		ID:           existing.ID,
		CustomerCode: existing.CustomerCode,
		FullName:     req.FullName,
		Phone:        req.Phone,
		Address:      req.Address,
		Status:       model.CustomerStatusActive,
		CreatedAt:    existing.CreatedAt,
		DeletedAt:    nil,
	}

	mockRepo := new(MockCustomerRepository)
	mockRepo.On("GetByPhone", mock.Anything, req.Phone).Return(existing, nil)
	mockRepo.On("Update", mock.Anything, existing.ID, updateReq, model.CustomerStatusActive, (*string)(nil)).Return(reactivated, nil)

	svc := service.NewCustomerService(mockRepo)
	res, err := svc.Create(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, reactivated, res)
	assert.Equal(t, existing.CustomerCode, res.CustomerCode)
	assert.Equal(t, model.CustomerStatusActive, res.Status)
	assert.Nil(t, res.DeletedAt)
	mockRepo.AssertNotCalled(t, "Create")
	mockRepo.AssertExpectations(t)
}

func TestCustomerService_Create_CaseD_Blocked(t *testing.T) {
	req := model.CreateCustomerRequest{
		FullName: "Blocked User",
		Phone:    "08123456789",
	}
	existing := model.Customer{
		ID:     "existing-uuid",
		Phone:  req.Phone,
		Status: model.CustomerStatusBlocked,
	}

	mockRepo := new(MockCustomerRepository)
	mockRepo.On("GetByPhone", mock.Anything, req.Phone).Return(existing, nil)

	svc := service.NewCustomerService(mockRepo)
	_, err := svc.Create(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, "customer is blocked", err.Error())
	mockRepo.AssertNotCalled(t, "Create")
	mockRepo.AssertNotCalled(t, "Update")
}

func TestCustomerService_DefensiveSequenceParsing(t *testing.T) {
	req := model.CreateCustomerRequest{FullName: "Defensive User", Phone: "08999999999"}
	dateStr := time.Now().Format("060102")

	tests := []struct {
		name         string
		latestCode   string
		latestErr    error
		expectedCode string
	}{
		{"empty string normal condition", "", nil, fmt.Sprintf("CST%s000001", dateStr)},
		{"malformed string", "INVALID_CODE", nil, fmt.Sprintf("CST%s000001", dateStr)},
		{"shorter than 6 characters", "CST1", nil, fmt.Sprintf("CST%s000001", dateStr)},
		{"unparseable trailing chars", fmt.Sprintf("CST%sABCDEF", dateStr), nil, fmt.Sprintf("CST%s000001", dateStr)},
		{"valid previous code 42", fmt.Sprintf("CST%s000042", dateStr), nil, fmt.Sprintf("CST%s000043", dateStr)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockCustomerRepository)
			mockRepo.On("GetByPhone", mock.Anything, req.Phone).Return(model.Customer{}, nil)
			mockRepo.On("GetLatestCustomerCodeToday", mock.Anything).Return(tt.latestCode, tt.latestErr)

			created := model.Customer{
				ID:           "uuid-gen",
				CustomerCode: tt.expectedCode,
				Status:       model.CustomerStatusActive,
			}
			mockRepo.On("Create", mock.Anything, req, tt.expectedCode, model.CustomerStatusActive).Return(created, nil)

			svc := service.NewCustomerService(mockRepo)
			res, err := svc.Create(context.Background(), req)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, res.CustomerCode)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCustomerService_Update_Immutability(t *testing.T) {
	id := "uuid-123"
	req := model.UpdateCustomerRequest{
		FullName: "New Name",
		Phone:    "08111111111",
	}
	deletedAtVal := "2026-05-01T00:00:00Z"
	existing := model.Customer{
		ID:           id,
		CustomerCode: "CST260101000001",
		FullName:     "Old Name",
		Phone:        "08000000000",
		Status:       model.CustomerStatusInactive,
		CreatedAt:    "2026-01-01T00:00:00Z",
		DeletedAt:    &deletedAtVal,
	}
	updated := model.Customer{
		ID:           id,
		CustomerCode: existing.CustomerCode,
		FullName:     req.FullName,
		Phone:        req.Phone,
		Status:       existing.Status,
		CreatedAt:    existing.CreatedAt,
		DeletedAt:    existing.DeletedAt,
	}

	mockRepo := new(MockCustomerRepository)
	mockRepo.On("Get", mock.Anything, id).Return(existing, nil)
	mockRepo.On("Update", mock.Anything, id, req, existing.Status, existing.DeletedAt).Return(updated, nil)

	svc := service.NewCustomerService(mockRepo)
	res, err := svc.Update(context.Background(), id, req)

	assert.NoError(t, err)
	assert.Equal(t, updated, res)
	assert.Equal(t, existing.CustomerCode, res.CustomerCode)
	assert.Equal(t, existing.Status, res.Status)
	assert.Equal(t, existing.DeletedAt, res.DeletedAt)
	mockRepo.AssertExpectations(t)
}

func TestCustomerService_Delete(t *testing.T) {
	mockRepo := new(MockCustomerRepository)
	mockRepo.On("Delete", mock.Anything, "uuid-123").Return(nil)

	svc := service.NewCustomerService(mockRepo)
	err := svc.Delete(context.Background(), "uuid-123")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCustomerService_UpdateStatus_ValidTransitions(t *testing.T) {
	id := "uuid-status"
	validTransitions := []struct {
		current model.CustomerStatus
		target  model.CustomerStatus
	}{
		{model.CustomerStatusActive, model.CustomerStatusInactive},
		{model.CustomerStatusActive, model.CustomerStatusBlocked},
		{model.CustomerStatusInactive, model.CustomerStatusActive},
		{model.CustomerStatusInactive, model.CustomerStatusBlocked},
		{model.CustomerStatusBlocked, model.CustomerStatusActive},
		{model.CustomerStatusBlocked, model.CustomerStatusInactive},
	}

	for _, tt := range validTransitions {
		t.Run(fmt.Sprintf("%s_to_%s", tt.current, tt.target), func(t *testing.T) {
			existing := model.Customer{ID: id, Status: tt.current}
			expected := model.Customer{ID: id, Status: tt.target}

			mockRepo := new(MockCustomerRepository)
			mockRepo.On("Get", mock.Anything, id).Return(existing, nil)
			mockRepo.On("UpdateStatus", mock.Anything, id, tt.target).Return(expected, nil)

			svc := service.NewCustomerService(mockRepo)
			res, err := svc.UpdateStatus(context.Background(), id, tt.target)

			assert.NoError(t, err)
			assert.Equal(t, tt.target, res.Status)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCustomerService_UpdateStatus_RedundantTransitions(t *testing.T) {
	id := "uuid-redundant"
	redundantTransitions := []model.CustomerStatus{
		model.CustomerStatusActive,
		model.CustomerStatusInactive,
		model.CustomerStatusBlocked,
	}

	for _, status := range redundantTransitions {
		t.Run(fmt.Sprintf("redundant_%s", status), func(t *testing.T) {
			existing := model.Customer{ID: id, Status: status}

			mockRepo := new(MockCustomerRepository)
			mockRepo.On("Get", mock.Anything, id).Return(existing, nil)

			svc := service.NewCustomerService(mockRepo)
			_, err := svc.UpdateStatus(context.Background(), id, status)

			assert.Error(t, err)
			assert.Equal(t, fmt.Sprintf("customer already %s", status), err.Error())
			mockRepo.AssertNotCalled(t, "UpdateStatus")
			mockRepo.AssertNotCalled(t, "Update")
		})
	}
}

func TestCustomerService_Get_and_List(t *testing.T) {
	mockRepo := new(MockCustomerRepository)
	cust := model.Customer{ID: "id-1", FullName: "John"}
	mockRepo.On("Get", mock.Anything, "id-1").Return(cust, nil)
	mockRepo.On("List", mock.Anything).Return([]model.Customer{cust}, nil)

	svc := service.NewCustomerService(mockRepo)

	res, err := svc.Get(context.Background(), "id-1")
	assert.NoError(t, err)
	assert.Equal(t, cust, res)

	listRes, err := svc.List(context.Background())
	assert.NoError(t, err)
	assert.Len(t, listRes, 1)
	assert.Equal(t, cust, listRes[0])
	mockRepo.AssertExpectations(t)
}
