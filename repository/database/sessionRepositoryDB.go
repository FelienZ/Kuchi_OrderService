package database

import (
	"context"
	"go-inventory/internal/store"
	"go-inventory/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepositoryInDB struct {
	Db *pgxpool.Pool
}

func NewSessionRepositoryInDB(db *pgxpool.Pool) *SessionRepositoryInDB {
	return &SessionRepositoryInDB{
		Db: db,
	}
}

func (r *SessionRepositoryInDB) Create(ctx context.Context, tx store.DBExecutor, s models.UserSession) (string, error) {
	var id string
	err := tx.QueryRow(ctx, "INSERT INTO user_session (id, user_id, expired_at) VALUES ($1,$2,$3) RETURNING id",
		s.ID, s.UserID, s.ExpiresAt).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *SessionRepositoryInDB) FindByID(ctx context.Context, tx store.DBExecutor, id string) (models.UserSession, error) {
	var s models.UserSession
	err := tx.QueryRow(ctx, "SELECT id, user_id, expired_at FROM user_session WHERE id=$1",
		id).Scan(&s.ID, &s.UserID, &s.ExpiresAt)
	if err == pgx.ErrNoRows {
		return models.UserSession{}, ErrNoRows
	}
	if err != nil {
		return models.UserSession{}, err
	}
	return s, nil
}

func (r *SessionRepositoryInDB) Delete(ctx context.Context, tx store.DBExecutor, id string) error {
	d, err := tx.Exec(ctx, "DELETE FROM user_session WHERE id=$1", id)
	if err != nil {
		return err
	}
	if d.RowsAffected() == 0 {
		return ErrNoDelete
	}
	return nil
}
