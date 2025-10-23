package middleware

import (
	"log"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	w          http.ResponseWriter
	statusCode int
}

func Logger(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		log.Printf("Received request: %s %s", r.Method, r.URL.Path, time.Since(start))

		next.ServeHTTP(w, r)
	})

}
