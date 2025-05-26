package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

// AllowedOrigins is a list of origins that are allowed to access the API
var AllowedOrigins = []string{
	"*", // allow all origins (not for production)
	// Add other allowed origins here
}

// CORSMiddleware handles CORS requests
func CORSMiddleware(next http.Handler) http.Handler {
	fmt.Println("CORSMiddleware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// fmt.Println("Request Method:", r.Method, " / Request Origin:", origin)

		// Check if the origin is allowed
		if isOriginAllowed(origin) {
			// fmt.Println("[ALLOWED]", origin)
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true") // Allow credentials if needed
		}

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// isOriginAllowed checks if the request origin is in the allowed origins list
func isOriginAllowed(origin string) bool {
	for _, allowedOrigin := range AllowedOrigins {
		if allowedOrigin == "*" {
			return true // Allow all origins if wildcard is present
		}
		if allowedOrigin == origin {
			return true // Exact match
		}
		// Check for localhost with any port
		if allowedOrigin == "http://localhost:*" {
			if strings.HasPrefix(origin, "http://localhost:") {
				return true
			}
		}
	}
	return false
}
