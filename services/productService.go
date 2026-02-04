package services

import (
	"fmt"
	"go-inventory/models"
	"go-inventory/repository"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type ProductServiceImpl struct {
	ProductRepo models.ProductRepository
	LoggerRepo  *repository.LogInMemory
}

func (s *ProductServiceImpl) Create(p models.Product) error {
	if p.Name == "" || p.Price <= 0 || p.Stock < 0 {
		return ErrInvalid
	}
	list := s.ProductRepo.FindAll()
	newProduct := models.Product{
		ID:        "product-" + strconv.Itoa(len(list)+1),
		Name:      p.Name,
		Price:     p.Price,
		Stock:     p.Stock,
		CreatedAt: time.Now().UTC(),
	}
	return s.ProductRepo.Save(newProduct)
}

func (s *ProductServiceImpl) GetByID(id string) (models.Product, error) {
	if id == "" {
		return models.Product{}, ErrInvalid
	}
	return s.ProductRepo.FindByID(id)
}

func (s *ProductServiceImpl) List() []models.Product {
	return s.ProductRepo.FindAll()
}

func (s *ProductServiceImpl) Sell(id string, qty int) error {
	if id == "" || qty <= 0 {
		return ErrInvalid
	}
	item, err := s.ProductRepo.FindByID(id)
	if err != nil {
		return err
	}
	if item.Stock < qty {
		return ErrNotEnough
	}
	newStock := item.Stock - qty
	errLog := s.LoggerRepo.CreateLog(models.TransactionLog{
		ID:        "log-" + uuid.NewString(),
		EntityID:  item.ID,
		Entity:    models.PRODUCT,
		Action:    "SELL",
		CreatedAt: time.Now().UTC(),
		Note:      fmt.Sprintf("Action:%s - Status:%s - Time:%s ", "Product Out", "Success", time.Now().UTC().Format(time.DateTime)),
	})
	if errLog != nil {
		fmt.Println(errLog.Error())
	}
	return s.ProductRepo.UpdateStock(id, newStock)
}

func (s *ProductServiceImpl) RecoverStock(id string, qty int) error {
	p, err := s.ProductRepo.FindByID(id)
	if err != nil {
		return err
	}
	p.Stock = qty
	return s.ProductRepo.UpdateStock(id, p.Stock)
}

func (s *ProductServiceImpl) UpdateProductData(id string, u models.UpdateProductRequest) error {
	if id == "" {
		return ErrInvalid
	}
	product, err := s.ProductRepo.FindByID(id)
	if err != nil {
		return err
	}
	if u.Name != nil {
		product.Name = *u.Name
	}
	if u.Price != nil {
		if *u.Price <= 0 {
			return ErrInvalid
		}
		product.Price = *u.Price
	}
	if u.Stock != nil {
		if *u.Stock < 0 {
			return ErrInvalid
		}
		product.Stock = *u.Stock
	}
	product.UpdatedAt = time.Now().UTC()
	return s.ProductRepo.Update(product)

}

func (s *ProductServiceImpl) DeleteProduct(id string) error {
	if id == "" {
		return ErrInvalid
	}
	return s.ProductRepo.Delete(id)
}
