package users

import (
	"encoding/json"
	exceptions "go-inventory/api/Exceptions"
	"go-inventory/models"
	"net/http"
)

type UserAPIServices struct {
	Service models.UserService
}

func (s *UserAPIServices) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newUser models.RegisterRequest
	if errDec := json.NewDecoder(r.Body).Decode(&newUser); errDec != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "failed at encode user request",
		})
		return
	}
	if err := s.Service.RegisterUser(r.Context(), newUser); err != nil {
		// fmt.Println("cek error register: ", err)
		exceptions.ErrorHandlerTranslator(err, w)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Register User Success",
	})
}
