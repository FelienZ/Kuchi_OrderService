package models

type ProductRepository interface {
	FindByID(id string) (Product, error)
	FindAll() []Product
	Save(product Product) error
	UpdateStock(id string, newStock int) error
}
