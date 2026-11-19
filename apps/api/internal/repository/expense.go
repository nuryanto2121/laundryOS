package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/nuryanto2121/laundryos/apps/api/internal/db/generated"
	"github.com/nuryanto2121/laundryos/apps/api/internal/model"
)

// ExpenseRepository handles database persistence for operational expenses (kas kecil).
type ExpenseRepository struct {
	queries *db.Queries
}

// NewExpenseRepository creates a new ExpenseRepository.
func NewExpenseRepository(pool *pgxpool.Pool) *ExpenseRepository {
	return &ExpenseRepository{queries: db.New(pool)}
}

// Create inserts a new expense record.
func (r *ExpenseRepository) Create(ctx context.Context, req model.CreateExpenseRequest) (db.Expense, error) {
	id := uuid.New()
	var outletUUID pgtype.UUID
	if req.OutletID != "" {
		parsed, err := uuid.Parse(req.OutletID)
		if err == nil {
			outletUUID = pgtype.UUID{Bytes: parsed, Valid: true}
		}
	}

	expenseDate := time.Now()
	if req.ExpenseDate != "" {
		if t, err := time.Parse("2006-01-02", req.ExpenseDate); err == nil {
			expenseDate = t
		}
	}

	var amountNum pgtype.Numeric
	amountStr := strconv.FormatFloat(req.Amount, 'f', 2, 64)
	if err := amountNum.Scan(amountStr); err != nil {
		return db.Expense{}, fmt.Errorf("invalid amount: %w", err)
	}

	createdBy := req.CreatedBy
	if createdBy == "" {
		createdBy = "Kasir"
	}

	return r.queries.CreateExpense(ctx, db.CreateExpenseParams{
		ID:          pgtype.UUID{Bytes: id, Valid: true},
		OutletID:    outletUUID,
		ExpenseDate: pgtype.Date{Time: expenseDate, Valid: true},
		Category:    req.Category,
		Amount:      amountNum,
		Description: req.Description,
		CreatedBy:   createdBy,
	})
}

// List retrieves expenses between startDate and endDate.
func (r *ExpenseRepository) List(ctx context.Context, start, end time.Time) ([]db.Expense, error) {
	return r.queries.ListExpenses(ctx, db.ListExpensesParams{
		ExpenseDate:   pgtype.Date{Time: start, Valid: true},
		ExpenseDate_2: pgtype.Date{Time: end, Valid: true},
	})
}

// Get retrieves an expense by UUID string.
func (r *ExpenseRepository) Get(ctx context.Context, idStr string) (db.Expense, error) {
	parsed, err := uuid.Parse(idStr)
	if err != nil {
		return db.Expense{}, fmt.Errorf("invalid expense id: %w", err)
	}
	return r.queries.GetExpense(ctx, pgtype.UUID{Bytes: parsed, Valid: true})
}

// Delete removes an expense by UUID string.
func (r *ExpenseRepository) Delete(ctx context.Context, idStr string) error {
	parsed, err := uuid.Parse(idStr)
	if err != nil {
		return fmt.Errorf("invalid expense id: %w", err)
	}
	return r.queries.DeleteExpense(ctx, pgtype.UUID{Bytes: parsed, Valid: true})
}

// GetSummary returns expense totals grouped by category.
func (r *ExpenseRepository) GetSummary(ctx context.Context, start, end time.Time) ([]db.GetExpenseSummaryRow, error) {
	return r.queries.GetExpenseSummary(ctx, db.GetExpenseSummaryParams{
		ExpenseDate:   pgtype.Date{Time: start, Valid: true},
		ExpenseDate_2: pgtype.Date{Time: end, Valid: true},
	})
}

// GetTotal returns the total expense sum for a date range.
func (r *ExpenseRepository) GetTotal(ctx context.Context, start, end time.Time) (float64, error) {
	num, err := r.queries.GetTotalExpenseByDateRange(ctx, db.GetTotalExpenseByDateRangeParams{
		ExpenseDate:   pgtype.Date{Time: start, Valid: true},
		ExpenseDate_2: pgtype.Date{Time: end, Valid: true},
	})
	if err != nil {
		return 0, err
	}
	f, _ := num.Float64Value()
	return f.Float64, nil
}
