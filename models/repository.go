package models

type ProductRepository interface {
	FindByID(id string) (Product, error)
	FindAll() []Product
	Save(product Product) error
	Update(product Product) error
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

type UserRepository interface {
	FindByID(id string) (User, error)
	FindByEmail(email string) (User, error)
	FindByUsername(username string) (User, error)
	FindAll() []User
	Create(u User) error
	Update(u User) error
	Delete(id string) error
}
