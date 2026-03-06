package database

import (
	"context"
	"fmt"
	"go-inventory/internal/store"
	"go-inventory/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepositoryInDB struct {
	Db *pgxpool.Pool
}

func NewOrderRepositoryInstanceInDB(db *pgxpool.Pool) *OrderRepositoryInDB {
	return &OrderRepositoryInDB{
		Db: db,
	}
}
func (r *OrderRepositoryInDB) FindByID(ctx context.Context, tx store.DBExecutor, id string) (models.Order, error) {
	var order models.Order
	var item models.OrderItemFromDB
	d, err := tx.Query(ctx, "SELECT o.id, o.user_id, o.status, o.created_at, o.updated_at, oi.product_id, oi.quantity, oi.price FROM orders o LEFT JOIN order_items oi ON o.id=oi.order_id WHERE o.id=$1",
		id)
	if err != nil {
		return models.Order{}, err
	}
	defer d.Close()
	var found bool
	for d.Next() {
		found = true
		if errScan := d.Scan(&order.ID, &order.UserID, &order.Status, &order.CreatedAt, &order.UpdatedAt,
			&item.ProductID, &item.Qty, &item.TotalPrice); errScan != nil {
			return models.Order{}, errScan
		}
		if item.ProductID != nil && item.Qty != nil && item.TotalPrice != nil {
			order.Item = append(order.Item, models.OrderItem{ProductID: *item.ProductID,
				Qty: *item.Qty, TotalPrice: *item.TotalPrice})
		}
	}
	if !found {
		return models.Order{}, ErrNoRows
	}
	if err := d.Err(); err != nil {
		return models.Order{}, err
	}
	return order, nil
}

func (r *OrderRepositoryInDB) FindAll(ctx context.Context, tx store.DBExecutor,
	filter models.GetOrderParameter) ([]models.Order, error) {
	var id string
	var listID []string
	var order models.Order
	var item models.OrderItemFromDB
	queryArgs := []any{}
	idx := 1
	baseQuery := "SELECT id FROM orders"
	if filter.Status != nil {
		// fmt.Println("status tidak nil, value: ", filter.Status)
		baseQuery += fmt.Sprintf(" WHERE status=$%d::status", idx) // harus match type enum db
		// stringStatus := fmt.Sprintf("'%s'", filter.Status)
		queryArgs = append(queryArgs, filter.Status)
		idx++
	}
	queryArgs = append(queryArgs, filter.Limit, filter.Offset)
	baseQuery += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", idx, idx+1)
	// fmt.Println("cek Query Lengkap: ", baseQuery)
	d, err := tx.Query(ctx, baseQuery, queryArgs...) // id
	if err != nil {
		return nil, err
	}
	defer d.Close()
	list := make(map[string]*models.Order)
	for d.Next() {
		if err := d.Scan(&id); err != nil {
			return nil, err
		}
		// fmt.Println("cek id scan: ", id)
		listID = append(listID, id)
	}
	// fmt.Println("Cek list ID Order: ", listID)
	if err := d.Err(); err != nil {
		return nil, err
	}
	if len(listID) == 0 {
		return []models.Order{}, nil // alias orderItem kosong
	}
	secondQuery := "SELECT o.id, o.user_id, o.status, o.created_at, o.updated_at, oi.product_id, oi.quantity, oi.price FROM orders o LEFT JOIN order_items oi ON o.id=oi.order_id WHERE o.id=ANY($1::uuid[]) ORDER BY o.created_at DESC"
	s, errSec := tx.Query(ctx, secondQuery, listID) // detail + item
	if errSec != nil {
		return nil, errSec
	}
	defer s.Close()
	for s.Next() {
		if err := s.Scan(&order.ID, &order.UserID, &order.Status, &order.CreatedAt, &order.UpdatedAt,
			&item.ProductID, &item.Qty, &item.TotalPrice); err != nil {
			return nil, err
		}
		if _, ok := list[order.ID]; !ok {
			list[order.ID] = &models.Order{
				ID: order.ID, UserID: order.UserID, Item: []models.OrderItem{}, Status: order.Status,
				CreatedAt: order.CreatedAt, UpdatedAt: order.UpdatedAt,
			}
		}
		if item.ProductID != nil && item.Qty != nil && item.TotalPrice != nil {
			list[order.ID].Item = append(list[order.ID].Item, models.OrderItem{ProductID: *item.ProductID, Qty: *item.Qty, TotalPrice: *item.TotalPrice})
		} // prevent append nil dari null value
	}
	//translate ke slice
	var result []models.Order
	for _, v := range listID {
		result = append(result, *list[v])
	}
	// fmt.Println("Cek result all order: ", result)
	return result, nil
}

