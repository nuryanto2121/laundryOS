package model

// OrderStatus represents the current lifecycle state of an order.
type OrderStatus string

const (
	OrderStatusNew           OrderStatus = "NEW"
	OrderStatusWaitingPickup OrderStatus = "WAITING_PICKUP"
	OrderStatusConfirmed     OrderStatus = "CONFIRMED"
	OrderStatusCancelled     OrderStatus = "CANCELLED"
	OrderStatusCompleted     OrderStatus = "COMPLETED"
	OrderStatusInProgress    OrderStatus = "IN_PROGRESS"
)
