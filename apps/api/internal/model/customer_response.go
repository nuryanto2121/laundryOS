package model

// Customer is the application-level representation of a customer.
// It uses only primitive Go types — no database-specific types.
type Customer struct {
	ID           string   `json:"id"`
	CustomerCode string   `json:"customer_code"`
	FullName     string   `json:"full_name"`
	Phone        string   `json:"phone"`
	Email        string   `json:"email"`
	Address      string   `json:"address"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
}
