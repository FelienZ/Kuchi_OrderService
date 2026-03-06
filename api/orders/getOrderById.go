package orders

import (
	"encoding/json"
	exceptions "go-inventory/api/Exceptions"
	"go-inventory/models"
	"net/http"
)

func (s *OrderAPIServices) GetOrderById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	order, err := s.Service.GetByID(r.Context(), id)
	if err != nil {
		exceptions.ErrorHandlerTranslator(err, w)
		return
	}
	json.NewEncoder(w).Encode(models.APIResponse[models.Order]{
		Message: "Success Get Order Data",
		Data:    order,
	})
}
