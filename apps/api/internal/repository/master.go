package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/nuryanto2121/laundryos/apps/api/internal/db/generated"
	"github.com/nuryanto2121/laundryos/apps/api/internal/model"
)

// MasterRepository provides access to master data (perfumes, discounts, payment methods, outlets, categories, employees).
type MasterRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

// NewMasterRepository creates a new MasterRepository.
func NewMasterRepository(pool *pgxpool.Pool) *MasterRepository {
	return &MasterRepository{pool: pool, queries: db.New(pool)}
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

// ListOutletsWithGPS retrieves all outlets with GPS coordinates and geofence radius.
func (r *MasterRepository) ListOutletsWithGPS(ctx context.Context) ([]model.Outlet, error) {
	if r.pool == nil {
		outlets, err := r.queries.ListOutlets(ctx)
		if err != nil {
			return nil, err
		}
		var res []model.Outlet
		for _, o := range outlets {
			res = append(res, model.Outlet{
				ID:                   o.ID.String(),
				Name:                 o.Name,
				Code:                 o.Code,
				Address:              o.Address.String,
				Phone:                o.Phone.String,
				Latitude:             -6.113425,
				Longitude:            106.914230,
				GeofenceRadiusMeters: 25,
				IsActive:             o.IsActive,
				CreatedAt:            o.CreatedAt.Time,
				UpdatedAt:            o.UpdatedAt.Time,
			})
		}
		return res, nil
	}

	query := `
		SELECT id, name, code, COALESCE(address, ''), COALESCE(phone, ''),
		       COALESCE(latitude, -6.113425), COALESCE(longitude, 106.914230),
		       COALESCE(geofence_radius_meters, 25), is_active, created_at, updated_at
		FROM outlets
		ORDER BY name ASC
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("master repo: list outlets with gps: %w", err)
	}
	defer rows.Close()

	var list []model.Outlet
	for rows.Next() {
		var o model.Outlet
		if err := rows.Scan(
			&o.ID, &o.Name, &o.Code, &o.Address, &o.Phone,
			&o.Latitude, &o.Longitude, &o.GeofenceRadiusMeters,
			&o.IsActive, &o.CreatedAt, &o.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, o)
	}
	return list, nil
}

// CreateOutlet creates a new outlet with GPS coordinates.
func (r *MasterRepository) CreateOutlet(ctx context.Context, req model.CreateOutletRequest) (model.Outlet, error) {
	lat := req.Latitude
	if lat == 0 {
		lat = -6.113425
	}
	lng := req.Longitude
	if lng == 0 {
		lng = 106.914230
	}
	radius := req.GeofenceRadiusMeters
	if radius <= 0 {
		radius = 25
	}

	query := `
		INSERT INTO outlets (name, code, address, phone, latitude, longitude, geofence_radius_meters, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, name, code, COALESCE(address, ''), COALESCE(phone, ''),
		          latitude, longitude, geofence_radius_meters, is_active, created_at, updated_at
	`
	var o model.Outlet
	err := r.pool.QueryRow(ctx, query,
		req.Name, req.Code, req.Address, req.Phone,
		lat, lng, radius, req.IsActive,
	).Scan(
		&o.ID, &o.Name, &o.Code, &o.Address, &o.Phone,
		&o.Latitude, &o.Longitude, &o.GeofenceRadiusMeters,
		&o.IsActive, &o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		return model.Outlet{}, fmt.Errorf("master repo: create outlet: %w", err)
	}
	return o, nil
}

// UpdateOutlet updates outlet details including GPS and geofence radius.
func (r *MasterRepository) UpdateOutlet(ctx context.Context, id string, req model.UpdateOutletRequest) (model.Outlet, error) {
	query := `
		UPDATE outlets
		SET name = COALESCE(NULLIF($2, ''), name),
		    code = COALESCE(NULLIF($3, ''), code),
		    address = COALESCE(NULLIF($4, ''), address),
		    phone = COALESCE(NULLIF($5, ''), phone),
		    latitude = CASE WHEN $6::numeric != 0::numeric THEN $6::numeric ELSE latitude END,
		    longitude = CASE WHEN $7::numeric != 0::numeric THEN $7::numeric ELSE longitude END,
		    geofence_radius_meters = CASE WHEN $8 > 0 THEN $8 ELSE geofence_radius_meters END,
		    is_active = $9,
		    updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, code, COALESCE(address, ''), COALESCE(phone, ''),
		          latitude, longitude, geofence_radius_meters, is_active, created_at, updated_at
	`
	var o model.Outlet
	err := r.pool.QueryRow(ctx, query,
		id, req.Name, req.Code, req.Address, req.Phone,
		req.Latitude, req.Longitude, req.GeofenceRadiusMeters, req.IsActive,
	).Scan(
		&o.ID, &o.Name, &o.Code, &o.Address, &o.Phone,
		&o.Latitude, &o.Longitude, &o.GeofenceRadiusMeters,
		&o.IsActive, &o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		return model.Outlet{}, fmt.Errorf("master repo: update outlet: %w", err)
	}
	return o, nil
}

// DeleteOutlet deactivates an outlet.
func (r *MasterRepository) DeleteOutlet(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `UPDATE outlets SET is_active = false, updated_at = NOW() WHERE id = $1`, id)
	return err
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

