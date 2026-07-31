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

// CustomerRepository handles persistence operations for customers.
// It is the only layer aware of sqlc and pgtype.
type CustomerRepository struct {
	queries *db.Queries
}

// NewCustomerRepository creates a new CustomerRepository.
func NewCustomerRepository(pool *pgxpool.Pool) *CustomerRepository {
	return &CustomerRepository{queries: db.New(pool)}
}

// Create inserts a new customer and returns the created record.
func (r *CustomerRepository) Create(ctx context.Context, req model.CreateCustomerRequest) (model.Customer, error) {
	id := uuid.New()

	arg := db.CreateCustomerParams{
		ID:           pgtype.UUID{Bytes: id, Valid: true},
		CustomerCode: req.CustomerCode,
		FullName:     req.FullName,
		Phone:        req.Phone,
		Email:        pgtype.Text{String: req.Email, Valid: req.Email != ""},
		Address:      pgtype.Text{String: req.Address, Valid: req.Address != ""},
		Latitude:     toNumeric(req.Latitude),
		Longitude:    toNumeric(req.Longitude),
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

// Update updates mutable customer fields and returns the updated record.
func (r *CustomerRepository) Update(ctx context.Context, id string, req model.UpdateCustomerRequest) (model.Customer, error) {
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
	}

	c, err := r.queries.UpdateCustomer(ctx, arg)
	if err != nil {
		return model.Customer{}, fmt.Errorf("repository: update customer: %w", err)
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

// toModel converts a sqlc-generated db.Customer to model.Customer.
func toModel(c db.Customer) model.Customer {
	m := model.Customer{
		CustomerCode: c.CustomerCode,
		FullName:     c.FullName,
		Phone:        c.Phone,
		Latitude:     numericToFloat64(c.Latitude),
		Longitude:    numericToFloat64(c.Longitude),
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
	return m
}
