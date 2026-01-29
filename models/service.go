package models

type ProductService interface {
	Create(p Product) error
	GetByID(id string) (Product, error)
	List() []Product
	Sell(id string, qty int) error
}

type OrderService interface {
	GetByID(id string) (Order, error)
	List() []Order
	GetByUserID(userid string) ([]Order, error)
	CreateOrder(o Order) error
	PayOrder(id string) error
	CancelOrder(id string) error
}
