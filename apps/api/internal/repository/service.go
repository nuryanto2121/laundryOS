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
		items[i] = model.ServiceItem{
			ID:          uuid.UUID(s.ID.Bytes).String(),
			ServiceCode: s.ServiceCode,
			ServiceName: s.ServiceName,
			DurationDay: int(s.DurationDay.Int32),
			PricePerKg:  serviceNumericToFloat64(s.PricePerKg),
			Active:      s.Active.Bool,
		}
	}
	return items, nil
}

// ListAllServices retrieves all services.
func (r *ServiceRepository) ListAllServices(ctx context.Context) ([]model.ServiceItem, error) {
	services, err := r.queries.ListAllServices(ctx)
	if err != nil {
		return nil, fmt.Errorf("repository: list all services: %w", err)
	}

	items := make([]model.ServiceItem, len(services))
	for i, s := range services {
		items[i] = model.ServiceItem{
			ID:          uuid.UUID(s.ID.Bytes).String(),
			ServiceCode: s.ServiceCode,
			ServiceName: s.ServiceName,
			DurationDay: int(s.DurationDay.Int32),
			PricePerKg:  serviceNumericToFloat64(s.PricePerKg),
			Active:      s.Active.Bool,
		}
	}
	return items, nil
}

// CreateService creates a new service in database.
func (r *ServiceRepository) CreateService(ctx context.Context, code, name string, duration int, price float64, active bool) (model.ServiceItem, error) {
	id := uuid.New()
	var idBytes [16]byte
	copy(idBytes[:], id[:])

	priceNumeric := pgtype.Numeric{}
	_ = priceNumeric.Scan(fmt.Sprintf("%.2f", price))

	created, err := r.queries.CreateService(ctx, db.CreateServiceParams{
		ID:          pgtype.UUID{Bytes: idBytes, Valid: true},
		ServiceCode: code,
		ServiceName: name,
		DurationDay: pgtype.Int4{Int32: int32(duration), Valid: true},
		PricePerKg:  priceNumeric,
		Active:      pgtype.Bool{Bool: active, Valid: true},
	})
	if err != nil {
		return model.ServiceItem{}, err
	}

	return model.ServiceItem{
		ID:          id.String(),
		ServiceCode: created.ServiceCode,
		ServiceName: created.ServiceName,
		DurationDay: int(created.DurationDay.Int32),
		PricePerKg:  price,
		Active:      created.Active.Bool,
	}, nil
}

// UpdateService updates an existing service.
func (r *ServiceRepository) UpdateService(ctx context.Context, id uuid.UUID, name string, duration int, price float64, active bool) (model.ServiceItem, error) {
	var idBytes [16]byte
	copy(idBytes[:], id[:])

	priceNumeric := pgtype.Numeric{}
	_ = priceNumeric.Scan(fmt.Sprintf("%.2f", price))

	updated, err := r.queries.UpdateService(ctx, db.UpdateServiceParams{
		ID:          pgtype.UUID{Bytes: idBytes, Valid: true},
		ServiceName: name,
		DurationDay: pgtype.Int4{Int32: int32(duration), Valid: true},
		PricePerKg:  priceNumeric,
		Active:      pgtype.Bool{Bool: active, Valid: true},
	})
	if err != nil {
		return model.ServiceItem{}, err
	}

	return model.ServiceItem{
		ID:          id.String(),
		ServiceCode: updated.ServiceCode,
		ServiceName: updated.ServiceName,
		DurationDay: int(updated.DurationDay.Int32),
		PricePerKg:  price,
		Active:      updated.Active.Bool,
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
