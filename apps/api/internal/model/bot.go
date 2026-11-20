package model

// IncomingMessageRequest represents a webhook payload forwarded by WhatsApp Gateway.
type IncomingMessageRequest struct {
	From        string `json:"from" validate:"required"`
	Text        string `json:"text" validate:"required"`
	PushName    string `json:"push_name"`
	Timestamp   int64  `json:"timestamp"`
	IsGroup     bool   `json:"is_group"`
	GroupJID    string `json:"group_jid"`
	Participant string `json:"participant"`
	QuotedText  string `json:"quoted_text"`
}

// BotResponse represents the processed FAQ response.
type BotResponse struct {
	Category string `json:"category"`
	Reply    string `json:"reply"`
}

// ServiceItem represents a laundry service fetched from database.
type ServiceItem struct {
	ID           string  `json:"id"`
	ServiceCode  string  `json:"service_code"`
	ServiceName  string  `json:"service_name"`
	DurationDay  int     `json:"duration_day"`
	PricePerKg   float64 `json:"price_per_kg"`
	Price        float64 `json:"price"`
	Active       bool    `json:"active"`
	IsActive     bool    `json:"is_active"`
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Category     string  `json:"category"`
	Unit         string  `json:"unit"`
}
