package products

import (
	"encoding/json"
	"net/http"
)

func (s *ProductServiceAPI) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.Service.DeleteProduct(id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "failed to delete product",
			"data":    err.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "success delete product",
	})
}
