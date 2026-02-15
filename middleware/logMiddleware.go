package middleware

import (
	"fmt"
	"net/http"
	"time"
)

type LogMiddleware struct {
	Next http.Handler
}

// tugas : log method, duration, statusCode
func (m *LogMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	lrw := &LogResponseWriter{
		ResponseWriter: w,
		statusCode:     200,
	} // penting override writer untuk statusCode, kalau coba write tanpa override -> superfluous (writeHead di handler, tapi middleware pakai lagi atau sebaliknya)
	m.Next.ServeHTTP(lrw, r) // log task diukur sampai selesai response (ini handler asli, next)
	fmt.Printf("path: %s, duration: %v, method: %s, statusCode: %d: \n", r.RequestURI, time.Since(start), r.Method, lrw.statusCode)
}
