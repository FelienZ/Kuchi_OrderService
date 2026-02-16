package middleware

import (
	"fmt"
	"net/http"
	"time"
)

type LogMiddleware struct {
	Next http.Handler
}

func (m *LogMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	lrw := &LogResponseWriter{
		ResponseWriter: w,
		statusCode:     200,
	}
	defer func() {
		fmt.Printf("[path]: %s, [duration]: %v, [method]: %s, [statusCode]: %d \n", r.RequestURI, time.Since(start), r.Method, lrw.statusCode)
	}() //deffensive log
	m.Next.ServeHTTP(lrw, r) // (next, recovery)
}
