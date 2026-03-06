package database

import (
	"context"
	"fmt"
	"go-inventory/internal/store"
	"go-inventory/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepositoryInDB struct {
	Db *pgxpool.Pool
}

func NewUserRepoInDBInstance(db *pgxpool.Pool) *UserRepositoryInDB {
	return &UserRepositoryInDB{
		Db: db,
	}
}

func (r *UserRepositoryInDB) FindByID(ctx context.Context, tx store.DBExecutor, id string) (models.User, error) {
	var u models.User
	err := tx.QueryRow(ctx, "SELECT id, user_name, email, password, role, created_at, updated_at FROM users WHERE id=$1", id).Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err == pgx.ErrNoRows {
		return models.User{}, ErrNoRows
	}
	if err != nil {
		return models.User{}, err
	}
	return u, nil
}

func (r *UserRepositoryInDB) FindByEmail(ctx context.Context, tx store.DBExecutor, email string) (models.User, error) {
	var u models.User
	err := tx.QueryRow(ctx, "SELECT id, user_name, email, password, role, created_at, updated_at FROM users WHERE email=$1", email).Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err == pgx.ErrNoRows {
		return models.User{}, ErrNoRows
	}
	if err != nil {
		return models.User{}, err
	}
	return u, nil
}

func (r *UserRepositoryInDB) FindByUsername(ctx context.Context, tx store.DBExecutor, username string) (models.User, error) {
	var u models.User
	err := tx.QueryRow(ctx, "SELECT id, user_name, email, password, role, created_at, updated_at FROM users WHERE user_name=$1", username).Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err == pgx.ErrNoRows {
		return models.User{}, ErrNoRows
	}
	if err != nil {
		return models.User{}, err
	}
	return u, nil
}

func (r *UserRepositoryInDB) FindAll(ctx context.Context, tx store.DBExecutor) ([]models.User, error) {
	var list []models.User
	d, err := tx.Query(ctx, "SELECT id, user_name, email, password, role, created_at, updated_at FROM users")
	if err != nil {
		return nil, err
	}
	defer d.Close()
	for d.Next() {
		var u models.User
		if err := d.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, u)
	}
	return list, d.Err()
}

func (r *UserRepositoryInDB) Create(ctx context.Context, tx store.DBExecutor, u models.User) (string, error) {
	var id string
	err := tx.QueryRow(ctx, "INSERT INTO users (id, user_name, email, password, role, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		u.ID, u.Username, u.Email, u.Password, u.Role, u.CreatedAt, u.UpdatedAt).Scan(&id)
	if err != nil {
		if _, errViolated := IsErrViolated(err); errViolated != nil {
			return "", errViolated
		}
		return "", err
	}
	return id, nil
}

func (r *UserRepositoryInDB) Update(ctx context.Context, tx store.DBExecutor, u models.User) error {
	d, err := tx.Exec(ctx, "UPDATE users SET user_name=$1, email=$2, password=$3, role=$4, updated_at=$5 WHERE id=$6",
		u.Username, u.Email, u.Password, u.Role, u.UpdatedAt, u.ID)
	if err != nil {
		if _, errViolated := IsErrViolated(err); errViolated != nil {
			return errViolated
		}
		return err
	}
	if d.RowsAffected() == 0 {
		fmt.Println("rows affect 0")
		return ErrNoUpdate
	}
	return nil
}

func (r *UserRepositoryInDB) Delete(ctx context.Context, tx store.DBExecutor, id string) error {
	d, err := tx.Exec(ctx, "DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return err
	}
	if d.RowsAffected() == 0 {
		return ErrNoDelete
	}
	return nil
}
