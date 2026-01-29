package repository

import "errors"

var (
	ErrNotFound = errors.New("Item Not Found")
	ErrConflict = errors.New("Item Already Exist")
)
