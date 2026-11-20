package model

import "time"

// POSOrderItemRequest represents an item in a POS order cart.
type POSOrderItemRequest struct {
	ServiceID    string  `json:"service_id"`
	ServiceName  string  `json:"service_name"`
	Category     string  `json:"category"`
	Quantity     float64 `json:"quantity"`
	Unit         string  `json:"unit"`
	PricePerUnit float64 `json:"price_per_unit"`
	Subtotal     float64 `json:"subtotal"`
	Notes        string  `json:"notes"`
}

// CreatePOSOrderRequest represents payload to create an order from POS.
type CreatePOSOrderRequest struct {
	// If existing customer
	CustomerID string `json:"customer_id"`
	// If new customer created inline
	CustomerName    string `json:"customer_name"`
	CustomerPhone   string `json:"customer_phone"`
	CustomerAddress string `json:"customer_address"`

	OutletID              string                `json:"outlet_id"`
	Items                 []POSOrderItemRequest `json:"items"`
	PerfumeID             string                `json:"perfume_id"`
	PerfumeName           string                `json:"perfume_name"`
	DiscountID            string                `json:"discount_id"`
	DiscountAmount        float64               `json:"discount_amount"`
	TotalWeight           float64               `json:"total_weight"`
	TotalPrice            float64               `json:"total_price"`
	PaidAmount            float64               `json:"paid_amount"`
	PaymentMethodID       string                `json:"payment_method_id"`
	PaymentMethodName     string                `json:"payment_method_name"`
	RackNumber            string                `json:"rack_number"`
	EstimatedCompletionAt *time.Time            `json:"estimated_completion_at"`
	Notes                 string                `json:"notes"`
	PickupRequired        bool                  `json:"pickup_required"`
	CreatedBy             string                `json:"created_by"`
}

// UpdatePaymentRequest represents payload to settle payment on an order.
type UpdatePaymentRequest struct {
	PaidAmount        float64 `json:"paid_amount"`
	PaymentMethodID   string  `json:"payment_method_id"`
	PaymentMethodName string  `json:"payment_method_name"`
	StaffName         string  `json:"staff_name"`
}

// UpdateStatusRequest represents payload to transition order status.
type UpdateStatusRequest struct {
	Status    string `json:"status"`
	ChangedBy string `json:"changed_by"`
	Notes     string `json:"notes"`
}

// CancelOrderRequest represents payload to cancel an order.
type CancelOrderRequest struct {
	Reason    string `json:"reason"`
	ChangedBy string `json:"changed_by"`
}

// UpdateRackRequest represents payload to assign or move clothes to a rack.
type UpdateRackRequest struct {
	RackNumber string `json:"rack_number"`
}

// AddOrderPhotoRequest represents payload to attach photo documentation.
type AddOrderPhotoRequest struct {
	PhotoType string `json:"photo_type"` // 'MASUK' | 'KELUAR'
	PhotoURL  string `json:"photo_url"`
	Notes     string `json:"notes"`
}

// CreateExpenseRequest represents payload to record an expense.
type CreateExpenseRequest struct {
	OutletID    string  `json:"outlet_id"`
	ExpenseDate string  `json:"expense_date"` // YYYY-MM-DD
	Category    string  `json:"category"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	CreatedBy   string  `json:"created_by"`
}

// UpdateExpenseRequest represents payload to edit/correct an expense.
type UpdateExpenseRequest struct {
	ExpenseDate string  `json:"expense_date"` // YYYY-MM-DD
	Category    string  `json:"category"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}
