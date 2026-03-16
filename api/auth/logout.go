package auth

import (
	"encoding/json"
	"net/http"
)

func (s *AuthAPIServices) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// validasi cookie di middleware auth
	_, errCookie := r.Cookie("session_id")
	if errCookie != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "error No sessionID",
		})
		return
	}
	// ini delete di repo
	/* errLogout := s.Services.Logout(r.Context(), cookieData.Value)
	if errLogout != nil {
		exceptions.ErrorHandlerTranslator(errLogout, w)
		return
	} */
	//reassign
	newCookie := new(http.Cookie)
	newCookie = &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
	}
	http.SetCookie(w, newCookie)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "success logout",
	})
}
