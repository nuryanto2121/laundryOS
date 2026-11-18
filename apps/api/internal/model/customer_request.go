package model

// CreateCustomerRequest holds the input fields for creating a customer.
type CreateCustomerRequest struct {
	FullName  string   `json:"full_name" validate:"required"`
	Phone     string   `json:"phone" validate:"required,min=8"`
	Email     string   `json:"email" validate:"omitempty,email"`
	Address   string   `json:"address"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

// UpdateCustomerRequest holds the input fields for updating a customer.
type UpdateCustomerRequest struct {
	FullName  string   `json:"full_name" validate:"required"`
	Phone     string   `json:"phone" validate:"required,min=8"`
	Email     string   `json:"email" validate:"omitempty,email"`
	Address   string   `json:"address"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

// UpdateCustomerStatusRequest holds the input field for changing customer status.
type UpdateCustomerStatusRequest struct {
	Status CustomerStatus `json:"status" validate:"required,oneof=ACTIVE INACTIVE BLOCKED"`
}
