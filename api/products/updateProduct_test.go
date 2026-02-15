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

func TestUpdateProduct(t *testing.T) {
	newProduct := models.Product{
		Name:  "new-brand",
		Price: 20000,
		Stock: 2,
	}
	bodyRequest, _ := json.Marshal(newProduct)
	req := httptest.NewRequest(http.MethodPut, "http://localhost:5000/api/products/product-1", bytes.NewBuffer(bodyRequest))
	req.SetPathValue("id", "product-1")
	s := http_test.NewProductServiceTestInstance(newProduct, nil)
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
	if _, ok := bodyString["message"]; !ok {
		t.Fatalf("Expected response body to have any success message")
	}
}

func TestUpdateProductInvalid(t *testing.T) {
	newProduct := models.Product{
		Name: "new-brand",
		// Price: 20000,
		Stock: 2,
	}
	bodyRequest, _ := json.Marshal(newProduct)
	req := httptest.NewRequest(http.MethodPut, "http://localhost:5000/api/products/product-1", bytes.NewBuffer(bodyRequest))
	req.SetPathValue("id", "")
	s := http_test.NewProductServiceTestInstance(newProduct, services.ErrProductInvalid)
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
	if bodyString["message"] != services.ErrProductInvalid.Error() {
		t.Fatalf("Expected response message to equal %s, but got %s", services.ErrProductInvalid.Error(), bodyString["message"])
	}
}

func TestUpdateProductNotFound(t *testing.T) {
	newProduct := models.Product{
		Name: "new-brand",
		// Price: 20000,
		Stock: 2,
	}
	bodyRequest, _ := json.Marshal(newProduct)
	req := httptest.NewRequest(http.MethodPut, "http://localhost:5000/api/products/product-1", bytes.NewBuffer(bodyRequest))
	req.SetPathValue("id", "product-x")
	s := http_test.NewProductServiceTestInstance(newProduct, services.ErrProductNotFound)
	Httpservice := ProductServiceAPI{Service: s}
	rec := httptest.NewRecorder()
	Httpservice.UpdateProduct(rec, req)

	response := rec.Result()
	body, _ := io.ReadAll(response.Body)

	if response.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected status 404, but got: %v", response.StatusCode)
	}
	var bodyString map[string]string
	json.Unmarshal(body, &bodyString)
	if bodyString["message"] != services.ErrProductNotFound.Error() {
		t.Fatalf("Expected response message to equal %s, but got %s", services.ErrProductNotFound.Error(), bodyString["message"])
	}
}
