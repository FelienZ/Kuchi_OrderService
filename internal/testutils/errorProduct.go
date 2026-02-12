package testutils

import "errors"

var (
	ErrProductNotFound  = errors.New("Product Data Not Found")
	ErrProductConflict  = errors.New("Conflict Action Product")
	ErrProductInvalid   = errors.New("Invalid Payload for Product Data")
	ErrProductNotEnough = errors.New("Not Enough Product Stock")
)
