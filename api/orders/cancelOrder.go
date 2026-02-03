package orders

import (
	"encoding/json"
	"net/http"
)

func (s *OrderAPIServices) CancelOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.Service.CancelOrder(id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "failed to cancel order",
			"data":    err.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "success cancel order",
	})
}
