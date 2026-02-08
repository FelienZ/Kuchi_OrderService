package services

import "errors"

var (
	ErrOrderNotFound = errors.New("Order Data Not Found")
	ErrOrderConflict = errors.New("Conflict Action Order")
	ErrOrderInvalid  = errors.New("Invalid Payload for Order Data")
)
