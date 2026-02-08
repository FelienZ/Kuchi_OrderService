package services

import (
	"go-inventory/internal/testutils"
	"testing"
)

func TestGetListProduct(t *testing.T) {
	// mock
	mockRepo := testutils.NewProductRepoTestInstance()
	productService := ProductServiceImpl{ProductRepo: mockRepo}

	lists := productService.List()
	if lists == nil {
		t.Fatalf("Expected Slice of Products, got: %v", lists)
	}
}
