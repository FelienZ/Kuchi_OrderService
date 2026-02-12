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
	LoggerRepo  models.LoggerService
}

func (s *ProductServiceImpl) Create(p models.Product) error {
	if p.Name == "" || p.Price <= 0 || p.Stock < 0 {
		return ErrProductInvalid
	}
	list := s.ProductRepo.FindAll()
	newProduct := models.Product{
		ID:        "product-" + strconv.Itoa(len(list)+1),
		Name:      p.Name,
		Price:     p.Price,
		Stock:     p.Stock,
		CreatedAt: time.Now().UTC(),
	}
	err := s.ProductRepo.Save(newProduct)
	if err == repository.ErrConflict {
		return ErrProductConflict
	}
	if err != nil {
		return err //kalo ternyata bukan conflict, tapi harusnya hanya ini
	}
	return nil
}

func (s *ProductServiceImpl) GetByID(id string) (models.Product, error) {
	if id == "" {
		return models.Product{}, ErrProductInvalid
	}
	p, err := s.ProductRepo.FindByID(id)
	if err == repository.ErrNotFound {
		return models.Product{}, ErrProductNotFound
	}
	if err != nil {
		return models.Product{}, err
	}
	return p, nil
}

func (s *ProductServiceImpl) List() []models.Product {
	return s.ProductRepo.FindAll()
}

func (s *ProductServiceImpl) Sell(id string, qty int) error {
	if id == "" || qty <= 0 {
		return ErrProductInvalid
	}
	item, err := s.ProductRepo.FindByID(id)
	if err == repository.ErrNotFound {
		return ErrProductNotFound
	}
	if err != nil {
		return err
	}
	if item.Stock < qty {
		return ErrProductNotEnough
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
	err = s.ProductRepo.UpdateStock(id, newStock)
	if err == repository.ErrNotFound {
		return ErrProductNotFound
	}
	if err != nil {
		return err
	}
	return nil
}

func (s *ProductServiceImpl) RecoverStock(id string, qty int) error {
	p, err := s.ProductRepo.FindByID(id)
	if err == repository.ErrNotFound {
		return ErrProductNotFound
	}
	if qty < 0 {
		return ErrProductInvalid
	}
	p.Stock = qty
	err = s.ProductRepo.UpdateStock(id, p.Stock)
	if err == repository.ErrNotFound {
		return ErrProductNotFound
	}
	if err != nil {
		return err
	}
	return nil
}

func (s *ProductServiceImpl) UpdateProductData(id string, u models.UpdateProductRequest) error {
	if id == "" {
		return ErrProductInvalid
	}
	product, err := s.ProductRepo.FindByID(id)
	if err == repository.ErrNotFound {
		return ErrProductNotFound
	}
	if u.Name != nil {
		product.Name = *u.Name
	}
	if u.Price != nil {
		if *u.Price <= 0 {
			return ErrProductInvalid
		}
		product.Price = *u.Price
	}
	if u.Stock != nil {
		if *u.Stock < 0 {
			return ErrProductInvalid
		}
		product.Stock = *u.Stock
	}
	product.UpdatedAt = time.Now().UTC()
	err = s.ProductRepo.Update(product)
	if err == repository.ErrNotFound {
		return ErrProductNotFound
	}
	if err != nil {
		return err
	}
	return nil
}

func (s *ProductServiceImpl) DeleteProduct(id string) error {
	if id == "" {
		return ErrProductInvalid
	}
	err := s.ProductRepo.Delete(id)
	if err == repository.ErrNotFound {
		return ErrProductNotFound
	}
	if err != nil {
		return err
	}
	return nil
}
