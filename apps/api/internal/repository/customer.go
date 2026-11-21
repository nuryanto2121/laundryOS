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

// CustomerRepository handles persistence operations for customers.
// It is the only layer aware of sqlc and pgtype.
type CustomerRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

// NewCustomerRepository creates a new CustomerRepository.
func NewCustomerRepository(pool *pgxpool.Pool) *CustomerRepository {
	return &CustomerRepository{pool: pool, queries: db.New(pool)}
}

// Create inserts a new customer and returns the created record.
func (r *CustomerRepository) Create(ctx context.Context, req model.CreateCustomerRequest, customerCode string, status model.CustomerStatus) (model.Customer, error) {
	id := uuid.New()

	arg := db.CreateCustomerParams{
		ID:           pgtype.UUID{Bytes: id, Valid: true},
		CustomerCode: customerCode,
		FullName:     req.FullName,
		Phone:        req.Phone,
		Email:        pgtype.Text{String: req.Email, Valid: req.Email != ""},
		Address:      pgtype.Text{String: req.Address, Valid: req.Address != ""},
		Latitude:     toNumeric(req.Latitude),
		Longitude:    toNumeric(req.Longitude),
		Status:       string(status),
	}

	c, err := r.queries.CreateCustomer(ctx, arg)
	if err != nil {
		return model.Customer{}, fmt.Errorf("repository: create customer: %w", err)
	}
	return toModel(c), nil
}

// Get retrieves a single non-deleted customer by string ID.
func (r *CustomerRepository) Get(ctx context.Context, id string) (model.Customer, error) {
	pgID, err := parseUUID(id)
	if err != nil {
		return model.Customer{}, err
	}
	c, err := r.queries.GetCustomer(ctx, pgID)
	if err != nil {
		return model.Customer{}, fmt.Errorf("repository: get customer: %w", err)
	}
	return toModel(c), nil
}

// GetByPhone retrieves a customer by phone number across all lifecycle states (including deleted/inactive).
func (r *CustomerRepository) GetByPhone(ctx context.Context, phone string) (model.Customer, error) {
	c, err := r.queries.GetCustomerByPhone(ctx, phone)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Customer{}, nil
		}
		return model.Customer{}, fmt.Errorf("repository: get customer by phone: %w", err)
	}
	return toModel(c), nil
}

// GetLatestCustomerCodeToday retrieves today's latest customer code based on created_at.
// Note: encountering pgx.ErrNoRows is a normal operational condition indicating the first customer registration of the day, not an error.
func (r *CustomerRepository) GetLatestCustomerCodeToday(ctx context.Context) (string, error) {
	code, err := r.queries.GetLatestCustomerCodeToday(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", fmt.Errorf("repository: get latest customer code today: %w", err)
	}
	return code, nil
}

// List returns all non-deleted customers.
func (r *CustomerRepository) List(ctx context.Context) ([]model.Customer, error) {
	rows, err := r.queries.ListCustomers(ctx)
	if err != nil {
		return nil, fmt.Errorf("repository: list customers: %w", err)
	}
	result := make([]model.Customer, len(rows))
	for i, c := range rows {
		result[i] = toModel(c)
	}
	return result, nil
}

