package model

import "time"

// Outlet represents a branch/outlet in LaundryOS multi-outlet system
type Outlet struct {
	ID                   string    `json:"id"`
	Name                 string    `json:"name"`
	Code                 string    `json:"code"`
	Address              string    `json:"address"`
	Phone                string    `json:"phone"`
	Latitude             float64   `json:"latitude"`
	Longitude            float64   `json:"longitude"`
	GeofenceRadiusMeters int       `json:"geofence_radius_meters"`
	IsActive             bool      `json:"is_active"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// CreateOutletRequest represents payload to register a new outlet/cabang
type CreateOutletRequest struct {
	Name                 string  `json:"name"`
	Code                 string  `json:"code"`
	Address              string  `json:"address"`
	Phone                string  `json:"phone"`
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	GeofenceRadiusMeters int     `json:"geofence_radius_meters"`
	IsActive             bool    `json:"is_active"`
}

// UpdateOutletRequest represents payload to edit outlet profile & GPS coordinates
type UpdateOutletRequest struct {
	Name                 string  `json:"name"`
	Code                 string  `json:"code"`
	Address              string  `json:"address"`
	Phone                string  `json:"phone"`
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	GeofenceRadiusMeters int     `json:"geofence_radius_meters"`
	IsActive             bool    `json:"is_active"`
}
