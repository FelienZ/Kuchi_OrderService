package services

import "errors"

var (
	ErrSessionInvalid            = errors.New("Invalid Session")
	ErrSessionInvalidCredentials = errors.New("Session Credentials Invalid")
)
