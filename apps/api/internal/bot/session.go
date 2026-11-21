package bot

import "time"

// Conversation state constants for WhatsApp interactive multi-turn flows.
const (
	StateIdle                    = "IDLE"
	StateAwaitingBranchSelection = "AWAITING_BRANCH_SELECTION"
	StateAwaitingAddress         = "AWAITING_ADDRESS"
	StateAwaitingService         = "AWAITING_SERVICE"
	StateAwaitingDetergent       = "AWAITING_DETERGENT"
	StateAwaitingSchedule        = "AWAITING_SCHEDULE"
	StateAwaitingConfirmation    = "AWAITING_CONFIRMATION"
	StateAwaitingOverloadAction  = "AWAITING_OVERLOAD_ACTION"
)

// SessionData holds the in-flight conversation state and draft order details for a customer.
type SessionData struct {
	Phone          string    `json:"phone"`
	CustomerName   string    `json:"customer_name"`
	State          string    `json:"state"`
	OutletID       string    `json:"outlet_id"`
	OutletName     string    `json:"outlet_name"`
	TargetGroupJID string    `json:"target_group_jid"`
	Address        string    `json:"address"`
	ServiceID      string    `json:"service_id"`
	ServiceName    string    `json:"service_name"`
	Schedule       string    `json:"schedule"`
	DetergentAddon string    `json:"detergent_addon"`
	DetergentPrice float64   `json:"detergent_price"`
	OrderNumber    string    `json:"order_number"`
	OverloadWeight float64   `json:"overload_weight"`
	OverloadExcess float64   `json:"overload_excess"`
	LastActive     time.Time `json:"last_active"`
}

// IsExpired checks if the session has exceeded the specified TTL.
func (s *SessionData) IsExpired(ttl time.Duration) bool {
	if s == nil {
		return true
	}
	return time.Since(s.LastActive) > ttl
}
