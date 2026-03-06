package orders

import (
	"encoding/json"
	exceptions "go-inventory/api/Exceptions"
	"net/http"
)

func (s *OrderAPIServices) CancelOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.Service.CancelOrder(r.Context(), id); err != nil {
		exceptions.ErrorHandlerTranslator(err, w)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "success cancel order",
	})
}
