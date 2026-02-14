package products

import (
	"go-inventory/services"
	"net/http"
)

func StatusCodeHelper(err error, w http.ResponseWriter) {
	switch err {
	case services.ErrProductConflict:
		w.WriteHeader(http.StatusConflict)
	case services.ErrProductInvalid:
		w.WriteHeader(http.StatusBadRequest)
	case services.ErrProductNotEnough:
		w.WriteHeader(http.StatusBadRequest)
	case services.ErrProductNotFound:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}
