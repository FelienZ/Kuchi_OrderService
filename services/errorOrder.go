package services

import (
	"errors"
	"go-inventory/repository/database"
)

var (
	ErrOrderNotFound = errors.New("Order Data Not Found")
	ErrOrderConflict = errors.New("Conflict Order Action")
	ErrOrderInvalid  = errors.New("Invalid Order Data")
)

func ErrorOrderDomainTranslator(err error) error {
	switch err {
	case database.ErrUniqueViolated:
		return ErrOrderConflict
	case database.ErrCheckViolated:
		return ErrOrderInvalid
	case database.ErrNotNullViolated:
		return ErrOrderInvalid
	case database.ErrForeignViolated:
		return ErrOrderConflict
	default:
		return nil
	}
}
