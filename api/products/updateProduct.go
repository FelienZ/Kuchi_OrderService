package products

import (
	"encoding/json"
	exceptions "go-inventory/api/Exceptions"
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
		})
		return
	}
	if errUpdate := s.Service.UpdateProductData(r.Context(), id, newProduct); errUpdate != nil {
		exceptions.ErrorHandlerTranslator(errUpdate, w)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "success update product data",
	})
}
