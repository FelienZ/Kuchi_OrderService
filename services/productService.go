package services

import (
	"fmt"
	"go-inventory/models"
	"go-inventory/repository"
	"time"

	"github.com/google/uuid"
)

type ProductServiceImpl struct {
	ProductRepo *repository.ProductInMemo
	LoggerRepo  *repository.LogInMemory
}

func (s *ProductServiceImpl) Create(p models.Product) error {
	if p.ID == "" || p.Name == "" || p.Price <= 0 || p.Stock < 0 {
		return ErrInvalid
	}
	newProduct := models.Product{
		ID:        p.ID,
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
		ID:       "log-" + uuid.NewString(),
		EntityID: item.ID,
		Entity:   models.PRODUCT,
		Action:   "SELL",
		Note:     fmt.Sprintf("Action:%s - Status:%s - Time:%s ", "Product Out", "Success", time.Now().Format(time.DateTime)),
	})
	if errLog != nil {
		fmt.Println(errLog.Error())
	}
	return s.ProductRepo.UpdateStock(id, newStock)
}
