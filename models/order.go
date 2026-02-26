package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type OrderItem struct {
	ProductID string `json:"product_id"`
	Qty       int    `json:"qty"`
}

type Status int

const (
	PENDING Status = iota
	CANCELED
	PAID
)

var statusToString = map[Status]string{
	PAID:     "PAID",
	CANCELED: "CANCELED",
	PENDING:  "PENDING",
}
var stringToStatus = map[string]Status{
	"PAID":     PAID,
	"CANCELED": CANCELED,
	"PENDING":  PENDING,
}

func (s Status) String() string {
	if s, ok := statusToString[s]; ok {
		return s
	}
	return "UNKNOWN"
}

func (s Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String()) // buat encode ke json
}

func (s *Status) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	// balik ke status
	if val, ok := stringToStatus[str]; ok {
		*s = val
		return nil
	}
	return fmt.Errorf("invalid status: %s", str)
}

type Order struct {
	ID        string      `json:"id"`
	UserID    string      `json:"user_id"`
	Item      []OrderItem `json:"item"`
	Status    Status      `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type OrderRequest struct {
	Item []OrderItem `json:"item"`
}

type GetOrderParameter struct {
	Status Status
	Limit  int
	Offset int
}
