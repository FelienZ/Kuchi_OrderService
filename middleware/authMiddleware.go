package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"go-inventory/internal/helper"
	"go-inventory/internal/jwt"
	"go-inventory/models"
	"net/http"
)

type AuthMiddleware struct {
	TokenManager *jwt.JWTToken
	UserService  models.UserService
}

func (m *AuthMiddleware) Wrap(Next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// validator session
		fmt.Println("Masuk auth midd")
		cookie, errCookie := r.Cookie("session_id")
		if errCookie != nil {
			fmt.Println("err get cookie: ", errCookie)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Invalid Authorization",
			})
			return
		}
		sessionData, errData := m.TokenManager.VerifyToken(cookie.Value)
		fmt.Println("Masuk auth midd sessData: ", sessionData)
		if errData != nil {
			fmt.Println("error validasi sess: ", errData)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "session is invalid",
			})
			newCookie := new(http.Cookie)
			newCookie = &http.Cookie{
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
				Value:    "",
			}
			http.SetCookie(w, newCookie)
			return
		}
		user, errUser := m.UserService.GetByID(r.Context(), sessionData.UserID)
		fmt.Println("Masuk auth midd user: ", user)
		if errUser != nil {
			// fmt.Println("error validasi sess: ", errData)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "session is invalid",
			})
			newCookie := new(http.Cookie)
			newCookie = &http.Cookie{
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
				Value:    "",
			}
			http.SetCookie(w, newCookie)
			return
		}
		userIdentity := models.UserIdentity{
			ID:       sessionData.UserID,
			Email:    user.Email,
			Username: user.Username,
			Role:     user.Role,
		}
		fmt.Println("Cek userIDentity: ", userIdentity)
		ctx := context.WithValue(r.Context(), helper.UserDataKey, userIdentity)
		Next.ServeHTTP(w, r.WithContext(ctx))
	})
}
