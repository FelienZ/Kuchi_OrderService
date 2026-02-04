package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type Entity int

const (
	ORDER Entity = iota
	PRODUCT
)

var entityToString = map[Entity]string{
	ORDER:   "ORDER",
	PRODUCT: "PRODUCT",
}
var stringToEntity = map[string]Entity{
	"ORDER":   ORDER,
	"PRODUCT": PRODUCT,
}

func (e Entity) String() string {
	if e, ok := entityToString[e]; ok {
		return e
	}
	return "UNKNOWN"
}

func (e Entity) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

func (e *Entity) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	if val, ok := stringToEntity[str]; ok {
		*e = val
		return nil
	}
	return fmt.Errorf("invalid entity: %s", str)
}

type TransactionLog struct {
	ID        string    `json:"id"`
	Entity    Entity    `json:"entity"`
	EntityID  string    `json:"entity_id"`
	Action    string    `json:"action"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
}
