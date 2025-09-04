package middleware

import (
	"net/http"
)

// add common headers to each response
func HeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// set common headers
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Applications", "API")
		w.Header().Set("Version", "1.0")

		// call the next handler
		next.ServeHTTP(w, r)
	})
}
