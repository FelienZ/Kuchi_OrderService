package services

import (
	"errors"
	"go-inventory/repository/database"
)

var (
	ErrUserInvalid  = errors.New("Invalid User Data")
	ErrUserConflict = errors.New("Conflict User Data")
	ErrUserNotFound = errors.New("Not Found User Data")
)

func ErrorUserDomainTranslator(err error) error {
	switch err {
	case database.ErrUniqueViolated:
		return ErrUserConflict
	case database.ErrCheckViolated:
		return ErrUserInvalid
	case database.ErrNotNullViolated:
		return ErrUserInvalid
	case database.ErrForeignViolated:
		return ErrUserConflict
	default:
		return nil
	}
}
