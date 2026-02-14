package services

import (
	"go-inventory/internal/testutils/service_test"
	"go-inventory/models"
	"testing"
)

func TestGetProductSuccess(t *testing.T) {
	// arrange produk (id created by service)
	product := models.Product{
		Name:  "Sebuah Produk",
		Price: 100000,
		Stock: 2,
	}

	//deps
	mockRepo := service_test.NewProductRepoTestInstance()

	// create service instance
	productService := ProductServiceImpl{ProductRepo: mockRepo}

	// act + assert
	err := productService.Create(product)
	if err != nil {
		t.Fatalf("Unexpected error create: %s", err)
	}
	// next bikin return id
	_, errFind := productService.GetByID("product-1")
	if errFind != nil {
		t.Fatalf("Unexpected error find: %s", errFind)
	}

}

func TestGetProductNotFound(t *testing.T) {
	// skip arrange produk

	//deps
	mockRepo := service_test.NewProductRepoTestInstance()

	// create service instance
	productService := ProductServiceImpl{ProductRepo: mockRepo}

	// act + assert
	_, errFind := productService.GetByID("product-123")
	if errFind != ErrProductNotFound {
		t.Fatalf("expected error find but got: %s", errFind)
	}

}
