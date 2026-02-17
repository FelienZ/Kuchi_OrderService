package users

import (
	"encoding/json"
	"go-inventory/models"
	"net/http"
)

type UserAPIServices struct {
	Service models.UserService
}

func (s *UserAPIServices) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newUser models.UserRequest
	if errDec := json.NewDecoder(r.Body).Decode(&newUser); errDec != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "failed at encode user request",
		})
	}
	if err := s.Service.RegisterUser(newUser); err != nil {
		errResponseHelper(err, w)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Register User Success",
	})
}
