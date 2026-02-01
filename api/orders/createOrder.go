package orders

import (
	"encoding/json"
	"go-inventory/models"
	"net/http"
)

func (s *OrderAPIServices) CreateOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid Create Order Request",
			"data":    err.Error(),
		})
		return
	}
	if errOrder := s.Service.CreateOrder(order); errOrder != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Fail to Create Order",
			"data":    errOrder.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Success Create Order",
	})
}
