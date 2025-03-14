package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"real-time-forum/models"
)

// ErrorHandler is a middleware that recovers from panics and logs the error.
func ErrorHandler(next http.HandlerFunc, errorLogger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error with additional context
				//log.Printf("Internal server error: %v, Request: %s %s", err, r.Method, r.URL.Path)
				errorLogger.Printf("Internal server error: %v, Request: %s %s", err, r.Method, r.URL.Path)

				// Create a custom error response
				response := models.ErrorResponse{
					Code:    http.StatusInternalServerError,
					Message: "An unexpected error occurred. Please try again later.",
				}

				fmt.Println(response)

				// Set the response header and write the JSON response
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				//json.NewEncoder(w).Encode(response)
				if err := json.NewEncoder(w).Encode(response); err != nil {
					errorLogger.Printf("Error encoding JSON response: %v", err)
				}
			}
		}()
		next(w, r)
	}
}
