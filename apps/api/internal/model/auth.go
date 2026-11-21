package model

import "time"

// LoginRequest payload for authenticating users
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthUser represents user information returned after login
type AuthUser struct {
	ID          string          `json:"id"`
	Username    string          `json:"username"`
	FullName    string          `json:"full_name"`
	Phone       string          `json:"phone,omitempty"`
	Email       string          `json:"email,omitempty"`
	Role        string          `json:"role"`
	OutletID    string          `json:"outlet_id,omitempty"`
	OutletName  string          `json:"outlet_name,omitempty"`
	Permissions map[string]bool `json:"permissions"`
}

// LoginResponse contains authentication token and profile
type LoginResponse struct {
	Token string   `json:"token"`
	User  AuthUser `json:"user"`
}

// Employee represents detailed employee model with ACL and outlet assignment
type Employee struct {
	ID          string          `json:"id"`
	Username    string          `json:"username"`
	FullName    string          `json:"full_name"`
	Phone       string          `json:"phone,omitempty"`
	Email       string          `json:"email,omitempty"`
	Role        string          `json:"role"`
	Password    string          `json:"-"`
	OutletID    string          `json:"outlet_id,omitempty"`
	OutletName  string          `json:"outlet_name,omitempty"`
	Permissions map[string]bool `json:"permissions"`
	IsActive    bool            `json:"is_active"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// CreateEmployeeRequest payload to create a new employee
type CreateEmployeeRequest struct {
	Username    string          `json:"username"`
	FullName    string          `json:"full_name"`
	Phone       string          `json:"phone,omitempty"`
	Email       string          `json:"email,omitempty"`
	Role        string          `json:"role"`
	Password    string          `json:"password"`
	OutletID    string          `json:"outlet_id,omitempty"`
	Permissions map[string]bool `json:"permissions,omitempty"`
	IsActive    bool            `json:"is_active"`
}

// UpdateEmployeeRequest payload to update an employee
type UpdateEmployeeRequest struct {
	FullName    string          `json:"full_name"`
	Phone       string          `json:"phone,omitempty"`
	Email       string          `json:"email,omitempty"`
	Role        string          `json:"role"`
	Password    string          `json:"password,omitempty"`
	OutletID    string          `json:"outlet_id,omitempty"`
	Permissions map[string]bool `json:"permissions,omitempty"`
	IsActive    bool            `json:"is_active"`
}
