package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPGConnection(ctx context.Context, connection string) (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(ctx, connection)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}
	// defer conn.Close()
	return conn, nil
}
