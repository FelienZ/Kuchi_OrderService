package products

import (
	"encoding/json"
	"go-inventory/models"
	"net/http"
)

func (s *ProductServiceAPI) GetProductById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	product, err := s.Service.GetByID(id)
	if err != nil {
		errResponseHelper(err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]models.Product{
		"data": product,
	})
}
