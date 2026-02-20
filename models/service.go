package models

type ProductService interface {
	Create(p Product) error
	GetByID(id string) (Product, error)
	List() []Product
	Sell(id string, qty int) error
	RecoverStock(id string, qty int) error
	UpdateProductData(id string, u UpdateProductRequest) error
	DeleteProduct(id string) error
}

type OrderService interface {
	GetByID(id string) (Order, error)
	List(GetOrderParameter) []Order
	GetByUserID(userid string) ([]Order, error)
	CreateOrder(o Order) error
	PayOrder(id string) error
	CancelOrder(id string) error
	DeleteOrder(id string) error
}

type LoggerService interface {
	CreateLog(l TransactionLog) error
	GetLogs() []TransactionLog
}

type UserService interface {
	RegisterUser(u RegisterRequest) error
	UpdateUser(id string, u UserUpdateRequest) error
	DeleteUser(id string) error
}

type UserSessionService interface {
	Login(l LoginRequest) (string, error)
	Logout(sessionID string) error
	ValidateSession(sessionID string) (string, error)
}
