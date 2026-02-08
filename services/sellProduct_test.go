package services

// skip test sell Product -> logger i/o logic
/* func TestSellProduct(t *testing.T) {
	// assert
	product := models.Product{
		Name:  "Barang",
		Price: 30000,
		Stock: 9,
	}

	// mock
	mockRepo := testutils.NewProductRepoTestInstance()
	logger := repository.NewLoggerInstance()
	productService := ProductServiceImpl{ProductRepo: mockRepo, LoggerRepo: logger}

	// assert + act
	if err := productService.Create(product); err != nil {
		t.Fatalf("Unexpected error create: %s", err)
	}
	sellQty := 4
	if err := productService.Sell("product-1", sellQty); err != nil {
		t.Fatalf("Unexpected error sell: %s", err)
	}
	// validasi
	newData, err := productService.GetByID("product-1")
	if err != nil {
		t.Fatalf("Unexpected error get after selling: %s", err)
	}
	newStock := product.Stock - sellQty
	if newData.Stock != newStock {
		t.Fatalf("expected stock after sell is: %d, but new stock is: %d", newStock, newData.Stock)
	}
}
*/
