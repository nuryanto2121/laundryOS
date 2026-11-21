package service

import (
	"context"
	"strings"

	"github.com/google/uuid"
	db "github.com/nuryanto2121/laundryos/apps/api/internal/db/generated"
	"github.com/nuryanto2121/laundryos/apps/api/internal/model"
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

// CreatePerfume creates a new perfume.
func (s *MasterService) CreatePerfume(ctx context.Context, name, description string, isActive bool) (db.Perfume, error) {
	return s.repo.CreatePerfume(ctx, name, description, isActive)
}

// UpdatePerfume updates an existing perfume.
func (s *MasterService) UpdatePerfume(ctx context.Context, idStr string, name, description string, isActive bool) (db.Perfume, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return db.Perfume{}, err
	}
	return s.repo.UpdatePerfume(ctx, id, name, description, isActive)
}

// DeletePerfume deactivates a perfume.
func (s *MasterService) DeletePerfume(ctx context.Context, idStr string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return err
	}
	return s.repo.DeletePerfume(ctx, id)
}

// ListDiscounts returns all active discounts.
func (s *MasterService) ListDiscounts(ctx context.Context) ([]db.Discount, error) {
	return s.repo.ListDiscounts(ctx)
}

// CreateDiscount creates a new discount.
func (s *MasterService) CreateDiscount(ctx context.Context, name, discountType string, value, minOrder float64, isActive bool) (db.Discount, error) {
	return s.repo.CreateDiscount(ctx, name, discountType, value, minOrder, isActive)
}

// UpdateDiscount updates an existing discount.
func (s *MasterService) UpdateDiscount(ctx context.Context, idStr string, name, discountType string, value, minOrder float64, isActive bool) (db.Discount, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return db.Discount{}, err
	}
	return s.repo.UpdateDiscount(ctx, id, name, discountType, value, minOrder, isActive)
}

// DeleteDiscount deactivates a discount.
func (s *MasterService) DeleteDiscount(ctx context.Context, idStr string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return err
	}
	return s.repo.DeleteDiscount(ctx, id)
}

// ListPaymentMethods returns all active payment methods.
func (s *MasterService) ListPaymentMethods(ctx context.Context) ([]db.PaymentMethod, error) {
	return s.repo.ListPaymentMethods(ctx)
}

// ListOutlets returns all active outlets with GPS coordinates.
func (s *MasterService) ListOutlets(ctx context.Context) ([]model.Outlet, error) {
	return s.repo.ListOutletsWithGPS(ctx)
}

// CreateOutlet creates a new outlet with GPS coordinates.
func (s *MasterService) CreateOutlet(ctx context.Context, req model.CreateOutletRequest) (model.Outlet, error) {
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		req.Name = "Outlet Laundry"
	}
	req.Code = strings.TrimSpace(req.Code)
	if req.Code == "" {
		req.Code = "OUT-" + strings.ToUpper(uuid.New().String()[:6])
	}
	return s.repo.CreateOutlet(ctx, req)
}

// UpdateOutlet updates outlet profile and GPS coordinates.
func (s *MasterService) UpdateOutlet(ctx context.Context, idStr string, req model.UpdateOutletRequest) (model.Outlet, error) {
	return s.repo.UpdateOutlet(ctx, idStr, req)
}

// DeleteOutlet deactivates an outlet.
func (s *MasterService) DeleteOutlet(ctx context.Context, idStr string) error {
	return s.repo.DeleteOutlet(ctx, idStr)
}

// GetDefaultOutlet returns the primary outlet.
func (s *MasterService) GetDefaultOutlet(ctx context.Context) (db.Outlet, error) {
	return s.repo.GetDefaultOutlet(ctx)
}

// ListExpenseCategories returns all active expense categories.
func (s *MasterService) ListExpenseCategories(ctx context.Context) ([]db.ExpenseCategory, error) {
	return s.repo.ListExpenseCategories(ctx)
}

// ListAllExpenseCategories returns all expense categories with optional sorting.
func (s *MasterService) ListAllExpenseCategories(ctx context.Context, sortBy, sortDir string) ([]db.ExpenseCategory, error) {
	return s.repo.ListAllExpenseCategories(ctx, sortBy, sortDir)
}
func (s *MasterService) CreateExpenseCategory(ctx context.Context, code, name string, isActive bool) (db.ExpenseCategory, error) {
	return s.repo.CreateExpenseCategory(ctx, code, name, isActive)
}

