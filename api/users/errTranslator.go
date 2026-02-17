package users

import (
	"encoding/json"
	"go-inventory/services"
	"net/http"
)

func errResponseHelper(err error, w http.ResponseWriter) {
	switch err {
	case services.ErrUserConflict:
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{
			"message": err.Error(),
		})
	case services.ErrUserInvalid:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": err.Error(),
		})
	case services.ErrUserNotFound:
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"message": err.Error(),
		})
	default:
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Failed to Get User Response, Internal Server Error",
		})
	}
}
