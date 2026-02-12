package services

import (
	"go-inventory/internal/testutils"
	"go-inventory/models"
	"testing"
)

// skip test sell Product -> logger i/o storage (sekarang pakai mock memory)
func TestSellProduct(t *testing.T) {
	// assert
	product := models.Product{
		Name:  "Barang",
		Price: 30000,
		Stock: 9,
	}

	// mock
	mockRepo := testutils.NewProductRepoTestInstance()
	logger := testutils.NewLoggerTestInstance()
	productService := ProductServiceImpl{ProductRepo: mockRepo, LoggerRepo: logger}

	// assert + act
	if err := productService.Create(product); err != nil {
		t.Fatalf("Unexpected error create: %v", err)
	}
	sellQty := 4
	if err := productService.Sell("product-1", sellQty); err != nil {
		t.Fatalf("Unexpected error sell: %v", err)
	}
	// validasi
	newData, err := productService.GetByID("product-1")
	if err != nil {
		t.Fatalf("Unexpected error get after selling: %v", err)
	}
	newStock := product.Stock - sellQty
	if newData.Stock != newStock {
		t.Fatalf("expected stock after sell is: %d, but new stock is: %d", newStock, newData.Stock)
	}
}

func TestSellProductNotFound(t *testing.T) {
	//skip bikin langsung 404
	productRepo := testutils.NewProductRepoTestInstance()
	logger := testutils.NewLoggerTestInstance()
	productService := ProductServiceImpl{ProductRepo: productRepo, LoggerRepo: logger}

	err := productService.Sell("product-1", 3)
	if err != ErrProductNotFound {
		t.Fatalf("Expected Error Not Found but Got: %v", err)
	}
}

func TestSellProductNotEnoughStock(t *testing.T) {
	// arrange
	product := models.Product{
		Name:  "apa aja",
		Stock: 2,
		Price: 10000,
	}

	repo := testutils.NewProductRepoTestInstance()
	logger := testutils.NewLoggerTestInstance()

	productService := ProductServiceImpl{ProductRepo: repo, LoggerRepo: logger}

	//act + assert
	if err := productService.Create(product); err != nil {
		t.Fatalf("unexpected error create: %v", err)
	}

	err := productService.Sell("product-1", 5)
	if err != ErrProductNotEnough {
		t.Fatalf("Expected error not enough, but got: %v", err)
	}
}
