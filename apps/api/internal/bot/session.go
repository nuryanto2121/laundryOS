package bot

import "time"

// Conversation state constants for WhatsApp interactive multi-turn flows.
const (
	StateIdle                 = "IDLE"
	StateAwaitingAddress      = "AWAITING_ADDRESS"
	StateAwaitingService      = "AWAITING_SERVICE"
	StateAwaitingSchedule     = "AWAITING_SCHEDULE"
	StateAwaitingConfirmation = "AWAITING_CONFIRMATION"
)

// SessionData holds the in-flight conversation state and draft order details for a customer.
type SessionData struct {
	Phone        string    `json:"phone"`
	CustomerName string    `json:"customer_name"`
	State        string    `json:"state"`
	Address      string    `json:"address"`
	ServiceID    string    `json:"service_id"`
	ServiceName  string    `json:"service_name"`
	Schedule     string    `json:"schedule"`
	LastActive   time.Time `json:"last_active"`
}

// IsExpired checks if the session has exceeded the specified TTL.
func (s *SessionData) IsExpired(ttl time.Duration) bool {
	if s == nil {
		return true
	}
	return time.Since(s.LastActive) > ttl
}
