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
	// statusVal := stringToStatus[status]
	var statusPtr *models.Status
	if status != "" {
		statusVal, ok := stringToStatus[status]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "unknown status type parameter",
			})
			return
		}
		statusPtr = &statusVal
	}
	orders, _ := s.Service.List(r.Context(), models.GetOrderParameter{Status: statusPtr, Limit: limit, Offset: offset})
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].CreatedAt.Before(orders[j].CreatedAt)
	})
	json.NewEncoder(w).Encode(orders)
}
