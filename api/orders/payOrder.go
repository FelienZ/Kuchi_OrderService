package orders

import (
	"encoding/json"
	"net/http"
)

func (s *OrderAPIServices) PayOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.Service.PayOrder(id); err != nil {
		errResponseHelper(err, w)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "succes pay order",
	})
}
