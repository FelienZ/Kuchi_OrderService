package products

import (
	"encoding/json"
	"go-inventory/models"
	"net/http"
)

func (s *ProductServiceAPI) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var newProduct models.UpdateProductRequest
	if errDecode := json.NewDecoder(r.Body).Decode(&newProduct); errDecode != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "unexpected error decode update product request",
			"data":    errDecode.Error(),
		})
		return
	}
	if errUpdate := s.Service.UpdateProductData(id, newProduct); errUpdate != nil {
		StatusCodeHelper(errUpdate, w)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "failed to update product data",
			"data":    errUpdate.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "success update product data",
	})
}
