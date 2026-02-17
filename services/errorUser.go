package services

import "errors"

var (
	ErrUserInvalid  = errors.New("Invalid User Data Payload")
	ErrUserConflict = errors.New("Conflict User Data")
	ErrUserNotFound = errors.New("Not Found User Data")
)
