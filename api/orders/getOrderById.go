package orders

import (
	"encoding/json"
	"go-inventory/models"
	"net/http"
)

func (s *OrderAPIServices) GetOrderById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	order, err := s.Service.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "failed to get order data",
			"data":    err.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]models.Order{
		"data": order,
	})
}
