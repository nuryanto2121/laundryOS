package service

import (
	"context"
	"strings"

	db "github.com/nuryanto2121/laundryos/apps/api/internal/db/generated"
	"github.com/nuryanto2121/laundryos/apps/api/internal/repository"
)

// MasterService handles business logic for master lookup data.
type MasterService struct {
	repo        *repository.MasterRepository
	serviceRepo *repository.ServiceRepository
}

// NewMasterService creates a new MasterService.
func NewMasterService(repo *repository.MasterRepository, serviceRepo *repository.ServiceRepository) *MasterService {
	return &MasterService{
		repo:        repo,
		serviceRepo: serviceRepo,
	}
}

// ListPerfumes returns all active perfumes.
func (s *MasterService) ListPerfumes(ctx context.Context) ([]db.Perfume, error) {
	return s.repo.ListPerfumes(ctx)
}

// ListDiscounts returns all active discounts.
func (s *MasterService) ListDiscounts(ctx context.Context) ([]db.Discount, error) {
	return s.repo.ListDiscounts(ctx)
}

// ListPaymentMethods returns all active payment methods.
func (s *MasterService) ListPaymentMethods(ctx context.Context) ([]db.PaymentMethod, error) {
	return s.repo.ListPaymentMethods(ctx)
}

// ListOutlets returns all active outlets.
func (s *MasterService) ListOutlets(ctx context.Context) ([]db.Outlet, error) {
	return s.repo.ListOutlets(ctx)
}

// GetDefaultOutlet returns the primary outlet.
func (s *MasterService) GetDefaultOutlet(ctx context.Context) (db.Outlet, error) {
	return s.repo.GetDefaultOutlet(ctx)
}

// ServiceMasterItem represents a service option in POS.
type ServiceMasterItem struct {
	ID          string  `json:"id"`
	ServiceCode string  `json:"service_code"`
	ServiceName string  `json:"service_name"`
	Category    string  `json:"category"`
	Unit        string  `json:"unit"`
	Price       float64 `json:"price"`
	DurationDay int     `json:"duration_day"`
}

// ListServices returns services with category and unit mapping, falling back to FAN Laundry defaults.
func (s *MasterService) ListServices(ctx context.Context) ([]ServiceMasterItem, error) {
	if s.serviceRepo != nil {
		active, err := s.serviceRepo.ListActiveServices(ctx)
		if err == nil && len(active) > 0 {
			res := make([]ServiceMasterItem, len(active))
			for i, a := range active {
				cat, unit := determineCategoryAndUnit(a.ServiceName)
				res[i] = ServiceMasterItem{
					ID:          a.ID,
					ServiceCode: a.ServiceCode,
					ServiceName: a.ServiceName,
					Category:    cat,
					Unit:        unit,
					Price:       a.PricePerKg,
					DurationDay: a.DurationDay,
				}
			}
			return res, nil
		}
	}

	// Default FAN Laundry service catalog
	return GetDefaultFANServices(), nil
}

func determineCategoryAndUnit(name string) (string, string) {
	lower := strings.ToLower(name)
	if strings.Contains(lower, "self service") || strings.Contains(lower, "cuci basah") || strings.Contains(lower, "cuci kering") {
		return "TABUNG", "Tabung"
	}
	if strings.Contains(lower, "sabun") || strings.Contains(lower, "pewangi") {
		return "SATUAN", "Pcs"
	}
	if strings.Contains(lower, "karpet") || strings.Contains(lower, "permadani") {
		return "METERAN", "Meter"
	}
	if strings.Contains(lower, "bed cover") || strings.Contains(lower, "bad cover") || strings.Contains(lower, "boneka") || strings.Contains(lower, "selimut") {
		return "SATUAN", "Pcs"
	}
	return "KILOAN", "Kg"
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) > 0 && indexOf(s, substr) >= 0)
}

func indexOf(s, substr string) int {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// GetDefaultFANServices returns standard FAN Laundry catalog.
func GetDefaultFANServices() []ServiceMasterItem {
	return []ServiceMasterItem{
		{ID: "srv-01", ServiceCode: "SS-CB", ServiceName: "Cuci Basah (Self Service)", Category: "TABUNG", Unit: "Tabung", Price: 8000, DurationDay: 0},
		{ID: "srv-02", ServiceCode: "SS-CK", ServiceName: "Cuci Kering (Self Service)", Category: "TABUNG", Unit: "Tabung", Price: 20000, DurationDay: 0},
		{ID: "srv-03", ServiceCode: "SS-SBP", ServiceName: "Sabun / Pewangi Sachet", Category: "SATUAN", Unit: "Pcs", Price: 1000, DurationDay: 0},
		{ID: "srv-04", ServiceCode: "CS-REG", ServiceName: "Cuci Setrika Reguler (3 Hari)", Category: "KILOAN", Unit: "Kg", Price: 7000, DurationDay: 3},
		{ID: "srv-05", ServiceCode: "CS-EXP", ServiceName: "Cuci Setrika Express (2 Hari)", Category: "KILOAN", Unit: "Kg", Price: 9000, DurationDay: 2},
		{ID: "srv-06", ServiceCode: "CS-KLT", ServiceName: "Cuci Setrika Kilat (1 Hari / 24 Jam)", Category: "KILOAN", Unit: "Kg", Price: 11000, DurationDay: 1},
		{ID: "srv-07", ServiceCode: "CS-SKL", ServiceName: "Cuci Setrika Super Kilat (12 Jam)", Category: "KILOAN", Unit: "Kg", Price: 13000, DurationDay: 1},
		{ID: "srv-08", ServiceCode: "CS-TRB", ServiceName: "Cuci Setrika Turbo (6 Jam)", Category: "KILOAN", Unit: "Kg", Price: 16000, DurationDay: 0},
		{ID: "srv-09", ServiceCode: "SS-REG", ServiceName: "Setrika Saja Reguler (3 Hari)", Category: "KILOAN", Unit: "Kg", Price: 6000, DurationDay: 3},
		{ID: "srv-10", ServiceCode: "SS-EXP", ServiceName: "Setrika Saja Express (2 Hari)", Category: "KILOAN", Unit: "Kg", Price: 8000, DurationDay: 2},
		{ID: "srv-11", ServiceCode: "SS-KLT", ServiceName: "Setrika Saja Kilat (1 Hari / 24 Jam)", Category: "KILOAN", Unit: "Kg", Price: 10000, DurationDay: 1},
		{ID: "srv-12", ServiceCode: "BC-REG", ServiceName: "Bad Cover Reguler (3 Hari)", Category: "SATUAN", Unit: "Pcs", Price: 30000, DurationDay: 3},
		{ID: "srv-13", ServiceCode: "BC-EXP", ServiceName: "Bad Cover Express (2 Hari)", Category: "SATUAN", Unit: "Pcs", Price: 35000, DurationDay: 2},
		{ID: "srv-14", ServiceCode: "KP-MED", ServiceName: "Karpet Permadani Sedang", Category: "METERAN", Unit: "Meter", Price: 20000, DurationDay: 7},
		{ID: "srv-15", ServiceCode: "KP-THK", ServiceName: "Karpet Permadani Tebal", Category: "METERAN", Unit: "Meter", Price: 30000, DurationDay: 7},
		{ID: "srv-16", ServiceCode: "BN-SML", ServiceName: "Boneka Sedang", Category: "SATUAN", Unit: "Pcs", Price: 35000, DurationDay: 3},
		{ID: "srv-17", ServiceCode: "BN-BIG", ServiceName: "Boneka Besar", Category: "SATUAN", Unit: "Pcs", Price: 60000, DurationDay: 3},
	}
}
