package services

import "errors"

var (
	ErrInvalid   = errors.New("Invalid Payload")
	ErrNotEnough = errors.New("Stock Not Enough")
	ErrConflict  = errors.New("Conflict Action")
)
