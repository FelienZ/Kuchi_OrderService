package users

import (
	"encoding/json"
	exceptions "go-inventory/api/Exceptions"
	"go-inventory/models"
	"net/http"
)

func (s *UserAPIServices) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")
	var newData models.UserUpdateRequest
	if errEnc := json.NewDecoder(r.Body).Decode(&newData); errEnc != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "failed to decode update User request",
		})
		return
	}
	if err := s.Service.UpdateUser(r.Context(), id, newData); err != nil {
		exceptions.ErrorHandlerTranslator(err, w)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "success update user data",
	})
}
