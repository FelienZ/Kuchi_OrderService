package products

import (
	"encoding/json"
	"net/http"
)

func (s *ProductServiceAPI) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.Service.DeleteProduct(id); err != nil {
		errResponseHelper(err, w)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "success delete product",
	})
}
