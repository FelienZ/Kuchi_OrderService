package users

import (
	"encoding/json"
	exceptions "go-inventory/api/Exceptions"
	"go-inventory/internal/helper"
	"go-inventory/models"
	"net/http"
)

func (s *UserAPIServices) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(helper.UserDataKey).(models.UserIdentity).ID
	user, err := s.Service.GetByID(r.Context(), userID)
	if err != nil {
		exceptions.ErrorHandlerTranslator(err, w)
		return
	}
	json.NewEncoder(w).Encode(models.APIResponse[models.User]{
		Message: "Success Get User Data",
		Data:    user,
	})
}
