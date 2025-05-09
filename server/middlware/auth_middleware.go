package middleware

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"real-time-forum/config"
)

// AuthMiddleware checks if the request has a valid JWT
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
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
