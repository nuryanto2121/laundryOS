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

// ServiceRepository handles persistence operations for laundry services.
type ServiceRepository struct {
	queries *db.Queries
}

// NewServiceRepository creates a new ServiceRepository.
func NewServiceRepository(pool *pgxpool.Pool) *ServiceRepository {
	return &ServiceRepository{queries: db.New(pool)}
}

// ListActiveServices retrieves all currently active services from the database.
func (r *ServiceRepository) ListActiveServices(ctx context.Context) ([]model.ServiceItem, error) {
	services, err := r.queries.ListActiveServices(ctx)
	if err != nil {
		return nil, fmt.Errorf("repository: list active services: %w", err)
	}

	items := make([]model.ServiceItem, len(services))
	for i, s := range services {
		var catID string
		if s.CategoryID.Valid {
			catID = uuid.UUID(s.CategoryID.Bytes).String()
		}
		items[i] = model.ServiceItem{
			ID:           uuid.UUID(s.ID.Bytes).String(),
			ServiceCode:  s.ServiceCode,
			ServiceName:  s.ServiceName,
			DurationDay:  int(s.DurationDay.Int32),
			PricePerKg:   serviceNumericToFloat64(s.PricePerKg),
			Price:        serviceNumericToFloat64(s.PricePerKg),
			Active:       s.Active.Bool,
			IsActive:     s.Active.Bool,
			CategoryID:   catID,
			CategoryName: s.CategoryName.String,
			Category:     s.CategoryName.String,
			Unit:         s.Unit.String,
		}
	}
	return items, nil
}

// ListAllServices retrieves all services with optional sorting.
func (r *ServiceRepository) ListAllServices(ctx context.Context, sortBy, sortDir string) ([]model.ServiceItem, error) {
	var sortByParam pgtype.Text
	if sortBy != "" {
		sortByParam = pgtype.Text{String: sortBy, Valid: true}
	}
	var sortDirParam pgtype.Text
	if sortDir != "" {
		sortDirParam = pgtype.Text{String: sortDir, Valid: true}
	}

	services, err := r.queries.ListAllServices(ctx, db.ListAllServicesParams{
		SortBy:  sortByParam,
		SortDir: sortDirParam,
	})
	if err != nil {
		return nil, fmt.Errorf("repository: list all services: %w", err)
	}

	items := make([]model.ServiceItem, len(services))
	for i, s := range services {
		var catID string
		if s.CategoryID.Valid {
			catID = uuid.UUID(s.CategoryID.Bytes).String()
		}
		items[i] = model.ServiceItem{
			ID:           uuid.UUID(s.ID.Bytes).String(),
			ServiceCode:  s.ServiceCode,
			ServiceName:  s.ServiceName,
			DurationDay:  int(s.DurationDay.Int32),
			PricePerKg:   serviceNumericToFloat64(s.PricePerKg),
			Price:        serviceNumericToFloat64(s.PricePerKg),
			Active:       s.Active.Bool,
			IsActive:     s.Active.Bool,
			CategoryID:   catID,
			CategoryName: s.CategoryName.String,
			Category:     s.CategoryName.String,
			Unit:         s.Unit.String,
		}
	}
	return items, nil
}

// CreateService creates a new service in database.
func (r *ServiceRepository) CreateService(ctx context.Context, code, name string, duration int, price float64, active bool, categoryID string, categoryName, unit string) (model.ServiceItem, error) {
	id := uuid.New()
	var idBytes [16]byte
	copy(idBytes[:], id[:])

	priceNumeric := pgtype.Numeric{}
	_ = priceNumeric.Scan(fmt.Sprintf("%.2f", price))

	var catUUID pgtype.UUID
	if u, err := uuid.Parse(categoryID); err == nil {
		catUUID = pgtype.UUID{Bytes: u, Valid: true}
	}

	created, err := r.queries.CreateService(ctx, db.CreateServiceParams{
		ID:           pgtype.UUID{Bytes: idBytes, Valid: true},
		ServiceCode:  code,
		ServiceName:  name,
		DurationDay:  pgtype.Int4{Int32: int32(duration), Valid: true},
		PricePerKg:   priceNumeric,
		Active:       pgtype.Bool{Bool: active, Valid: true},
		CategoryID:   catUUID,
		CategoryName: pgtype.Text{String: categoryName, Valid: categoryName != ""},
		Unit:         pgtype.Text{String: unit, Valid: unit != ""},
	})
	if err != nil {
		return model.ServiceItem{}, err
	}

	return model.ServiceItem{
		ID:           id.String(),
		ServiceCode:  created.ServiceCode,
		ServiceName:  created.ServiceName,
		DurationDay:  int(created.DurationDay.Int32),
		PricePerKg:   price,
		Price:        price,
		Active:       created.Active.Bool,
		IsActive:     created.Active.Bool,
		CategoryID:   categoryID,
		CategoryName: created.CategoryName.String,
		Category:     created.CategoryName.String,
		Unit:         created.Unit.String,
	}, nil
}

// UpdateService updates an existing service.
func (r *ServiceRepository) UpdateService(ctx context.Context, id uuid.UUID, name string, duration int, price float64, active bool, categoryID string, categoryName, unit string) (model.ServiceItem, error) {
	var idBytes [16]byte
	copy(idBytes[:], id[:])

	priceNumeric := pgtype.Numeric{}
	_ = priceNumeric.Scan(fmt.Sprintf("%.2f", price))

	var catUUID pgtype.UUID
	if u, err := uuid.Parse(categoryID); err == nil {
		catUUID = pgtype.UUID{Bytes: u, Valid: true}
	}

	updated, err := r.queries.UpdateService(ctx, db.UpdateServiceParams{
		ID:           pgtype.UUID{Bytes: idBytes, Valid: true},
		ServiceName:  name,
		DurationDay:  pgtype.Int4{Int32: int32(duration), Valid: true},
		PricePerKg:   priceNumeric,
		Active:       pgtype.Bool{Bool: active, Valid: true},
		CategoryID:   catUUID,
		CategoryName: pgtype.Text{String: categoryName, Valid: categoryName != ""},
		Unit:         pgtype.Text{String: unit, Valid: unit != ""},
	})
	if err != nil {
		return model.ServiceItem{}, err
	}

	return model.ServiceItem{
		ID:           id.String(),
		ServiceCode:  updated.ServiceCode,
		ServiceName:  updated.ServiceName,
		DurationDay:  int(updated.DurationDay.Int32),
		PricePerKg:   price,
		Price:        price,
		Active:       updated.Active.Bool,
		IsActive:     updated.Active.Bool,
		CategoryID:   categoryID,
		CategoryName: updated.CategoryName.String,
		Category:     updated.CategoryName.String,
		Unit:         updated.Unit.String,
	}, nil
}

// DeleteService deactivates a service.
func (r *ServiceRepository) DeleteService(ctx context.Context, id uuid.UUID) error {
	var idBytes [16]byte
	copy(idBytes[:], id[:])
	return r.queries.DeleteService(ctx, pgtype.UUID{Bytes: idBytes, Valid: true})
}

func serviceNumericToFloat64(value pgtype.Numeric) float64 {
	if !value.Valid {
		return 0
	}
	result, err := value.Float64Value()
	if err != nil || !result.Valid {
		return 0
	}
	return result.Float64
}
