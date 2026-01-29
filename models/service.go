package models

type ProductService interface {
	Create(p Product) error
	GetByID(id string) (Product, error)
	List() []Product
	Sell(id string, qty int) error
}
