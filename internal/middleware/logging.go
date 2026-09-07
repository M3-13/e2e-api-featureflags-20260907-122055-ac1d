package middleware

import (
	"log"
	"net/http"
	"time"
)

// statusRecorder wraps an http.ResponseWriter and records the status code that
// was written, defaulting to 200 when the handler never calls WriteHeader.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// Logging wraps next and logs one line per request to stderr (via the standard
// logger) containing exactly: method, path (without query string), status code
// and duration. Request bodies and query parameters (such as user) are never
// logged.
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, rec.status, time.Since(start))
	})
}
