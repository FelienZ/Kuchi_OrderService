package middleware

import (
	"net/http"
)

type CORSOptions struct {
	Next http.Handler
}

func (m *CORSOptions) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Access-Control-Allow-Origin", "http://192.168.1.11:3000")
	source := r.Header.Get("Origin")
	if source != "" {
		w.Header().Set("Access-Control-Allow-Origin", source)
	} else {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	m.Next.ServeHTTP(w, r)
}
