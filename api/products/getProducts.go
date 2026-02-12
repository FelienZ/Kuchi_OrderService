package products

import (
	"encoding/json"
	"go-inventory/models"
	"net/http"
	"sort"
)

type ProductServiceAPI struct {
	Service models.ProductService
}

func (s *ProductServiceAPI) GetProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := s.Service.List()
	// ascend by createdAt
	sort.Slice(data, func(i, j int) bool {
		return data[i].CreatedAt.Before(data[j].CreatedAt)
	})
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
