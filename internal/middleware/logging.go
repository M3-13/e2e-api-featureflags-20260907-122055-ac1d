package middleware

import "net/http"

// Logging wraps next. The real implementation is a later feature ticket; this
// pass-through keeps every route behind the middleware from the start.
func Logging(next http.Handler) http.Handler {
	return next
}
