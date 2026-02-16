package orders

import (
	"encoding/json"
	"net/http"
)

func (s *OrderAPIServices) CancelOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.Service.CancelOrder(id); err != nil {
		errResponseHelper(err, w)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "success cancel order",
	})
}
