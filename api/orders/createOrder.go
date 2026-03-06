package orders

import (
	"encoding/json"
	exceptions "go-inventory/api/Exceptions"
	"go-inventory/internal/helper"
	"go-inventory/models"
	"net/http"
)

func (s *OrderAPIServices) CreateOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var orderRequest models.OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&orderRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid Create Order Request",
		})
		return
	}
	userID, ok := r.Context().Value(helper.UserDataKey).(models.UserIdentity)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid Authorization",
		})
		return
	}
	newOrder := models.OrderPayload{
		Item:   orderRequest.Item,
		UserID: userID.ID,
	}
	orderId, errOrder := s.Service.CreateOrder(r.Context(), newOrder)
	if errOrder != nil {
		exceptions.ErrorHandlerTranslator(errOrder, w)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.APIResponse[string]{
		Message: "Success Create Order",
		Data:    orderId,
	})
}
