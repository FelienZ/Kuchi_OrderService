package products

import (
	"encoding/json"
	exceptions "go-inventory/api/Exceptions"
	"go-inventory/models"
	"net/http"
)

func (s *ProductServiceAPI) GetProductById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	product, err := s.Service.GetByID(r.Context(), id)
	if err != nil {
		exceptions.ErrorHandlerTranslator(err, w)
		return
	}
	json.NewEncoder(w).Encode(models.APIResponse[models.Product]{
		Message: "Success Get Product",
		Data:    product,
	})
}
