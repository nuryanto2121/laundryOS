package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

	newID := uuid.New()
	query := `
		INSERT INTO outlets (id, name, code, address, phone, latitude, longitude, geofence_radius_meters, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, name, code, COALESCE(address, ''), COALESCE(phone, ''),
		          latitude, longitude, geofence_radius_meters, is_active, created_at, updated_at
	`
	var o model.Outlet
	err := r.pool.QueryRow(ctx, query,
		newID, req.Name, req.Code, req.Address, req.Phone,
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

// Authenticate checks user credentials and returns AuthUser.
func (r *MasterRepository) Authenticate(ctx context.Context, username, password string) (model.AuthUser, error) {
	if r.pool == nil {
		return model.AuthUser{}, errors.New("database pool not initialized")
	}

	query := `
		SELECT e.id, e.username, e.full_name, COALESCE(e.phone, ''), COALESCE(e.email, ''),
		       e.role, COALESCE(e.password, e.pin), COALESCE(e.outlet_id::text, ''),
		       COALESCE(o.name, ''), COALESCE(e.permissions, '{}'::jsonb), e.is_active
		FROM employees e
		LEFT JOIN outlets o ON o.id = e.outlet_id
		WHERE LOWER(e.username) = LOWER($1) OR (e.phone IS NOT NULL AND e.phone = $1)
		LIMIT 1
	`
	var (
		user       model.AuthUser
		storedPass string
		permBytes  []byte
		isActive   bool
	)
	err := r.pool.QueryRow(ctx, query, strings.TrimSpace(username)).Scan(
		&user.ID, &user.Username, &user.FullName, &user.Phone, &user.Email,
		&user.Role, &storedPass, &user.OutletID,
		&user.OutletName, &permBytes, &isActive,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.AuthUser{}, errors.New("username atau nomor handphone tidak ditemukan")
		}
		return model.AuthUser{}, fmt.Errorf("master repo: auth: %w", err)
	}

	if !isActive {
		return model.AuthUser{}, errors.New("akun pegawai nonaktif, silakan hubungi owner")
	}

	if password != storedPass {
		return model.AuthUser{}, errors.New("password salah")
	}

	user.Permissions = make(map[string]bool)
	if len(permBytes) > 0 {
		_ = json.Unmarshal(permBytes, &user.Permissions)
	}

	return user, nil
}

