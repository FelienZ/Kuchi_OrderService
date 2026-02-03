package orders

import (
	"encoding/json"
	"net/http"
)

func (s *OrderAPIServices) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := s.Service.DeleteOrder(id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "fail to delete order",
			"data":    err.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "success delete order",
	})
}
