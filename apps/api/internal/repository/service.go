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
