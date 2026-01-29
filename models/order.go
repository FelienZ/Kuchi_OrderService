package models

import "time"

type OrderItem struct {
	ProductID string `json:"product_id"`
	Qty       int    `json:"qty"`
}

type Status int

const (
	PENDING Status = iota
	CANCELLED
	PAID
)

type Order struct {
	ID        string      `json:"id"`
	UserID    string      `json:"user_id"`
	Item      []OrderItem `json:"item"`
	Status    Status      `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}
