package services

import (
	"go-inventory/internal/testutils"
	"go-inventory/models"
	"go-inventory/repository"
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
	mockRepo := testutils.NewProductRepoTestInstance()

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
	mockRepo := testutils.NewProductRepoTestInstance()

	// create service instance
	productService := ProductServiceImpl{ProductRepo: mockRepo}

	// act + assert
	_, errFind := productService.GetByID("product-123")
	if errFind != repository.ErrNotFound {
		t.Fatalf("expected error find but got: %s", errFind)
	}

}