// ListPaginated returns customers with search filter, pagination, and total count.
func (r *CustomerRepository) ListPaginated(ctx context.Context, search, sortBy, sortDir string, limit, offset int) ([]model.Customer, int64, error) {
	var searchParam pgtype.Text
	if search != "" {
		searchParam = pgtype.Text{String: search, Valid: true}
	}
	var sortByParam pgtype.Text
	if sortBy != "" {
		sortByParam = pgtype.Text{String: sortBy, Valid: true}
	}
	var sortDirParam pgtype.Text
	if sortDir != "" {
		sortDirParam = pgtype.Text{String: sortDir, Valid: true}
	}

	rows, err := r.queries.ListCustomersPaginated(ctx, db.ListCustomersPaginatedParams{
		Limit:   int32(limit),
		Offset:  int32(offset),
		Search:  searchParam,
		SortBy:  sortByParam,
		SortDir: sortDirParam,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repository: list customers paginated: %w", err)
	}

	count, err := r.queries.CountCustomers(ctx, searchParam)
	if err != nil {
		return nil, 0, fmt.Errorf("repository: count customers: %w", err)
	}

	result := make([]model.Customer, len(rows))
	for i, c := range rows {
		result[i] = toModel(c)
	}
	return result, count, nil
}

// Update updates customer fields and returns the updated record.
func (r *CustomerRepository) Update(ctx context.Context, id string, req model.UpdateCustomerRequest, status model.CustomerStatus, deletedAt *string) (model.Customer, error) {
	pgID, err := parseUUID(id)
	if err != nil {
		return model.Customer{}, err
	}

	arg := db.UpdateCustomerParams{
		ID:        pgID,
		FullName:  req.FullName,
		Phone:     req.Phone,
		Email:     pgtype.Text{String: req.Email, Valid: req.Email != ""},
		Address:   pgtype.Text{String: req.Address, Valid: req.Address != ""},
		Latitude:  toNumeric(req.Latitude),
		Longitude: toNumeric(req.Longitude),
		Status:    string(status),
		DeletedAt: parseTimestamp(deletedAt),
	}

	c, err := r.queries.UpdateCustomer(ctx, arg)
	if err != nil {
		return model.Customer{}, fmt.Errorf("repository: update customer: %w", err)
	}
	return toModel(c), nil
}

// UpdateStatus persists a status change for an existing customer.
func (r *CustomerRepository) UpdateStatus(ctx context.Context, id string, status model.CustomerStatus) (model.Customer, error) {
	pgID, err := parseUUID(id)
	if err != nil {
		return model.Customer{}, err
	}

	arg := db.UpdateCustomerStatusParams{
		ID:     pgID,
		Status: string(status),
	}

	c, err := r.queries.UpdateCustomerStatus(ctx, arg)
	if err != nil {
		return model.Customer{}, fmt.Errorf("repository: update customer status: %w", err)
	}
	return toModel(c), nil
}

// Delete soft-deletes a customer by string ID.
func (r *CustomerRepository) Delete(ctx context.Context, id string) error {
	pgID, err := parseUUID(id)
	if err != nil {
		return err
	}
	if err := r.queries.DeleteCustomer(ctx, pgID); err != nil {
		return fmt.Errorf("repository: delete customer: %w", err)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Private mapping helpers — all pgtype conversions are isolated here
// ---------------------------------------------------------------------------

// parseUUID converts a string UUID to pgtype.UUID.
func parseUUID(id string) (pgtype.UUID, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("invalid customer id: %w", err)
	}
	return pgtype.UUID{Bytes: parsed, Valid: true}, nil
}

// toNumeric converts *float64 to pgtype.Numeric.
// pgtype.Numeric.Scan only accepts string or []byte — not float64 directly.
func toNumeric(f *float64) pgtype.Numeric {
	if f == nil {
		return pgtype.Numeric{}
	}
	var n pgtype.Numeric
	_ = n.Scan(strconv.FormatFloat(*f, 'f', -1, 64))
	return n
}

// numericToFloat64 converts pgtype.Numeric to *float64.
func numericToFloat64(n pgtype.Numeric) *float64 {
	if !n.Valid {
		return nil
	}
	f, _ := n.Float64Value()
	if !f.Valid {
		return nil
	}
	v := f.Float64
	return &v
}

// parseTimestamp converts a formatted string pointer to pgtype.Timestamp.
func parseTimestamp(s *string) pgtype.Timestamp {
	if s == nil || *s == "" {
		return pgtype.Timestamp{}
	}
	t, err := time.Parse(time.RFC3339, *s)
	if err != nil {
		return pgtype.Timestamp{}
	}
	return pgtype.Timestamp{Time: t.UTC(), Valid: true}
}

// toModel converts a sqlc-generated db.Customer to model.Customer.
func toModel(c db.Customer) model.Customer {
	m := model.Customer{
		CustomerCode: c.CustomerCode,
		FullName:     c.FullName,
		Phone:        c.Phone,
		Latitude:     numericToFloat64(c.Latitude),
		Longitude:    numericToFloat64(c.Longitude),
		Status:       model.CustomerStatus(c.Status),
	}
	if c.ID.Valid {
		m.ID = uuid.UUID(c.ID.Bytes).String()
	}
	if c.Email.Valid {
		m.Email = c.Email.String
	}
	if c.Address.Valid {
		m.Address = c.Address.String
	}
	if c.CreatedAt.Valid {
		m.CreatedAt = c.CreatedAt.Time.UTC().Format(time.RFC3339)
	}
	if c.UpdatedAt.Valid {
		m.UpdatedAt = c.UpdatedAt.Time.UTC().Format(time.RFC3339)
	}
	if c.DeletedAt.Valid {
		s := c.DeletedAt.Time.UTC().Format(time.RFC3339)
		m.DeletedAt = &s
	}
	return m
}

// GetCustomerDefaultOutletID returns the customer's preferred outlet_id, or from their latest order.
func (r *CustomerRepository) GetCustomerDefaultOutletID(ctx context.Context, customerID string) string {
	var outletID string
	_ = r.pool.QueryRow(ctx, `SELECT COALESCE(outlet_id::text, '') FROM customers WHERE id = $1`, customerID).Scan(&outletID)
	if outletID != "" {
		return outletID
	}
	_ = r.pool.QueryRow(ctx, `SELECT COALESCE(outlet_id::text, '') FROM orders WHERE customer_id = $1 ORDER BY created_at DESC LIMIT 1`, customerID).Scan(&outletID)
	return outletID
}

// SetCustomerOutletID updates the customer's preferred outlet.
func (r *CustomerRepository) SetCustomerOutletID(ctx context.Context, customerID string, outletID string) error {
	_, err := r.pool.Exec(ctx, `UPDATE customers SET outlet_id = $2, updated_at = NOW() WHERE id = $1`, customerID, outletID)
	return err
}
