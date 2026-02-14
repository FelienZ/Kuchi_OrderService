package services

import (
	"go-inventory/internal/testutils/service_test"
	"go-inventory/models"
	"testing"
)

func TestDeleteProduct(t *testing.T) {
	product := models.Product{
		Name:  "sebuah",
		Price: 3000,
		Stock: 40,
	}

	repo := service_test.NewProductRepoTestInstance()
	logger := service_test.NewLoggerTestInstance()

	service := ProductServiceImpl{ProductRepo: repo, LoggerRepo: logger}

	if err := service.Create(product); err != nil {
		t.Fatalf("Unexpected error create: %v", err)
	}

	if err := service.DeleteProduct("product-1"); err != nil {
		t.Fatalf("Unexpected error delete: %v", err)
	}

	_, err := service.GetByID("product-1")
	if err != ErrProductNotFound {
		t.Fatalf("expected error not found after delete, but got: %v", err)
	}
}
func TestDeleteProductNotFound(t *testing.T) {
	product := models.Product{
		Name:  "sebuah",
		Price: 3000,
		Stock: 40,
	}

	repo := service_test.NewProductRepoTestInstance()
	logger := service_test.NewLoggerTestInstance()

	service := ProductServiceImpl{ProductRepo: repo, LoggerRepo: logger}

	if err := service.Create(product); err != nil {
		t.Fatalf("Unexpected error create: %v", err)
	}

	if err := service.DeleteProduct("product-2"); err == nil {
		t.Fatalf("expected error not found, but got: %v", err)
	}

	_, err := service.GetByID("product-1")
	if err != nil {
		t.Fatalf("unexpected error get product: %v", err)
	}
}
