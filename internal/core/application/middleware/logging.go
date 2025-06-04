package middleware

import (
	"bytes"
	"log"
	"net/http"
)

type LoggingMiddleware struct {
}

func NewLoggingMiddleware() *LoggingMiddleware {
	return &LoggingMiddleware{}
}

type statusRecorder struct {
	http.ResponseWriter
	status   int
	response *bytes.Buffer
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	r.response.Write(b)
	return r.ResponseWriter.Write(b)
}

func (lm *LoggingMiddleware) Log(next http.Handler) http.Handler {
	exemptedPaths := []string{"/metrics"}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request path is in the exemptedPaths list
		for _, path := range exemptedPaths {
			if r.URL.Path == path {
				next.ServeHTTP(w, r) // Call the next handler without logging
				return
			}
		}

		responseBuffer := &bytes.Buffer{}
		recorder := &statusRecorder{
			ResponseWriter: w,
			status:         http.StatusOK,
			response:       responseBuffer,
		}

		defer func() {
			log.Printf(
				"%s %s %d \nResponse Body=%s",
				r.Method,
				r.URL.Path,
				recorder.status,
				recorder.response.String(),
			)
		}()

		next.ServeHTTP(recorder, r)
	})
}
