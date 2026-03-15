package models

import "context"

type ProductService interface {
	Create(ctx context.Context, p Product) (string, error)
	GetByID(ctx context.Context, id string) (Product, error)
	List(ctx context.Context) ([]Product, error)
	UpdateProductData(ctx context.Context, id string, u UpdateProductRequest) error
	DeleteProduct(ctx context.Context, id string) error
}

type OrderService interface {
	GetByID(ctx context.Context, id string) (Order, error)
	List(ctx context.Context, param GetOrderParameter) ([]Order, error)
	GetByUserID(ctx context.Context, userid string) ([]Order, error)
	CreateOrder(ctx context.Context, o OrderPayload) (string, error)
	PayOrder(ctx context.Context, id string) error
	CancelOrder(ctx context.Context, id string) error
	DeleteOrder(ctx context.Context, id string) error
}

type LoggerService interface {
	CreateLog(ctx context.Context, l TransactionLog) error
	GetLogs(ctx context.Context) []TransactionLog
}

type UserService interface {
	GetByID(ctx context.Context, id string) (User, error)
	RegisterUser(ctx context.Context, u RegisterRequest) error
	UpdateUser(ctx context.Context, id string, u UserUpdateRequest) error
	DeleteUser(ctx context.Context, id string) error
}

type UserSessionService interface {
	Login(ctx context.Context, l LoginRequest) (LoginResult, error)
	/* Logout(ctx context.Context, sessionID string) error
	ValidateSession(ctx context.Context, sessionID string) (UserIdentity, error) */
}
