package middleware

import (
	"log"
	"net/http"
	"time"
)

type responseWriterWrapper struct {
	http.ResponseWriter
	StatusCode int
}

func (r *responseWriterWrapper) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		responseWritter := &responseWriterWrapper{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}

		next.ServeHTTP(responseWritter, r)

		log.Printf(
			"| %d | %s | %s | %s | %s |",
			responseWritter.StatusCode,
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			time.Since(startTime),
		)
	})
}
