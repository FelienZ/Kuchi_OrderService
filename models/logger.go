package models

import "time"

type Entity int

const (
	ORDER Entity = iota
	PRODUCT
)

type TransactionLog struct {
	ID        string    `json:"id"`
	Entity    Entity    `json:"entity"`
	EntityID  string    `json:"entity_id"`
	Action    string    `json:"action"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
}
