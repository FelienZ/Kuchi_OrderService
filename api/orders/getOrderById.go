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
		errResponseHelper(err, w)
		return
	}
	json.NewEncoder(w).Encode(map[string]models.Order{
		"data": order,
	})
}
