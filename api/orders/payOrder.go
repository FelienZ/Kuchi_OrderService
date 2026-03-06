package orders

import (
	"encoding/json"
	exceptions "go-inventory/api/Exceptions"
	"net/http"
)

func (s *OrderAPIServices) PayOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.Service.PayOrder(r.Context(), id); err != nil {
		// fmt.Println("cek err handler payORder: ", err)
		exceptions.ErrorHandlerTranslator(err, w)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "succes pay order",
	})
}
