package orders

import (
	"encoding/json"
	"go-inventory/services"
	"net/http"
)

func errResponseHelper(err error, w http.ResponseWriter) {
	switch err {
	case services.ErrOrderConflict:
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{
			"message": err.Error(),
		})
	case services.ErrOrderInvalid:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": err.Error(),
		})
	case services.ErrOrderNotFound:
		w.WriteHeader(http.StatusNotFound)
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