// ListEmployeesDetailed retrieves all employees with full profile and ACL permissions.
func (r *MasterRepository) ListEmployeesDetailed(ctx context.Context) ([]model.Employee, error) {
	if r.pool == nil {
		// fallback to basic query
		emps, err := r.queries.ListEmployees(ctx)
		if err != nil {
			return nil, err
		}
		var list []model.Employee
		for _, e := range emps {
			list = append(list, model.Employee{
				ID:        e.ID.String(),
				Username:  e.Username,
				FullName:  e.FullName,
				Phone:     e.Phone.String,
				Role:      e.Role,
				IsActive:  e.IsActive,
				CreatedAt: e.CreatedAt.Time,
				UpdatedAt: e.UpdatedAt.Time,
			})
		}
		return list, nil
	}

	query := `
		SELECT e.id, e.username, e.full_name, COALESCE(e.phone, ''), COALESCE(e.email, ''),
		       e.role, COALESCE(e.outlet_id::text, ''), COALESCE(o.name, ''),
		       COALESCE(e.permissions, '{}'::jsonb), e.is_active, e.created_at, e.updated_at
		FROM employees e
		LEFT JOIN outlets o ON o.id = e.outlet_id
		ORDER BY e.created_at ASC
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("master repo: list employees detailed: %w", err)
	}
	defer rows.Close()

	var list []model.Employee
	for rows.Next() {
		var (
			emp       model.Employee
			permBytes []byte
		)
		if err := rows.Scan(
			&emp.ID, &emp.Username, &emp.FullName, &emp.Phone, &emp.Email,
			&emp.Role, &emp.OutletID, &emp.OutletName,
			&permBytes, &emp.IsActive, &emp.CreatedAt, &emp.UpdatedAt,
		); err != nil {
			return nil, err
		}
		emp.Permissions = make(map[string]bool)
		if len(permBytes) > 0 {
			_ = json.Unmarshal(permBytes, &emp.Permissions)
		}
		list = append(list, emp)
	}
	return list, nil
}

// CreateEmployeeDetailed creates an employee with password, outlet, and ACL permissions.
func (r *MasterRepository) CreateEmployeeDetailed(ctx context.Context, req model.CreateEmployeeRequest) (model.Employee, error) {
	if r.pool == nil {
		return model.Employee{}, errors.New("database pool not initialized")
	}

	id := uuid.New()
	password := req.Password
	if password == "" {
		password = "password123"
	}
	role := req.Role
	if role == "" {
		role = "KASIR"
	}

	permBytes, _ := json.Marshal(req.Permissions)
	if len(req.Permissions) == 0 {
		permBytes = []byte("{}")
	}

	var outletID *uuid.UUID
	if req.OutletID != "" {
		if parsed, err := uuid.Parse(req.OutletID); err == nil {
			outletID = &parsed
		}
	}

	query := `
		INSERT INTO employees (id, username, full_name, phone, email, role, password, pin, outlet_id, permissions, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, username, full_name, COALESCE(phone, ''), COALESCE(email, ''), role,
		          COALESCE(outlet_id::text, ''), permissions, is_active, created_at, updated_at
	`
	var (
		emp          model.Employee
		retPermBytes []byte
	)
	err := r.pool.QueryRow(ctx, query,
		id, req.Username, req.FullName, req.Phone, req.Email,
		role, password, password, outletID, permBytes, req.IsActive,
	).Scan(
		&emp.ID, &emp.Username, &emp.FullName, &emp.Phone, &emp.Email, &emp.Role,
		&emp.OutletID, &retPermBytes, &emp.IsActive, &emp.CreatedAt, &emp.UpdatedAt,
	)
	if err != nil {
		return model.Employee{}, fmt.Errorf("master repo: create employee: %w", err)
	}
	emp.Permissions = make(map[string]bool)
	_ = json.Unmarshal(retPermBytes, &emp.Permissions)
	return emp, nil
}

// UpdateEmployeeDetailed updates employee details, password, outlet, and permissions.
func (r *MasterRepository) UpdateEmployeeDetailed(ctx context.Context, idStr string, req model.UpdateEmployeeRequest) (model.Employee, error) {
	if r.pool == nil {
		return model.Employee{}, errors.New("database pool not initialized")
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return model.Employee{}, errors.New("id pegawai tidak valid")
	}

	var outletID *uuid.UUID
	if req.OutletID != "" {
		if parsed, err := uuid.Parse(req.OutletID); err == nil {
			outletID = &parsed
		}
	}

	permBytes, _ := json.Marshal(req.Permissions)
	if len(req.Permissions) == 0 {
		permBytes = []byte("{}")
	}

	query := `
		UPDATE employees
		SET full_name = COALESCE(NULLIF($2, ''), full_name),
		    phone = COALESCE(NULLIF($3, ''), phone),
		    email = COALESCE(NULLIF($4, ''), email),
		    role = COALESCE(NULLIF($5, ''), role),
		    password = CASE WHEN $6 != '' THEN $6 ELSE password END,
		    pin = CASE WHEN $6 != '' THEN $6 ELSE pin END,
		    outlet_id = $7,
		    permissions = $8,
		    is_active = $9,
		    updated_at = NOW()
		WHERE id = $1
		RETURNING id, username, full_name, COALESCE(phone, ''), COALESCE(email, ''), role,
		          COALESCE(outlet_id::text, ''), permissions, is_active, created_at, updated_at
	`
	var (
		emp          model.Employee
		retPermBytes []byte
	)
	err = r.pool.QueryRow(ctx, query,
		id, req.FullName, req.Phone, req.Email, req.Role,
		req.Password, outletID, permBytes, req.IsActive,
	).Scan(
		&emp.ID, &emp.Username, &emp.FullName, &emp.Phone, &emp.Email, &emp.Role,
		&emp.OutletID, &retPermBytes, &emp.IsActive, &emp.CreatedAt, &emp.UpdatedAt,
	)
	if err != nil {
		return model.Employee{}, fmt.Errorf("master repo: update employee: %w", err)
	}
	emp.Permissions = make(map[string]bool)
	_ = json.Unmarshal(retPermBytes, &emp.Permissions)
	return emp, nil
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

