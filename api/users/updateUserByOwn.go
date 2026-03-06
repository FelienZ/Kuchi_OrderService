package users

import (
	"encoding/json"
	exceptions "go-inventory/api/Exceptions"
	"go-inventory/internal/helper"
	"go-inventory/models"
	"net/http"
)

func (s *UserAPIServices) UpdateUserByOwn(w http.ResponseWriter, r *http.Request) {
	ctx, ok := r.Context().Value(helper.UserDataKey).(models.UserIdentity)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid Authorization",
		})
		return
	}
	var newUser models.UserUpdateRequest
	if errDec := json.NewDecoder(r.Body).Decode(&newUser); errDec != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "invalid action to decode update User payload",
		})
		return
	}
	if errUpdate := s.Service.UpdateUser(r.Context(), ctx.ID, newUser); errUpdate != nil {
		exceptions.ErrorHandlerTranslator(errUpdate, w)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "success update user data",
	})
}
