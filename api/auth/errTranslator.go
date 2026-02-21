package auth

import (
	"encoding/json"
	"go-inventory/services"
	"net/http"
)

func errResponseHelper(err error, w http.ResponseWriter) {
	switch err {
	case services.ErrSessionInvalid:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": err.Error(),
		})
	case services.ErrSessionInvalidCredentials:
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"message": err.Error(),
		})
	default:
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Internal Server Error",
		})
	}
}
