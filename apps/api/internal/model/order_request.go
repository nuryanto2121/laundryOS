package model

// CreateOrderRequest holds the input fields for creating an order.
type CreateOrderRequest struct {
	CustomerID     string  `json:"customer_id" validate:"required,uuid"`
	ServiceID      string  `json:"service_id" validate:"required,uuid"`
	OutletID       string  `json:"outlet_id"`
	TotalWeight    float64 `json:"total_weight" validate:"gte=0"`
	TotalPrice     float64 `json:"total_price" validate:"gte=0"`
	PickupRequired bool    `json:"pickup_required"`
	Notes          string  `json:"notes"`
}

// UpdateOrderRequest holds the mutable input fields for updating an order.
// Order number and status remain immutable until their dedicated workflows exist.
type UpdateOrderRequest struct {
	CustomerID     string  `json:"customer_id" validate:"required,uuid"`
	ServiceID      string  `json:"service_id" validate:"required,uuid"`
	OutletID       string  `json:"outlet_id"`
	TotalWeight    float64 `json:"total_weight" validate:"gte=0"`
	TotalPrice     float64 `json:"total_price" validate:"gte=0"`
	PickupRequired bool    `json:"pickup_required"`
	Notes          string  `json:"notes"`
}
