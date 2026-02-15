package products

import (
	"bytes"
	"encoding/json"
	"go-inventory/internal/testutils/http_test"
	"go-inventory/models"
	"go-inventory/services"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateProduct(t *testing.T) {
	newProduct := []byte(`{"Name": "sebuah Product", "price": 2000, "stock": 2}`)
	var product models.Product
	json.Unmarshal(newProduct, &product)
	req := httptest.NewRequest(http.MethodPost, "http://localhost:5000/api/products", bytes.NewBuffer(newProduct))

	s := http_test.NewProductServiceTestInstance(product, nil)

	HttpService := ProductServiceAPI{Service: s}
	rec := httptest.NewRecorder()
	HttpService.CreateProduct(rec, req)
	response := rec.Result()
	body, _ := io.ReadAll(response.Body)

	if response.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status created, but got: %v", response.StatusCode)
	}
	var bodyString map[string]string
	json.Unmarshal(body, &bodyString)
	if _, ok := bodyString["message"]; !ok {
		t.Fatalf("Expected response body to have any success message")
	}
}
func TestCreateProductInvalid(t *testing.T) {
	// newProduct := []byte(`{"Name": "sebuah Product", "stock": 2}`)
	product := models.Product{
		Name: "apa ajah",
		// Price: 5000,
		Stock: 10,
	}
	bodyReq, _ := json.Marshal(product) // sama" untuk encode go -> json (atas manual)
	req := httptest.NewRequest(http.MethodPost, "http://localhost:5000/api/products", bytes.NewBuffer(bodyReq))

	s := http_test.NewProductServiceTestInstance(product, services.ErrProductInvalid)

	HttpService := ProductServiceAPI{Service: s}
	rec := httptest.NewRecorder()
	HttpService.CreateProduct(rec, req)
	response := rec.Result()
	body, _ := io.ReadAll(response.Body)
	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected error Invalid Payload, but got: %v", response.StatusCode)
	}

	var bodyString map[string]string
	json.Unmarshal(body, &bodyString)
	if bodyString["message"] != "Invalid Payload for Product Data" {
		t.Fatalf("expected message to equal %s, but got %s", "Invalid Payload for Product Data", bodyString["message"])
	}

}

func TestCreateProductConflict(t *testing.T) {
	product := models.Product{
		Name:  "produk",
		Price: 2000,
		Stock: 3,
	}
	//assume di repo udah ada
	jsonProduct, _ := json.Marshal(product)
	req := httptest.NewRequest(http.MethodPost, "http://localhost:5000/api/products", bytes.NewBuffer(jsonProduct))
	s := http_test.NewProductServiceTestInstance(product, services.ErrProductConflict)
	HttpService := ProductServiceAPI{Service: s}
	rec := httptest.NewRecorder()
	HttpService.CreateProduct(rec, req)
	response := rec.Result()
	body, _ := io.ReadAll(response.Body)
	if response.StatusCode != http.StatusConflict {
		t.Fatalf("Expected error conflict, but Got %v", response.StatusCode)
	}
	var bodyString map[string]string // bentuk go response
	json.Unmarshal(body, &bodyString)
	if bodyString["message"] != "Conflict Action Product" {
		t.Fatalf("Expected response message to equal %s, but got %s", "Conflict Action Product", bodyString["message"])
	}
}
