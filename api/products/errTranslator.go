package products

import (
	"encoding/json"
	"go-inventory/services"
	"net/http"
)

func errResponseHelper(err error, w http.ResponseWriter) {
	switch err {
	case services.ErrProductConflict:
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{
			"message": err.Error(),
		})
	case services.ErrProductInvalid:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": err.Error(),
		})
	case services.ErrProductNotEnough:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": err.Error(),
		})
	case services.ErrProductNotFound:
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
