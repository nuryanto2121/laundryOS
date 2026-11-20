package service

import (
	"context"
	"time"

	db "github.com/nuryanto2121/laundryos/apps/api/internal/db/generated"
	"github.com/nuryanto2121/laundryos/apps/api/internal/model"
	"github.com/nuryanto2121/laundryos/apps/api/internal/repository"
)

// ExpenseService manages operational expenses (petty cash).
type ExpenseService struct {
	repo *repository.ExpenseRepository
}

// NewExpenseService creates a new ExpenseService.
func NewExpenseService(repo *repository.ExpenseRepository) *ExpenseService {
	return &ExpenseService{repo: repo}
}

// Create records a new operational expense.
func (s *ExpenseService) Create(ctx context.Context, req model.CreateExpenseRequest) (db.Expense, error) {
	return s.repo.Create(ctx, req)
}

// Update corrects an existing operational expense.
func (s *ExpenseService) Update(ctx context.Context, id string, req model.UpdateExpenseRequest) (db.Expense, error) {
	return s.repo.Update(ctx, id, req)
}

// List retrieves expenses between a given date range.
func (s *ExpenseService) List(ctx context.Context, start, end time.Time) ([]db.Expense, error) {
	return s.repo.List(ctx, start, end)
}

// ListPaginated retrieves expenses with date range, category filter, and returns a paginated envelope.
func (s *ExpenseService) ListPaginated(ctx context.Context, start, end time.Time, category string, page, limit int) (model.PaginatedResponse, error) {
	offset, meta := model.CalculatePagination(page, limit, 0)
	expenses, total, err := s.repo.ListPaginated(ctx, start, end, category, meta.Limit, offset)
	if err != nil {
		return model.PaginatedResponse{}, err
	}
	_, meta = model.CalculatePagination(page, limit, total)
	return model.PaginatedResponse{
		Data: expenses,
		Meta: meta,
	}, nil
}

// Get retrieves a single expense by ID.
func (s *ExpenseService) Get(ctx context.Context, id string) (db.Expense, error) {
	return s.repo.Get(ctx, id)
}

// Delete removes an expense by ID.
func (s *ExpenseService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// GetSummary returns expense totals grouped by category.
func (s *ExpenseService) GetSummary(ctx context.Context, start, end time.Time) ([]db.GetExpenseSummaryRow, error) {
	return s.repo.GetSummary(ctx, start, end)
}

// GetTotal returns the total expense sum in a date range.
func (s *ExpenseService) GetTotal(ctx context.Context, start, end time.Time) (float64, error) {
	return s.repo.GetTotal(ctx, start, end)
}
