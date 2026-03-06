package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type Role int

const (
	Member Role = iota
	Admin
)

var MapRoleToString = map[Role]string{
	Member: "Member",
	Admin:  "Admin",
}
var MapStringToRole = map[string]Role{
	"Member": Member,
	"Admin":  Admin,
}

func (r Role) String() string {
	if val, ok := MapRoleToString[r]; ok {
		return val
	}
	return "UNKNOWN"
}

func (r *Role) Scan(value any) error {
	switch v := value.(type) {
	case string:
		*r = MapStringToRole[v]
	case []byte:
		*r = MapStringToRole[string(v)]
	default:
		return fmt.Errorf("unsupported type %T for Role", value)
	}
	return nil
}

func (r Role) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func (r *Role) UnmarshalJSON(u []byte) error {
	var roleString string
	if err := json.Unmarshal(u, &roleString); err != nil {
		return err
	}
	if val, ok := MapStringToRole[roleString]; ok {
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

type UserIdentity struct {
	ID       string `json:"id"`
	Role     Role   `json:"role"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserUpdateRequest struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}
