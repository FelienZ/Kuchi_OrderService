package middleware

import (
	"encoding/json"
	"go-inventory/internal/helper"
	"go-inventory/models"
	"net/http"
	"slices"
)

func RoleMiddleware(allowedRole ...models.Role) func(http.Handler) http.Handler {
	return func(Next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userIdentity, ok := r.Context().Value(helper.UserDataKey).(models.UserIdentity)
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{
					"message": "Invalid Authorization",
				})
				return
			}
			if slices.Contains(allowedRole, userIdentity.Role) {
				Next.ServeHTTP(w, r)
				return
			}
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Permission is Denied for current Role",
			})
		})
	}
}
