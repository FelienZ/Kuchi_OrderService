package auth

import (
	"encoding/json"
	exceptions "go-inventory/api/Exceptions"
	"go-inventory/models"
	"net/http"
	"time"
)

type AuthAPIServices struct {
	Services models.UserSessionService
}

func (s *AuthAPIServices) LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var loginData models.LoginRequest
	errDec := json.NewDecoder(r.Body).Decode(&loginData)
	if errDec != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Failed to encode login request",
		})
		return
	}
	sessionData, errLogin := s.Services.Login(r.Context(), loginData)
	if errLogin != nil {
		// fmt.Println("cek err login: ", errLogin)
		exceptions.ErrorHandlerTranslator(errLogin, w)
		return
	}
	cookie := new(http.Cookie)
	cookie = &http.Cookie{
		Name:     "session_id",
		Value:    sessionData.AccessToken,
		Path:     "/",
		MaxAge:   15 * int(time.Minute),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	// w.WriteHeader(http.StatusCreated) statok
	json.NewEncoder(w).Encode(models.APIResponse[models.LoginResult]{
		Message: "Success Login",
		Data:    sessionData,
	})
}
