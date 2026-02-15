package products

import (
	"encoding/json"
	"go-inventory/internal/testutils/http_test"
	"go-inventory/models"
	"go-inventory/services"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetProductById(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://localhost:5000/api/product/prouct-1", nil)
	req.SetPathValue("id", "product-1")
	// assume udah ada
	s := http_test.NewProductServiceTestInstance(models.Product{
		Name:  "produk",
		Price: 2000,
		Stock: 2,
	}, nil)
	HttpService := ProductServiceAPI{Service: s}
	rec := httptest.NewRecorder()
	HttpService.GetProductById(rec, req)
	response := rec.Result()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("Unexpected error get by id: %v", response.StatusCode)
	}
	body, _ := io.ReadAll(response.Body)
	var bodyString map[string]string
	json.Unmarshal(body, &bodyString)
	if _, ok := bodyString["data"]; !ok {
		t.Fatalf("Expected response body to have product data field")
	}
}
func TestGetProductByIdInvalid(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://localhost:5000/api/product/prouct-1", nil)
	req.SetPathValue("id", "")
	s := http_test.NewProductServiceTestInstance(models.Product{
		Name:  "produk",
		Price: 2000,
		Stock: 2,
	}, services.ErrProductInvalid)
	HttpService := ProductServiceAPI{Service: s}
	rec := httptest.NewRecorder()
	HttpService.GetProductById(rec, req)
	response := rec.Result()
	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected error invalid payload, but got: %v", response.StatusCode)
	}
	body, _ := io.ReadAll(response.Body)
	var bodyString map[string]string
	json.Unmarshal(body, &bodyString)
	if _, ok := bodyString["message"]; !ok {
		t.Fatalf("Expected response body to have any related error invalid request message")
	}
}
func TestGetProductByIdNotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://localhost:5000/api/product/prouct-100", nil)
	req.SetPathValue("id", "prouct-100")
	s := http_test.NewProductServiceTestInstance(models.Product{}, services.ErrProductNotFound)
	HttpService := ProductServiceAPI{Service: s}
	rec := httptest.NewRecorder()
	HttpService.GetProductById(rec, req)
	response := rec.Result()
	if response.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected error not found, but got: %v", response.StatusCode)
	}
	body, _ := io.ReadAll(response.Body)
	var bodyString map[string]string
	json.Unmarshal(body, &bodyString)
	if _, ok := bodyString["message"]; !ok {
		t.Fatalf("Expected response body to have any related error not found message")
	}
}