// UpdateExpenseCategory updates an existing expense category.
func (s *MasterService) UpdateExpenseCategory(ctx context.Context, idStr, name string, isActive bool) (db.ExpenseCategory, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return db.ExpenseCategory{}, err
	}
	return s.repo.UpdateExpenseCategory(ctx, id, name, isActive)
}

// DeleteExpenseCategory deactivates an expense category.
func (s *MasterService) DeleteExpenseCategory(ctx context.Context, idStr string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return err
	}
	return s.repo.DeleteExpenseCategory(ctx, id)
}

// ListEmployees returns all active employees.
func (s *MasterService) ListEmployees(ctx context.Context) ([]db.Employee, error) {
	return s.repo.ListEmployees(ctx)
}

// CreateEmployee creates a new employee.
func (s *MasterService) CreateEmployee(ctx context.Context, username, fullName, phone, role, pin string, isActive bool) (db.Employee, error) {
	return s.repo.CreateEmployee(ctx, username, fullName, phone, role, pin, isActive)
}

// UpdateEmployee updates an employee.
func (s *MasterService) UpdateEmployee(ctx context.Context, idStr, fullName, phone, role, pin string, isActive bool) (db.Employee, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return db.Employee{}, err
	}
	return s.repo.UpdateEmployee(ctx, id, fullName, phone, role, pin, isActive)
}

// DeleteEmployee deactivates an employee.
func (s *MasterService) DeleteEmployee(ctx context.Context, idStr string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return err
	}
	return s.repo.DeleteEmployee(ctx, id)
}

// ListServiceCategories returns all active service categories.
func (s *MasterService) ListServiceCategories(ctx context.Context) ([]db.ServiceCategory, error) {
	return s.repo.ListServiceCategories(ctx)
}

// ListAllServiceCategories returns all service categories.
func (s *MasterService) ListAllServiceCategories(ctx context.Context) ([]db.ServiceCategory, error) {
	return s.repo.ListAllServiceCategories(ctx)
}

// CreateServiceCategory creates a new service category.
func (s *MasterService) CreateServiceCategory(ctx context.Context, code, name, defaultUnit string, displayOrder int, isActive bool) (db.ServiceCategory, error) {
	return s.repo.CreateServiceCategory(ctx, code, name, defaultUnit, displayOrder, isActive)
}

// UpdateServiceCategory updates an existing service category.
func (s *MasterService) UpdateServiceCategory(ctx context.Context, idStr, name, defaultUnit string, displayOrder int, isActive bool) (db.ServiceCategory, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return db.ServiceCategory{}, err
	}
	return s.repo.UpdateServiceCategory(ctx, id, name, defaultUnit, displayOrder, isActive)
}

// DeleteServiceCategory deactivates a service category.
func (s *MasterService) DeleteServiceCategory(ctx context.Context, idStr string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return err
	}
	return s.repo.DeleteServiceCategory(ctx, id)
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
	Active      bool    `json:"active"`
}

// ListServices returns services with category and unit mapping, falling back to FAN Laundry defaults.
func (s *MasterService) ListServices(ctx context.Context, sortBy, sortDir string) ([]ServiceMasterItem, error) {
	if s.serviceRepo != nil {
		active, err := s.serviceRepo.ListAllServices(ctx, sortBy, sortDir)
		if err == nil && len(active) > 0 {
			res := make([]ServiceMasterItem, len(active))
			for i, a := range active {
				cat := a.CategoryName
				unit := a.Unit
				if cat == "" || unit == "" {
					detCat, detUnit := determineCategoryAndUnit(a.ServiceName)
					if cat == "" {
						cat = detCat
					}
					if unit == "" {
						unit = detUnit
					}
				}
				res[i] = ServiceMasterItem{
					ID:          a.ID,
					ServiceCode: a.ServiceCode,
					ServiceName: a.ServiceName,
					Category:    cat,
					Unit:        unit,
					Price:       a.PricePerKg,
					DurationDay: a.DurationDay,
					Active:      a.Active,
				}
			}
			return res, nil
		}
	}

	// Default FAN Laundry service catalog
	return GetDefaultFANServices(), nil
}

// CreateService creates a new service in catalog.
func (s *MasterService) CreateService(ctx context.Context, code, name string, duration int, price float64, active bool, categoryID, categoryName, unit string) (model.ServiceItem, error) {
	if s.serviceRepo != nil {
		return s.serviceRepo.CreateService(ctx, code, name, duration, price, active, categoryID, categoryName, unit)
	}
	return model.ServiceItem{}, nil
}

