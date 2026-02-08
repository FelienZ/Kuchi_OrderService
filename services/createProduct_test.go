package services

import (
	"go-inventory/internal/testutils"
	"go-inventory/models"
	"testing"
)

// Create Product Test

func TestCreateProductSuccess(t *testing.T) {
	// arrange repo
	mockRepo := testutils.NewProductRepoTestInstance()
	productService := ProductServiceImpl{ProductRepo: mockRepo}
	product := models.Product{
		Name:  "Buku",
		Price: 10000,
		Stock: 2,
	} //valid
	err := productService.Create(product)
	if err != nil {
		t.Fatalf("Unexpected Error Create: %s", err.Error())
	}
}
func TestCreateProductFail(t *testing.T) {
	// arrange repo
	mockRepo := testutils.NewProductRepoTestInstance()
	productService := ProductServiceImpl{ProductRepo: mockRepo}
	product2 := models.Product{
		Name:  "Barang error",
		Price: 0,
		Stock: -1,
	}
	// act
	err := productService.Create(product2)
	// assert
	if err != ErrProductInvalid {
		t.Fatalf("Expected Error Invalid")
	}
}
