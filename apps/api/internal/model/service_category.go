package model

// ServiceCategoryItem represents a category of laundry services/products.
type ServiceCategoryItem struct {
	ID           string `json:"id"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	DefaultUnit  string `json:"default_unit"`
	DisplayOrder int    `json:"display_order"`
	IsActive     bool   `json:"is_active"`
}
