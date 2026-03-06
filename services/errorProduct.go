package services

import (
	"errors"
	"go-inventory/repository/database"
)

var (
	ErrProductNotFound  = errors.New("Product Data Not Found")
	ErrProductConflict  = errors.New("Conflict Product Action")
	ErrProductInvalid   = errors.New("Invalid Product Data")
	ErrProductNotEnough = errors.New("Not Enough Product Stock")
)

func ErrorProductDomainTranslator(err error) error {
	switch err {
	case database.ErrUniqueViolated:
		return ErrProductConflict
	case database.ErrCheckViolated:
		return ErrProductInvalid
	case database.ErrNotNullViolated:
		return ErrProductInvalid
	case database.ErrForeignViolated:
		return ErrProductConflict
	default:
		return nil
	}
}
