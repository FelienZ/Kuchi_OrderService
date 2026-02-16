package orders

import (
	"encoding/json"
	"net/http"
)

func (s *OrderAPIServices) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.Service.DeleteOrder(id); err != nil {
		errResponseHelper(err, w)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "success delete order",
	})
}
