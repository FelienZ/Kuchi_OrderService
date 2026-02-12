package products

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-inventory/internal/testutils"
	"go-inventory/models"
	"go-inventory/services"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateProduct(t *testing.T) {
	newProduct := []byte(`{"Name": "sebuah Product", "price": 2000, "stock": 2}`)
	req := httptest.NewRequest(http.MethodPost, "http://localhost:5000/api/products", bytes.NewBuffer(newProduct))
	// bikin mock buat service (next, lepas boundary repo)
	productRepo := testutils.NewProductRepoTestInstance()
	loggerRepo := testutils.NewLoggerTestInstance()
	s := services.ProductServiceImpl{ProductRepo: productRepo, LoggerRepo: loggerRepo}

	HttpService := ProductServiceAPI{Service: &s}
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

	fmt.Println(response.StatusCode)
	fmt.Println(string(body))
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
	productRepo := testutils.NewProductRepoTestInstance()
	loggerRepo := testutils.NewLoggerTestInstance()
	s := services.ProductServiceImpl{ProductRepo: productRepo, LoggerRepo: loggerRepo}

	HttpService := ProductServiceAPI{Service: &s}
	rec := httptest.NewRecorder()
	HttpService.CreateProduct(rec, req)
	response := rec.Result()
	body, _ := io.ReadAll(response.Body)

	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected error Invalid Payload, but got: %v", response.StatusCode)
	}

	var bodyString map[string]string
	json.Unmarshal(body, &bodyString)
	if bodyString["message"] != "Failed to Create Product" {
		t.Fatalf("expected message to equal %s", "Failed to Create Product")
	}

	fmt.Println(response.StatusCode)
	fmt.Println(string(body))
}
