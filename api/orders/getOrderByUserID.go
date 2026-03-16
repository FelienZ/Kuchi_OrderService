package orders

import (
	"encoding/json"
	exceptions "go-inventory/api/Exceptions"
	"go-inventory/internal/helper"
	"go-inventory/models"
	"net/http"
)

func (s *OrderAPIServices) GetOrderByUserID(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(helper.UserDataKey).(models.UserIdentity).ID
	orders, err := s.Service.GetByUserID(r.Context(), userID)
	if err != nil {
		exceptions.ErrorHandlerTranslator(err, w)
		return
	}
	json.NewEncoder(w).Encode(models.APIResponse[[]models.Order]{
		Message: "Success Get Order Data",
		Data:    orders,
	})
}
