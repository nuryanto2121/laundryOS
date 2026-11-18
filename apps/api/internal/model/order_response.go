package model

// Order is the application-level representation of an order.
type Order struct {
	ID             string      `json:"id"`
	OrderNumber    string      `json:"order_number"`
	CustomerID     string      `json:"customer_id"`
	ServiceID      string      `json:"service_id"`
	Status         OrderStatus `json:"status"`
	TotalWeight    float64     `json:"total_weight"`
	TotalPrice     float64     `json:"total_price"`
	PickupRequired bool        `json:"pickup_required"`
	Notes          string      `json:"notes"`
	CreatedAt      string      `json:"created_at"`
	UpdatedAt      string      `json:"updated_at"`
}
