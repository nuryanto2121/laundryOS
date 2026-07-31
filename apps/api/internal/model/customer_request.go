package model

// CreateCustomerRequest holds the input fields for creating a customer.
type CreateCustomerRequest struct {
	CustomerCode string   `json:"customer_code"`
	FullName     string   `json:"full_name"`
	Phone        string   `json:"phone"`
	Email        string   `json:"email"`
	Address      string   `json:"address"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
}

// UpdateCustomerRequest holds the input fields for updating a customer.
type UpdateCustomerRequest struct {
	FullName  string   `json:"full_name"`
	Phone     string   `json:"phone"`
	Email     string   `json:"email"`
	Address   string   `json:"address"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}
