package orders

import (
	"encoding/json"
	"go-inventory/services"
	"net/http"
	"sort"
)

type OrderAPIServices struct {
	Service *services.OrderServiceImpl
}

func (s *OrderAPIServices) GetOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	orders := s.Service.List()
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].CreatedAt.Before(orders[j].CreatedAt)
	})
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}
