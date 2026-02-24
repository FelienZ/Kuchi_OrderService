package middleware

import (
	"context"
	"encoding/json"
	"go-inventory/internal/helper"
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
		userData, errSession := m.SessionService.ValidateSession(cookie.Value)
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
		userIdentity := models.UserIdentity{
			ID:    userData.ID,
			Role:  userData.Role,
			Email: userData.Email,
		}
		ctx := context.WithValue(r.Context(), helper.UserDataKey, userIdentity)
		Next.ServeHTTP(w, r.WithContext(ctx))
	})
}
