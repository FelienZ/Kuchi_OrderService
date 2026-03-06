package products

import (
	"encoding/json"
	exceptions "go-inventory/api/Exceptions"
	"net/http"
)

func (s *ProductServiceAPI) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.Service.DeleteProduct(r.Context(), id); err != nil {
		exceptions.ErrorHandlerTranslator(err, w)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "success delete product",
	})
}
