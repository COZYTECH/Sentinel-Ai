package events

import "time"

type TransactionCreatedEvent struct {
	Event         string    `json:"event"`
	TransactionID string    `json:"transaction_id"`
	UserID        string    `json:"user_id"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	Country       string    `json:"country"`
	DeviceID      string    `json:"device_id"`
	IPAddress     string    `json:"ip_address"`
	Timestamp     time.Time `json:"timestamp"`
}

type EmployeeActivityEvent struct {
	Event      string    `json:"event"`
	EmployeeID string    `json:"employee_id"`
	Action     string    `json:"action"`
	ResourceID string    `json:"resource_id"`
	IPAddress  string    `json:"ip_address"`
	Timestamp  time.Time `json:"timestamp"`
}