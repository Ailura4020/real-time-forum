package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

// AllowedOrigins is a list of origins that are allowed to access the API
var AllowedOrigins = []string{
	"*",
}

//"http://localhost:8080", // Frontend running on localhost
//"https://example.com",   // Other allowed origins
//"https://another-example.com",

// CORSMiddleware handles CORS requests
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		//fmt.Println("[ORIGIN]", origin)
		fmt.Println("Request Method:", r.Method, " / Request Origin:", origin)

		// Check if the origin is allowed
		if isOriginAllowed(origin) {
			fmt.Println("[ALLOWED]", origin)
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
			//fmt.Println("[ALLOWED ORIGIN WILDCARD]", origin)
			return true // Allow all origins if wildcard is present
		}
		if allowedOrigin == origin {
			//fmt.Println("[ALLOWED ORIGIN ==]", origin)
			return true // Exact match
		}
		// Check for localhost with any port
		if allowedOrigin == "http://localhost:*" {
			//fmt.Println("[ALLOWED ORIGIN localhost]", origin)
			if strings.HasPrefix(origin, "http://localhost:") {
				return true
			}
		}
	}
	return false
}

//package middleware
//
//import "net/http"
//
//// CORSMiddleware handles CORS requests
//func CORSMiddleware(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		w.Header().Set("Access-Control-Allow-Origin", "*")
//		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
//		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
//
//		if r.Method == "OPTIONS" {
//			w.WriteHeader(http.StatusOK)
//			return
//		}
//
//		next.ServeHTTP(w, r)
//	})
//}
