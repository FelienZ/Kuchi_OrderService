package http_test

import (
	"go-inventory/models"
)

type ProductServiceImpl struct {
	Product models.Product
	Err     error
}

func NewProductServiceTestInstance(product models.Product, err error) *ProductServiceImpl {
	return &ProductServiceImpl{Product: product, Err: err}
}

func (s *ProductServiceImpl) Create(p models.Product) error {
	return s.Err // mock return errornya
}

func (s *ProductServiceImpl) GetByID(id string) (models.Product, error) {
	return s.Product, s.Err
}

func (s *ProductServiceImpl) List() []models.Product {
	return []models.Product{s.Product}
}

func (s *ProductServiceImpl) Sell(id string, qty int) error {
	return s.Err
}

func (s *ProductServiceImpl) UpdateProductData(id string, u models.UpdateProductRequest) error {
	return s.Err
}

func (s *ProductServiceImpl) DeleteProduct(id string) error {
	return s.Err
}

func (s *ProductServiceImpl) RecoverStock(id string, qty int) error {
	return s.Err
}
