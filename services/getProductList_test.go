package services

import (
	"go-inventory/internal/testutils/service_test"
	"testing"
)

func TestGetListProduct(t *testing.T) {
	// mock
	mockRepo := service_test.NewProductRepoTestInstance()
	productService := ProductServiceImpl{ProductRepo: mockRepo}

	lists := productService.List()
	if lists == nil {
		t.Fatalf("Expected Slice of Products, got: %v", lists)
	}
}
