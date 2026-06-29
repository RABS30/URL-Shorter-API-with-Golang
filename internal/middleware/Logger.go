package middleware

import (
	"log"
	"net/http"
	"time"
)

const ErrorLogKey contextKey = "errorDetails"

type ResponseWriterWrapper struct {
	http.ResponseWriter
	StatusCode int
	ErrorLog   string
}

func (r *ResponseWriterWrapper) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *ResponseWriterWrapper) WriteError(err string) {
	r.ErrorLog = err
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		wrapper := &ResponseWriterWrapper{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
			ErrorLog:       "-",
		}

		next.ServeHTTP(wrapper, r)

		log.Printf(
			"| %d | %s | %s | %s | %s | %s |",
			wrapper.StatusCode,
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			time.Since(startTime),
			wrapper.ErrorLog,
		)
	})
}
