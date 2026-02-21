package middleware

import (
	"context"
	"encoding/json"
	"go-inventory/models"
	"go-inventory/services"
	"net/http"
)

type AuthMiddleware struct {
	SessionService models.UserSessionService
}

func (m *AuthMiddleware) Wrap(Next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// validator session
		cookie, errCookie := r.Cookie("session_id")
		if errCookie != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Invalid Authorization",
			})
			return
		}
		userID, errSession := m.SessionService.ValidateSession(cookie.Value)
		if errSession == services.ErrSessionInvalidCredentials {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "session is invalid",
			})
			// reassign
			newCookie := new(http.Cookie)
			newCookie = &http.Cookie{
				Name:     "session_id",
				Value:    "",
				MaxAge:   -1,
				Path:     "/",
				HttpOnly: true,
			}
			http.SetCookie(w, newCookie)
			return
		}
		if errSession != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "session is invalid",
			})
			return
		}
		ctx := context.WithValue(r.Context(), "userID", userID)
		Next.ServeHTTP(w, r.WithContext(ctx))
	})
}
