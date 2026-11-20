package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/nuryanto2121/laundryos/apps/api/internal/db/generated"
)

// MasterRepository provides access to master data (perfumes, discounts, payment methods, outlets, categories, employees).
type MasterRepository struct {
	queries *db.Queries
}

// NewMasterRepository creates a new MasterRepository.
func NewMasterRepository(pool *pgxpool.Pool) *MasterRepository {
	return &MasterRepository{queries: db.New(pool)}
}

func toPgUUID(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

// ListPerfumes retrieves all active perfumes.
func (r *MasterRepository) ListPerfumes(ctx context.Context) ([]db.Perfume, error) {
	perfumes, err := r.queries.ListActivePerfumes(ctx)
	if err != nil {
		return nil, fmt.Errorf("master repo: list perfumes: %w", err)
	}
	return perfumes, nil
}

// ListAllPerfumes retrieves all perfumes including inactive.
func (r *MasterRepository) ListAllPerfumes(ctx context.Context) ([]db.Perfume, error) {
	return r.queries.ListAllPerfumes(ctx)
}

// CreatePerfume creates a new perfume option.
func (r *MasterRepository) CreatePerfume(ctx context.Context, name, description string, isActive bool) (db.Perfume, error) {
	id := uuid.New()
	desc := pgtype.Text{}
	if description != "" {
		desc = pgtype.Text{String: description, Valid: true}
	}
	return r.queries.CreatePerfume(ctx, db.CreatePerfumeParams{
		ID:          toPgUUID(id),
		Name:        name,
		Description: desc,
		IsActive:    isActive,
	})
}

// UpdatePerfume updates an existing perfume.
func (r *MasterRepository) UpdatePerfume(ctx context.Context, id uuid.UUID, name, description string, isActive bool) (db.Perfume, error) {
	desc := pgtype.Text{}
	if description != "" {
		desc = pgtype.Text{String: description, Valid: true}
	}
	return r.queries.UpdatePerfume(ctx, db.UpdatePerfumeParams{
		ID:          toPgUUID(id),
		Name:        name,
		Description: desc,
		IsActive:    isActive,
	})
}

// DeletePerfume deactivates a perfume.
func (r *MasterRepository) DeletePerfume(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeletePerfume(ctx, toPgUUID(id))
}

// ListDiscounts retrieves all active discounts.
func (r *MasterRepository) ListDiscounts(ctx context.Context) ([]db.Discount, error) {
	discounts, err := r.queries.ListActiveDiscounts(ctx)
	if err != nil {
		return nil, fmt.Errorf("master repo: list discounts: %w", err)
	}
	return discounts, nil
}

// ListAllDiscounts retrieves all discounts.
func (r *MasterRepository) ListAllDiscounts(ctx context.Context) ([]db.Discount, error) {
	return r.queries.ListAllDiscounts(ctx)
}

// CreateDiscount creates a new discount.
func (r *MasterRepository) CreateDiscount(ctx context.Context, name, discountType string, value, minOrder float64, isActive bool) (db.Discount, error) {
	id := uuid.New()
	valNumeric := pgtype.Numeric{}
	_ = valNumeric.Scan(fmt.Sprintf("%.2f", value))
	minNumeric := pgtype.Numeric{}
	_ = minNumeric.Scan(fmt.Sprintf("%.2f", minOrder))

	return r.queries.CreateDiscount(ctx, db.CreateDiscountParams{
		ID:             toPgUUID(id),
		Name:           name,
		DiscountType:   discountType,
		Value:          valNumeric,
		MinOrderAmount: minNumeric,
		IsActive:       isActive,
	})
}

// UpdateDiscount updates an existing discount.
func (r *MasterRepository) UpdateDiscount(ctx context.Context, id uuid.UUID, name, discountType string, value, minOrder float64, isActive bool) (db.Discount, error) {
	valNumeric := pgtype.Numeric{}
	_ = valNumeric.Scan(fmt.Sprintf("%.2f", value))
	minNumeric := pgtype.Numeric{}
	_ = minNumeric.Scan(fmt.Sprintf("%.2f", minOrder))

	return r.queries.UpdateDiscount(ctx, db.UpdateDiscountParams{
		ID:             toPgUUID(id),
		Name:           name,
		DiscountType:   discountType,
		Value:          valNumeric,
		MinOrderAmount: minNumeric,
		IsActive:       isActive,
	})
}

// DeleteDiscount deactivates a discount.
func (r *MasterRepository) DeleteDiscount(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteDiscount(ctx, toPgUUID(id))
}

// ListPaymentMethods retrieves all active payment methods.
func (r *MasterRepository) ListPaymentMethods(ctx context.Context) ([]db.PaymentMethod, error) {
	methods, err := r.queries.ListActivePaymentMethods(ctx)
	if err != nil {
		return nil, fmt.Errorf("master repo: list payment methods: %w", err)
	}
	return methods, nil
}

// ListOutlets retrieves all active outlets.
func (r *MasterRepository) ListOutlets(ctx context.Context) ([]db.Outlet, error) {
	outlets, err := r.queries.ListOutlets(ctx)
	if err != nil {
		return nil, fmt.Errorf("master repo: list outlets: %w", err)
	}
	return outlets, nil
}

// GetDefaultOutlet retrieves the primary default outlet.
func (r *MasterRepository) GetDefaultOutlet(ctx context.Context) (db.Outlet, error) {
	outlet, err := r.queries.GetDefaultOutlet(ctx)
	if err != nil {
		return db.Outlet{}, fmt.Errorf("master repo: get default outlet: %w", err)
	}
	return outlet, nil
}

// ListExpenseCategories retrieves active expense categories.
func (r *MasterRepository) ListExpenseCategories(ctx context.Context) ([]db.ExpenseCategory, error) {
	return r.queries.ListExpenseCategories(ctx)
}

// ListAllExpenseCategories retrieves all expense categories with optional sorting.
func (r *MasterRepository) ListAllExpenseCategories(ctx context.Context, sortBy, sortDir string) ([]db.ExpenseCategory, error) {
	var sortByParam pgtype.Text
	if sortBy != "" {
		sortByParam = pgtype.Text{String: sortBy, Valid: true}
	}
	var sortDirParam pgtype.Text
	if sortDir != "" {
		sortDirParam = pgtype.Text{String: sortDir, Valid: true}
	}
	return r.queries.ListAllExpenseCategories(ctx, db.ListAllExpenseCategoriesParams{
		SortBy:  sortByParam,
		SortDir: sortDirParam,
	})
}

// CreateExpenseCategory creates a new expense category.
func (r *MasterRepository) CreateExpenseCategory(ctx context.Context, code, name string, isActive bool) (db.ExpenseCategory, error) {
	id := uuid.New()
	return r.queries.CreateExpenseCategory(ctx, db.CreateExpenseCategoryParams{
		ID:       toPgUUID(id),
		Code:     code,
		Name:     name,
		IsActive: isActive,
	})
}

// UpdateExpenseCategory updates an existing expense category.
func (r *MasterRepository) UpdateExpenseCategory(ctx context.Context, id uuid.UUID, name string, isActive bool) (db.ExpenseCategory, error) {
	return r.queries.UpdateExpenseCategory(ctx, db.UpdateExpenseCategoryParams{
		ID:       toPgUUID(id),
		Name:     name,
		IsActive: isActive,
	})
}

// DeleteExpenseCategory deactivates an expense category.
func (r *MasterRepository) DeleteExpenseCategory(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteExpenseCategory(ctx, toPgUUID(id))
}

// ListEmployees retrieves all active employees.
func (r *MasterRepository) ListEmployees(ctx context.Context) ([]db.Employee, error) {
	return r.queries.ListEmployees(ctx)
}

// ListAllEmployees retrieves all employees.
func (r *MasterRepository) ListAllEmployees(ctx context.Context) ([]db.Employee, error) {
	return r.queries.ListAllEmployees(ctx)
}

// CreateEmployee creates a new employee.
func (r *MasterRepository) CreateEmployee(ctx context.Context, username, fullName, phone, role, pin string, isActive bool) (db.Employee, error) {
	id := uuid.New()
	phoneText := pgtype.Text{}
	if phone != "" {
		phoneText = pgtype.Text{String: phone, Valid: true}
	}
	return r.queries.CreateEmployee(ctx, db.CreateEmployeeParams{
		ID:       toPgUUID(id),
		Username: username,
		FullName: fullName,
		Phone:    phoneText,
		Role:     role,
		Pin:      pin,
		IsActive: isActive,
	})
}

// UpdateEmployee updates an employee.
func (r *MasterRepository) UpdateEmployee(ctx context.Context, id uuid.UUID, fullName, phone, role, pin string, isActive bool) (db.Employee, error) {
	phoneText := pgtype.Text{}
	if phone != "" {
		phoneText = pgtype.Text{String: phone, Valid: true}
	}
	return r.queries.UpdateEmployee(ctx, db.UpdateEmployeeParams{
		ID:       toPgUUID(id),
		FullName: fullName,
		Phone:    phoneText,
		Role:     role,
		Pin:      pin,
		IsActive: isActive,
	})
}

// DeleteEmployee deactivates an employee.
func (r *MasterRepository) DeleteEmployee(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteEmployee(ctx, toPgUUID(id))
}

// ListServiceCategories retrieves active service categories.
func (r *MasterRepository) ListServiceCategories(ctx context.Context) ([]db.ServiceCategory, error) {
	return r.queries.ListActiveServiceCategories(ctx)
}

// ListAllServiceCategories retrieves all service categories.
func (r *MasterRepository) ListAllServiceCategories(ctx context.Context) ([]db.ServiceCategory, error) {
	return r.queries.ListAllServiceCategories(ctx)
}

// CreateServiceCategory creates a new service category.
func (r *MasterRepository) CreateServiceCategory(ctx context.Context, code, name, defaultUnit string, displayOrder int, isActive bool) (db.ServiceCategory, error) {
	id := uuid.New()
	return r.queries.CreateServiceCategory(ctx, db.CreateServiceCategoryParams{
		ID:           toPgUUID(id),
		Code:         code,
		Name:         name,
		DefaultUnit:  defaultUnit,
		DisplayOrder: pgtype.Int4{Int32: int32(displayOrder), Valid: true},
		IsActive:     pgtype.Bool{Bool: isActive, Valid: true},
	})
}

// UpdateServiceCategory updates an existing service category.
func (r *MasterRepository) UpdateServiceCategory(ctx context.Context, id uuid.UUID, name, defaultUnit string, displayOrder int, isActive bool) (db.ServiceCategory, error) {
	return r.queries.UpdateServiceCategory(ctx, db.UpdateServiceCategoryParams{
		ID:           toPgUUID(id),
		Name:         name,
		DefaultUnit:  defaultUnit,
		DisplayOrder: pgtype.Int4{Int32: int32(displayOrder), Valid: true},
		IsActive:     pgtype.Bool{Bool: isActive, Valid: true},
	})
}

// DeleteServiceCategory deactivates a service category.
func (r *MasterRepository) DeleteServiceCategory(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteServiceCategory(ctx, toPgUUID(id))
}

