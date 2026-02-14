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

func TestCreateProduct(t *testing.T) {
	newProduct := []byte(`{"Name": "sebuah Product", "price": 2000, "stock": 2}`)
	req := httptest.NewRequest(http.MethodPost, "http://localhost:5000/api/products", bytes.NewBuffer(newProduct))

	s := http_test.NewProductServiceTestInstance()

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
	if bodyString["message"] != "Success Created" {
		t.Fatalf("expected message to equal %s", "Success Created")
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

	s := http_test.NewProductServiceTestInstance()

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
		t.Fatalf("expected message to equal %s", "Invalid Payload for Product Data")
	}

}
