package products

import (
	"bytes"
	"encoding/json"
	"go-inventory/internal/testutils/http_test"
	"go-inventory/models"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateProduct(t *testing.T) {
	newProduct := models.Product{
		Name:  "new-brand",
		Price: 20000,
		Stock: 2,
	}
	bodyRequest, _ := json.Marshal(newProduct)
	req := httptest.NewRequest(http.MethodPut, "http://localhost:5000/api/products/product-1", bytes.NewBuffer(bodyRequest))
	req.SetPathValue("id", "product-1")
	s := http_test.NewProductServiceTestInstance()
	Httpservice := ProductServiceAPI{Service: s}
	rec := httptest.NewRecorder()
	Httpservice.UpdateProduct(rec, req)

	response := rec.Result()
	body, _ := io.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		t.Fatalf("Expected status ok, but got: %v", response.StatusCode)
	}
	var bodyString map[string]string
	json.Unmarshal(body, &bodyString)
	if bodyString["message"] != "success update product data" {
		t.Fatalf("Expected response message to equal %s, but got %s", "success update product data", bodyString["message"])
	}
}

func TestUpdateProductFail(t *testing.T) {
	newProduct := models.Product{
		Name:  "new-brand",
		Price: 20000,
		Stock: 2,
	}
	bodyRequest, _ := json.Marshal(newProduct)
	req := httptest.NewRequest(http.MethodPut, "http://localhost:5000/api/products/product-1", bytes.NewBuffer(bodyRequest))
	req.SetPathValue("id", "")
	s := http_test.NewProductServiceTestInstance()
	Httpservice := ProductServiceAPI{Service: s}
	rec := httptest.NewRecorder()
	Httpservice.UpdateProduct(rec, req)

	response := rec.Result()
	body, _ := io.ReadAll(response.Body)

	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected status 400, but got: %v", response.StatusCode)
	}
	var bodyString map[string]string
	json.Unmarshal(body, &bodyString)
	if bodyString["message"] != "Invalid Payload for Product Data" {
		t.Fatalf("Expected response message to equal %s, but got %s", "Invalid Payload for Product Data", bodyString["message"])
	}
}

// pr next kalo case find tapi not found
