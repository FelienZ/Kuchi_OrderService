package models

import "time"

type UserSession struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	ExpiresAt time.Time `json:"expired_at"`
}

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResult struct {
	Identity  UserIdentity `json:"identity"`
	SessionID string       `json:"session_id"`
}
