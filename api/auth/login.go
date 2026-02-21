package auth

import (
	"encoding/json"
	"go-inventory/models"
	"net/http"
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
	sessionID, errLogin := s.Services.Login(loginData)
	if errLogin != nil {
		errResponseHelper(errLogin, w)
		return
	}
	cookie := new(http.Cookie)
	cookie = &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		MaxAge:   15 * 60,
		HttpOnly: true, //xss cookie prevent js access
		// SameSite: http.SameSiteStrictMode, // csrf preventive cookie accessed outside site(def: lax -> get bisa)
	}
	http.SetCookie(w, cookie)
	// w.WriteHeader(http.StatusCreated) statok
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Success Login",
	})
}
