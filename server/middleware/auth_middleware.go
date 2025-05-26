package middleware

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"real-time-forum/config"
	"real-time-forum/utils"
	// "strings"
)

// AuthMiddleware checks if the request has a valid JWT
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	fmt.Println("AuthMiddleware")
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Format should be "Bearer {token}"
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		tokenString := authHeader[7:]
		token, err := ValidateJWT(tokenString)
		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Token is valid, proceed to the next handler
		next(w, r)
	}
}

//func AuthMiddlewareWithRedirectFunc(next http.HandlerFunc, redirectPath string) http.HandlerFunc {
//	// Convert the HandlerFunc to a Handler by using http.HandlerFunc adapter
//	return AuthMiddlewareWithRedirect(http.HandlerFunc(next), redirectPath)
//}

//func AuthMiddlewareWithRedirect(next http.Handler, redirectPath string) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		// Try to get token from cookies
//		cookie, err := r.Cookie("jwt_token")
//		tokenString := ""
//
//		if err == nil && cookie.Value != "" {
//			tokenString = cookie.Value
//		} else {
//			// No token in cookie, check Authorization header
//			authHeader := r.Header.Get("Authorization")
//			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
//				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
//			} else {
//				// No token found, redirect to login page
//				http.Redirect(w, r, redirectPath, http.StatusFound)
//				return
//			}
//		}
//
//		// Use your existing function to validate token and extract user ID
//		userID, err := utils.ExtractUserIDFromToken(tokenString)
//		if err != nil {
//			// Invalid token, redirect to login page
//			http.Redirect(w, r, redirectPath, http.StatusFound)
//			return
//		}
//
//		// Check for blacklisted tokens
//		if utils.IsTokenBlacklisted(tokenString) {
//			// Token is blacklisted, redirect to login
//			http.Redirect(w, r, redirectPath, http.StatusFound)
//			return
//		}
//
//		// Token is valid, add user ID to request context
//		ctx := context.WithValue(r.Context(), "user_id", userID)
//		ctx = context.WithValue(ctx, "token", tokenString)
//
//		// Serve the protected content
//		next.ServeHTTP(w, r.WithContext(ctx))
//	}
//}

// ValidateJWT verifies a JWT token
func ValidateJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return config.GetJWTSecret(), nil
	})

	return token, err
}

// ContextKey is a type for context keys used in middleware
// This prevents collisions with other context keys
// You may want to move this to a shared package if used elsewhere

type ContextKey string

const (
	ContextUserID ContextKey = "user_id"
	ContextToken  ContextKey = "token"
)

// WebSocketAuthMiddleware checks JWT, extracts user info, and adds it to context for WebSocket handlers
func WebSocketAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	fmt.Println("WebSocketAuthMiddleware")
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract token from query parameter
		tokenString := r.URL.Query().Get("token")
		if tokenString == "" {
			http.Error(w, "Unauthorized: No token provided in query", http.StatusUnauthorized)
			return
		}

		fmt.Println("WebSocketAuthMiddleware: token from query =", tokenString)

		// Validate token and extract user ID
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil || claims == nil || claims.ID == 0 {
			fmt.Println("WebSocketAuthMiddleware: token validation failed:", err)
			http.Error(w, "Unauthorized: Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Add user ID and token to context
		ctx := context.WithValue(r.Context(), ContextUserID, claims.ID)
		ctx = context.WithValue(ctx, ContextToken, tokenString)

		// Call the next handler with the new context
		next(w, r.WithContext(ctx))
	}
}
