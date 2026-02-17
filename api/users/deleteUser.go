package users

import (
	"encoding/json"
	"net/http"
)

func (s *UserAPIServices) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")
	if err := s.Service.DeleteUser(id); err != nil {
		errResponseHelper(err, w)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "success delete user",
	})
}
