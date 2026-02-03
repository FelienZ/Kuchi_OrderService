package models

type ProductRepository interface {
	FindByID(id string) (Product, error)
	FindAll() []Product
	Save(product Product) error
	Delete(id string) error
	UpdateStock(id string, newStock int) error
}

type OrderRepository interface {
	FindByID(id string) (Order, error)
	FindAll() []Order
	FindByUserID(userid string) []Order
	Save(order Order) error
	Update(order Order) error
	Delete(id string) error
}
