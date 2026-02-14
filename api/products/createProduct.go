package products

import (
	"encoding/json"
	"go-inventory/models"
	"net/http"
)

func (s *ProductServiceAPI) CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newProduct models.Product
	if errEncode := json.NewDecoder(r.Body).Decode(&newProduct); errEncode != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid Create Product Request",
		})
		return
	}
	// panggil create, nanti create panggil save + method baru untuk ovw
	if errCreate := s.Service.Create(newProduct); errCreate != nil {
		errResponseHelper(errCreate, w)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Success Created",
	})
}
