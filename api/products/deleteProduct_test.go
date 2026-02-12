package products

import (
	"encoding/json"
	"go-inventory/internal/testutils"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteProduct(t *testing.T) {
	// bodyRequest, _ := json.Marshal(models.Product{Name: "barang", Price: 20000, Stock: 3})
	req := httptest.NewRequest(http.MethodDelete, "http://localhost:5000/api/products/product-1", nil)
	// tidak kenal pathValue karena bukan serveMux, mirip polosan http request. path value punya servemux, di sini set manual:
	req.SetPathValue("id", "product-1") // key + value

	s := testutils.NewProductServiceTestInstance()
	HttpService := ProductServiceAPI{Service: s}

	rec := httptest.NewRecorder()
	HttpService.DeleteProduct(rec, req)
	response := rec.Result()

	body, _ := io.ReadAll(response.Body)
	if response.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, but got: %v", response.StatusCode)
	}
	var bodyString map[string]string
	json.Unmarshal(body, &bodyString)
	if string(bodyString["message"]) != "success delete product" {
		t.Fatalf("expected response message after delete to equal %s but got %s", "success delete product", bodyString)
	}
}

func TestDeleteInvalid(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "http://localhost:5000/api/products/product-1", nil)
	// req.SetPathValue("id", "") // bahkan skip

	rec := httptest.NewRecorder()
	s := testutils.NewProductServiceTestInstance()
	service := ProductServiceAPI{Service: s}

	service.DeleteProduct(rec, req)

	response := rec.Result()
	body, _ := io.ReadAll(response.Body)

	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected statusCode to equal 400, but got: %v", response.StatusCode)
	}
	var bodyString map[string]string
	json.Unmarshal(body, &bodyString)
	if bodyString["message"] != "failed to delete product" {
		t.Fatalf("Expected response message to equal: %s but got: %s", "failed to delete product", bodyString["message"])
	}
}
