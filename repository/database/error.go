package database

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrNoRows          = errors.New("None of Data Found")
	ErrNoUpdate        = errors.New("None of Rows Affected after Update")
	ErrNoDelete        = errors.New("None of Rows Affected after Delete ")
	ErrUniqueViolated  = errors.New("Conflict Data for Action, Data Exist")
	ErrForeignViolated = errors.New("Conflict Data for Action, Reference invalid")
	ErrNotNullViolated = errors.New("Invalid Data for Action, Null Value not permitted")
	ErrCheckViolated   = errors.New("Invalid Data, Value is invalid")
)

func IsErrViolated(err error) (string, error) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		// fmt.Println("Cek constraint violated: ", pgErr.ConstraintName)
		switch pgErr.Code {
		case "23505":
			return "", ErrUniqueViolated
		case "23503":
			return "", ErrForeignViolated
		case "23502":
			return "", ErrNotNullViolated
		case "23514":
			return "", ErrCheckViolated
		}
	}
	return "", nil
}
