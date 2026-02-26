package orders

import (
	"encoding/json"
	"go-inventory/models"
	"net/http"
	"sort"
	"strconv"
)

type OrderAPIServices struct {
	Service models.OrderService
}

var stringToStatus = map[string]models.Status{
	"PAID":     models.PAID,
	"CANCELED": models.CANCELED,
	"PENDING":  models.PENDING,
}

func (s *OrderAPIServices) GetOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	status := r.URL.Query().Get("status")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	orders := s.Service.List(models.GetOrderParameter{Status: stringToStatus[status], Limit: limit, Offset: offset})
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].CreatedAt.Before(orders[j].CreatedAt)
	})
	json.NewEncoder(w).Encode(orders)
}
