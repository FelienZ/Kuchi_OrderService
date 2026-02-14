package orders

import (
	"go-inventory/services"
	"net/http"
)

func StatusCodeHelper(err error, w http.ResponseWriter) {
	switch err {
	case services.ErrOrderConflict:
		w.WriteHeader(http.StatusConflict)
	case services.ErrOrderInvalid:
		w.WriteHeader(http.StatusBadRequest)
	case services.ErrOrderNotFound:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}
