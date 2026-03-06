package database

import (
	"context"
	"go-inventory/internal/store"
	"go-inventory/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepositoryInDB struct {
	Db *pgxpool.Pool
}

func NewProductRepositoryInDBInstance(db *pgxpool.Pool) *ProductRepositoryInDB {
	return &ProductRepositoryInDB{
		Db: db,
	}
}

func (r *ProductRepositoryInDB) FindByID(ctx context.Context, tx store.DBExecutor, id string) (models.Product, error) {
	var p models.Product
	err := tx.QueryRow(ctx, "SELECT id, name, price, stock, created_at, updated_at FROM products WHERE id=$1", id).Scan(&p.ID, &p.Name,
		&p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt)
	if err == pgx.ErrNoRows {
		return models.Product{}, ErrNoRows
	}
	if err != nil {
		return models.Product{}, err
	}
	return p, nil
}

func (r *ProductRepositoryInDB) FindAll(ctx context.Context, tx store.DBExecutor) ([]models.Product, error) {
	var list []models.Product
	rows, err := tx.Query(ctx, "SELECT id, name, price, stock, created_at, updated_at FROM products ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, rows.Err()
}

func (r *ProductRepositoryInDB) Save(ctx context.Context, tx store.DBExecutor, p models.Product) (string, error) {
	var id string
	// fmt.Println("cek product: ", p)
	err := tx.QueryRow(ctx, "INSERT INTO products (id, name, price, stock, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6) RETURNING id",
		p.ID, p.Name, p.Price, p.Stock, p.CreatedAt, p.UpdatedAt).Scan(&id)
	if err != nil {
		if _, errViolation := IsErrViolated(err); errViolation != nil {
			return "", errViolation
		}
		// fmt.Println("cek err: ", err)
		return "", err
	}
	return id, nil
}

func (r *ProductRepositoryInDB) Update(ctx context.Context, tx store.DBExecutor, p models.Product) error {
	u, err := tx.Exec(ctx, "UPDATE products SET name=$1, price=$2, stock=$3, updated_at=$4 WHERE id = $5",
		p.Name, p.Price, p.Stock, p.UpdatedAt, p.ID)
	if err != nil {
		if _, errViolation := IsErrViolated(err); errViolation != nil {
			return errViolation
		}
		return err
	}
	if u.RowsAffected() == 0 {
		return ErrNoUpdate
	}
	return nil
}

func (r *ProductRepositoryInDB) Delete(ctx context.Context, tx store.DBExecutor, id string) error {
	d, err := tx.Exec(ctx, "DELETE FROM products WHERE id=$1", id)
	if err != nil {
		return err
	}
	if d.RowsAffected() == 0 {
		return ErrNoDelete
	}
	return nil
}

func (r *ProductRepositoryInDB) DecreaseStock(ctx context.Context, tx store.DBExecutor, id string, qty int) error {
	d, err := tx.Exec(ctx, "UPDATE products SET stock=stock-$2, updated_at=CURRENT_TIMESTAMP WHERE id=$1 AND stock>=$2", id, qty)
	if err != nil {
		if _, errViolation := IsErrViolated(err); errViolation != nil {
			return errViolation
		}
		return err
	}
	if d.RowsAffected() == 0 {
		return ErrNoUpdate
	}
	return nil
}
