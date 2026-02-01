package orders

import (
	"encoding/json"
	"go-inventory/services"
	"net/http"
)

type OrderAPIServices struct {
	Service *services.OrderServiceImpl
}

func (s *OrderAPIServices) GetOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	orders := s.Service.List()
	json.NewEncoder(w).Encode(orders)
}
