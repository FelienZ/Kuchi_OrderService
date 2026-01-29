package repository

import "errors"

var (
	ErrNotFound = errors.New("Product Not Found")
	ErrConflict = errors.New("Product ID Already Exist")
)
