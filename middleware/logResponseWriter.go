package middleware

import (
	"net/http"
)

type LogResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (m *LogResponseWriter) WriteHeader(code int) {
	m.statusCode = code
	m.ResponseWriter.WriteHeader(code)
}