func (r *OrderRepositoryInDB) FindByUserID(ctx context.Context, tx store.DBExecutor, userid string) ([]models.Order, error) {
	d, err := tx.Query(ctx, "SELECT o.id, o.user_id, o.status, o.created_at, o.updated_at, oi.product_id, oi.quantity, oi.price FROM orders o LEFT JOIN order_items oi ON o.id=oi.order_id WHERE o.user_id=$1",
		userid)
	if err != nil {
		return nil, err
	}
	defer d.Close()
	list := make(map[string]*models.Order)
	var found bool
	for d.Next() {
		found = true
		var order models.Order
		var item models.OrderItemFromDB
		if err := d.Scan(&order.ID, &order.UserID, &order.Status, &order.CreatedAt, &order.UpdatedAt,
			&item.ProductID, &item.Qty, &item.TotalPrice); err != nil {
			return nil, err
		}
		if _, ok := list[order.ID]; !ok {
			list[order.ID] = &models.Order{
				ID: order.ID, UserID: order.UserID, Item: []models.OrderItem{}, Status: order.Status,
				CreatedAt: order.CreatedAt, UpdatedAt: order.UpdatedAt,
			}
		}
		if item.ProductID != nil && item.Qty != nil && item.TotalPrice != nil {
			list[order.ID].Item = append(list[order.ID].Item, models.OrderItem{ProductID: *item.ProductID,
				Qty: *item.Qty, TotalPrice: *item.TotalPrice})
		}
	}
	if !found {
		return []models.Order{}, nil
	}
	if err := d.Err(); err != nil {
		return nil, err
	}
	var result []models.Order
	for _, v := range list {
		result = append(result, *v)
	}
	return result, nil
}

func (r *OrderRepositoryInDB) Save(ctx context.Context, tx store.DBExecutor, o models.Order) (string, error) {
	var id string
	err := tx.QueryRow(ctx, "INSERT INTO orders (id, user_id, status, created_at, updated_at) VALUES ($1,$2,$3,$4,$5) RETURNING id",
		o.ID, o.UserID, o.Status, o.CreatedAt, o.UpdatedAt).Scan(&id)
	if err != nil {
		return "", err
	}
	for _, v := range o.Item {
		_, errItem := tx.Exec(ctx, "INSERT INTO order_items (order_id, product_id, quantity, price) VALUES ($1,$2,$3,$4)",
			o.ID, v.ProductID, v.Qty, v.TotalPrice)
		if errItem != nil {
			return "", errItem
		}
	}
	return id, nil
}

func (r *OrderRepositoryInDB) Delete(ctx context.Context, tx store.DBExecutor, id string) error {
	// order_id ondelete cascade
	d, err := tx.Exec(ctx, "DELETE FROM orders WHERE id=$1", id)
	if err != nil {
		return err
	}
	if d.RowsAffected() == 0 {
		return ErrNoDelete
	}
	return nil
}

func (r *OrderRepositoryInDB) Update(ctx context.Context, tx store.DBExecutor, o models.Order) error {
	_, err := tx.Exec(ctx, "DELETE FROM order_items WHERE order_id=$1", o.ID) //order_item boleh tidak ada
	if err != nil {
		return err
	}
	for _, v := range o.Item {
		_, err := tx.Exec(ctx, "INSERT INTO order_items (order_id, product_id, quantity, price) VALUES ($1,$2,$3,$4)",
			o.ID, v.ProductID, v.Qty, v.TotalPrice)
		if err != nil {
			return err
		}
	}
	u, errUpdate := tx.Exec(ctx, "UPDATE orders SET status=$1, updated_at=$2 WHERE id=$3", o.Status, o.UpdatedAt, o.ID)
	if errUpdate != nil {
		return errUpdate
	}
	if u.RowsAffected() == 0 {
		return ErrNoUpdate
	}
	return nil
}

func (r *OrderRepositoryInDB) UpdateOrderStatus(ctx context.Context, tx store.DBExecutor, prevStatus models.Status,
	newStatus models.Status, id string) error {
	u, err := tx.Exec(ctx, "UPDATE orders SET status=$3, updated_at=CURRENT_TIMESTAMP WHERE id=$1 AND status=$2", id, prevStatus, newStatus)
	if err != nil {
		return err
	}
	if u.RowsAffected() == 0 {
		return ErrNoUpdate
	}
	return nil
}
