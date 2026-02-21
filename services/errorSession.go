package services

import "errors"

var (
	ErrSessionInvalid            = errors.New("Invalid Session Payload")
	ErrSessionInvalidCredentials = errors.New("Session Invalid")
)
