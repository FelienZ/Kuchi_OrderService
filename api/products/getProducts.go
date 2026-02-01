package products

import (
	"encoding/json"
	"go-inventory/services"
	"net/http"
)

type ProductServiceAPI struct {
	Service *services.ProductServiceImpl
}

func (s *ProductServiceAPI) GetProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := s.Service.List()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
