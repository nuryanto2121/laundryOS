package model

// CreateOrderRequest holds the input fields for creating an order.
type CreateOrderRequest struct {
	CustomerID     string  `json:"customer_id"`
	ServiceID      string  `json:"service_id"`
	TotalWeight    float64 `json:"total_weight"`
	TotalPrice     float64 `json:"total_price"`
	PickupRequired bool    `json:"pickup_required"`
	Notes          string  `json:"notes"`
}

// UpdateOrderRequest holds the mutable input fields for updating an order.
// Order number and status remain immutable until their dedicated workflows exist.
type UpdateOrderRequest struct {
	CustomerID     string  `json:"customer_id"`
	ServiceID      string  `json:"service_id"`
	TotalWeight    float64 `json:"total_weight"`
	TotalPrice     float64 `json:"total_price"`
	PickupRequired bool    `json:"pickup_required"`
	Notes          string  `json:"notes"`
}
