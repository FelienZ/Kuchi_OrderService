package testutils

import (
	"go-inventory/models"
)

type ProductServiceImpl struct{}

func NewProductServiceTestInstance() *ProductServiceImpl {
	return &ProductServiceImpl{}
}

func (s *ProductServiceImpl) Create(p models.Product) error {
	if p.Name == "" || p.Price <= 0 || p.Stock < 0 {
		return ErrProductInvalid
	}
	return nil
}

func (s *ProductServiceImpl) GetByID(id string) (models.Product, error) {
	if id == "" {
		return models.Product{}, ErrProductInvalid
	}
	return models.Product{}, nil
}

func (s *ProductServiceImpl) List() []models.Product {
	return []models.Product{}
}

func (s *ProductServiceImpl) Sell(id string, qty int) error {
	if id == "" || qty <= 0 {
		return ErrProductInvalid
	}
	return nil
}

func (s *ProductServiceImpl) UpdateProductData(id string, u models.UpdateProductRequest) error {
	if id == "" {
		return ErrProductInvalid
	}
	if u.Name != nil {
		if *u.Name == "" {
			return ErrProductInvalid
		}
	}
	if u.Price != nil {
		if *u.Price <= 0 {
			return ErrProductInvalid
		}
	}
	if u.Stock != nil {
		if *u.Stock < 0 {
			return ErrProductInvalid
		}
	}
	return nil
}

func (s *ProductServiceImpl) DeleteProduct(id string) error {
	if id == "" {
		return ErrProductInvalid
	}
	return nil
}

func (s *ProductServiceImpl) RecoverStock(id string, qty int) error {
	if id == "" || qty < 0 {
		return ErrProductInvalid
	}
	return nil
}
