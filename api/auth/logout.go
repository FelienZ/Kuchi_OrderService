package auth

import (
	"encoding/json"
	"net/http"
)

func (s *AuthAPIServices) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cookieData, errCookie := r.Cookie("session_id")
	if errCookie != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "error No sessionID",
		})
		return
	}
	// ini delete di repo
	errLogout := s.Services.Logout(cookieData.Value)
	if errLogout != nil {
		errResponseHelper(errLogout, w)
		return
	}
	//reassign
	newCookie := new(http.Cookie)
	newCookie = &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, newCookie)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "success logout",
	})
}
