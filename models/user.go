package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type Role int

const (
	Guest Role = iota
	Member
	Admin
)

var mapRoleToString = map[Role]string{
	Guest:  "Guest",
	Member: "Member",
	Admin:  "Admin",
}
var mapStringToRole = map[string]Role{
	"Guest":  Guest,
	"Member": Member,
	"Admin":  Admin,
}

func (r Role) String() string {
	if val, ok := mapRoleToString[r]; ok {
		return val
	}
	return "UNKNOWN"
}

func (r Role) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func (r *Role) UnmarshalJSON(u []byte) error {
	var roleString string
	if err := json.Unmarshal(u, &roleString); err != nil {
		return err
	}
	if val, ok := mapStringToRole[roleString]; ok {
		*r = val
		return nil
	}
	return fmt.Errorf("invalid role: %s", roleString)
}

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserUpdateRequest struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}
