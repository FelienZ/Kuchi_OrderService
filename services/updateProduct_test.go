package services

import (
	"go-inventory/internal/testutils"
	"go-inventory/models"
	"testing"
)

func TestUpdateProductSuccess(t *testing.T) {
	// arrange data produk awal
	product := models.Product{
		Name:  "Sebuah Barang",
		Price: 10000,
		Stock: 2,
	}

	// mock repo & service
	mockRepo := testutils.NewProductRepoTestInstance()
	productService := ProductServiceImpl{ProductRepo: mockRepo}

	// act + assert
	if err := productService.Create(product); err != nil {
		t.Fatalf("Unexpected error create: %s", err)
	}
	newData := models.Product{
		Name:  "Gitar",
		Price: 800000,
		Stock: 3,
	}
	if errUpdate := productService.UpdateProductData("product-1", models.UpdateProductRequest{
		Name:  &newData.Name,
		Price: &newData.Price,
		Stock: &newData.Stock,
	}); errUpdate != nil {
		t.Fatalf("Unexpected error update: %s", errUpdate)
	}
	// assert equal untuk update
	updated, err := productService.GetByID("product-1")
	if err != nil {
		t.Fatalf("unexpected error get after update: %s", err)
	}
	if updated.Name != "Gitar" {
		t.Fatalf("Expected updated name Gitar, but new name is: %s", updated.Name)
	}
	if updated.Price != 800000 {
		t.Fatalf("Expected updated price to 800000, but price is %d", updated.Price)
	}
	if updated.Stock != 3 {
		t.Fatalf("Expected updated stock to 3, but stock is: %d", updated.Stock)
	}
}

func TestUpdateProductFail(t *testing.T) {
	// arrange data produk awal
	product := models.Product{
		Name:  "Sebuah Barang",
		Price: 10000,
		Stock: 2,
	}

	// mock repo & service
	mockRepo := testutils.NewProductRepoTestInstance()
	productService := ProductServiceImpl{ProductRepo: mockRepo}

	// act + assert
	if err := productService.Create(product); err != nil {
		t.Fatalf("Unexpected error create: %s", err)
	}
	newData := models.Product{
		Name:  "Gitar",
		Price: 800000,
		Stock: -1,
	}
	if errUpdate := productService.UpdateProductData("product-1", models.UpdateProductRequest{
		Name:  &newData.Name,
		Price: &newData.Price,
		Stock: &newData.Stock,
	}); errUpdate == nil {
		t.Fatalf("expected error update, but got: %s", errUpdate)
	}
}
