package services

import (
	"context"
	"go-inventory/models"
	"go-inventory/repository/database"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductServiceImpl struct {
	Db          *pgxpool.Pool
	ProductRepo models.ProductRepository
	LoggerRepo  models.LoggerService
}

func (s *ProductServiceImpl) Create(ctx context.Context, p models.Product) (string, error) {
	if p.Name == "" || p.Price <= 0 || p.Stock < 0 {
		return "", ErrProductInvalid
	}
	newProduct := models.Product{
		ID:        uuid.NewString(),
		Name:      p.Name,
		Price:     p.Price,
		Stock:     p.Stock,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	id, err := s.ProductRepo.Save(ctx, s.Db, newProduct)
	if errViolation := ErrorProductDomainTranslator(err); errViolation != nil {
		return "", errViolation
	}
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *ProductServiceImpl) GetByID(ctx context.Context, id string) (models.Product, error) {
	if id == "" {
		return models.Product{}, ErrProductInvalid
	}
	p, err := s.ProductRepo.FindByID(ctx, s.Db, id)
	// fmt.Println("Cek err findProduct By ID: ", err)
	if err == database.ErrNoRows {
		return models.Product{}, ErrProductNotFound
	}
	if err != nil {
		return models.Product{}, err
	}
	return p, nil
}

func (s *ProductServiceImpl) List(ctx context.Context) ([]models.Product, error) {
	return s.ProductRepo.FindAll(ctx, s.Db)
}

func (s *ProductServiceImpl) UpdateProductData(ctx context.Context, id string, u models.UpdateProductRequest) error {
	if id == "" {
		return ErrProductInvalid
	}
	product, err := s.ProductRepo.FindByID(ctx, s.Db, id)
	if err == database.ErrNoRows {
		return ErrProductNotFound
	}
	if err != nil {
		return err
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
	err = s.ProductRepo.Update(ctx, s.Db, product)
	if errViolation := ErrorProductDomainTranslator(err); errViolation != nil {
		return errViolation
	}
	if err != nil {
		return err
	}
	return nil
}

func (s *ProductServiceImpl) DeleteProduct(ctx context.Context, id string) error {
	if id == "" {
		return ErrProductInvalid
	}
	err := s.ProductRepo.Delete(ctx, s.Db, id)
	if err != nil {
		return err
	}
	return nil
}
