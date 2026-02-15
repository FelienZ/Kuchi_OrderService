package middleware

import (
	"encoding/json"
	"net/http"
)

type RecoveryMiddleware struct {
	Next http.Handler
}

func (m *RecoveryMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Internal Server Error",
			})
		}
	}()
	m.Next.ServeHTTP(w, r)
}
