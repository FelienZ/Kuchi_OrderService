package orders

import (
	"encoding/json"
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
	userID, ok := r.Context().Value(helper.UserIDKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid Authorization",
		})
		return
	}
	newOrder := models.Order{
		Item:   orderRequest.Item,
		UserID: userID,
	}
	if errOrder := s.Service.CreateOrder(newOrder); errOrder != nil {
		errResponseHelper(errOrder, w)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Success Create Order",
	})
}