// UpdateService updates an existing service in catalog.
func (s *MasterService) UpdateService(ctx context.Context, idStr, name string, duration int, price float64, active bool, categoryID, categoryName, unit string) (model.ServiceItem, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return model.ServiceItem{}, err
	}
	if s.serviceRepo != nil {
		return s.serviceRepo.UpdateService(ctx, id, name, duration, price, active, categoryID, categoryName, unit)
	}
	return model.ServiceItem{}, nil
}

// DeleteService deactivates a service.
func (s *MasterService) DeleteService(ctx context.Context, idStr string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return err
	}
	if s.serviceRepo != nil {
		return s.serviceRepo.DeleteService(ctx, id)
	}
	return nil
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

// GetDefaultFANServices returns standard FAN Laundry catalog.
func GetDefaultFANServices() []ServiceMasterItem {
	return []ServiceMasterItem{
		{ID: "srv-01", ServiceCode: "SS-CB", ServiceName: "Cuci Basah (Self Service)", Category: "TABUNG", Unit: "Tabung", Price: 8000, DurationDay: 0, Active: true},
		{ID: "srv-02", ServiceCode: "SS-CK", ServiceName: "Cuci Kering (Self Service)", Category: "TABUNG", Unit: "Tabung", Price: 20000, DurationDay: 0, Active: true},
		{ID: "srv-03", ServiceCode: "SS-SBP", ServiceName: "Sabun / Pewangi Sachet", Category: "SATUAN", Unit: "Pcs", Price: 1000, DurationDay: 0, Active: true},
		{ID: "srv-04", ServiceCode: "CS-REG", ServiceName: "Cuci Setrika Reguler (3 Hari)", Category: "KILOAN", Unit: "Kg", Price: 7000, DurationDay: 3, Active: true},
		{ID: "srv-05", ServiceCode: "CS-EXP", ServiceName: "Cuci Setrika Express (2 Hari)", Category: "KILOAN", Unit: "Kg", Price: 9000, DurationDay: 2, Active: true},
		{ID: "srv-06", ServiceCode: "CS-KLT", ServiceName: "Cuci Setrika Kilat (1 Hari / 24 Jam)", Category: "KILOAN", Unit: "Kg", Price: 11000, DurationDay: 1, Active: true},
		{ID: "srv-07", ServiceCode: "CS-SKL", ServiceName: "Cuci Setrika Super Kilat (12 Jam)", Category: "KILOAN", Unit: "Kg", Price: 13000, DurationDay: 1, Active: true},
		{ID: "srv-08", ServiceCode: "CS-TRB", ServiceName: "Cuci Setrika Turbo (6 Jam)", Category: "KILOAN", Unit: "Kg", Price: 16000, DurationDay: 0, Active: true},
		{ID: "srv-09", ServiceCode: "SS-REG", ServiceName: "Setrika Saja Reguler (3 Hari)", Category: "KILOAN", Unit: "Kg", Price: 6000, DurationDay: 3, Active: true},
		{ID: "srv-10", ServiceCode: "SS-EXP", ServiceName: "Setrika Saja Express (2 Hari)", Category: "KILOAN", Unit: "Kg", Price: 8000, DurationDay: 2, Active: true},
		{ID: "srv-11", ServiceCode: "SS-KLT", ServiceName: "Setrika Saja Kilat (1 Hari / 24 Jam)", Category: "KILOAN", Unit: "Kg", Price: 10000, DurationDay: 1, Active: true},
		{ID: "srv-12", ServiceCode: "BC-REG", ServiceName: "Bad Cover Reguler (3 Hari)", Category: "SATUAN", Unit: "Pcs", Price: 30000, DurationDay: 3, Active: true},
		{ID: "srv-13", ServiceCode: "BC-EXP", ServiceName: "Bad Cover Express (2 Hari)", Category: "SATUAN", Unit: "Pcs", Price: 35000, DurationDay: 2, Active: true},
		{ID: "srv-14", ServiceCode: "KP-MED", ServiceName: "Karpet Permadani Sedang", Category: "METERAN", Unit: "Meter", Price: 20000, DurationDay: 7, Active: true},
		{ID: "srv-15", ServiceCode: "KP-THK", ServiceName: "Karpet Permadani Tebal", Category: "METERAN", Unit: "Meter", Price: 30000, DurationDay: 7, Active: true},
		{ID: "srv-16", ServiceCode: "BN-SML", ServiceName: "Boneka Sedang", Category: "SATUAN", Unit: "Pcs", Price: 35000, DurationDay: 3, Active: true},
		{ID: "srv-17", ServiceCode: "BN-BIG", ServiceName: "Boneka Besar", Category: "SATUAN", Unit: "Pcs", Price: 60000, DurationDay: 3, Active: true},
	}
}
