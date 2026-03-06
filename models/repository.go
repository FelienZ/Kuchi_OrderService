package models

import (
	"context"
	"go-inventory/internal/store"
)

type ProductRepository interface {
	FindByID(ctx context.Context, tx store.DBExecutor, id string) (Product, error)
	FindAll(ctx context.Context, tx store.DBExecutor) ([]Product, error)
	Save(ctx context.Context, tx store.DBExecutor, p Product) (string, error)
	Update(ctx context.Context, tx store.DBExecutor, p Product) error
	Delete(ctx context.Context, tx store.DBExecutor, id string) error
	DecreaseStock(ctx context.Context, tx store.DBExecutor, id string, qty int) error
}

type OrderRepository interface {
	FindByID(ctx context.Context, tx store.DBExecutor, id string) (Order, error)
	FindAll(ctx context.Context, tx store.DBExecutor, filter GetOrderParameter) ([]Order, error)
	FindByUserID(ctx context.Context, tx store.DBExecutor, userid string) ([]Order, error)
	Save(ctx context.Context, tx store.DBExecutor, order Order) (string, error)
	Update(ctx context.Context, tx store.DBExecutor, order Order) error
	Delete(ctx context.Context, tx store.DBExecutor, id string) error
	UpdateOrderStatus(ctx context.Context, tx store.DBExecutor, prevStatus Status, newStatus Status, id string) error
}

type UserRepository interface {
	FindByID(ctx context.Context, tx store.DBExecutor, id string) (User, error)
	FindByEmail(ctx context.Context, tx store.DBExecutor, email string) (User, error)
	FindByUsername(ctx context.Context, tx store.DBExecutor, username string) (User, error)
	FindAll(ctx context.Context, tx store.DBExecutor) ([]User, error)
	Create(ctx context.Context, tx store.DBExecutor, u User) (string, error)
	Update(ctx context.Context, tx store.DBExecutor, u User) error
	Delete(ctx context.Context, tx store.DBExecutor, id string) error
}

type SessionRepository interface {
	Create(ctx context.Context, tx store.DBExecutor, s UserSession) (string, error)
	FindByID(ctx context.Context, tx store.DBExecutor, id string) (UserSession, error)
	Delete(ctx context.Context, tx store.DBExecutor, id string) error
}
