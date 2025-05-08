//package utils
//
//import (
//	"errors"
//	"fmt"
//	"net/http"
//	"real-time-forum/config"
//	"real-time-forum/models"
//	"strings"
//	"time"
//
//	"github.com/golang-jwt/jwt/v5"
//)
//
//// Claims defines the structure for JWT claims
//type Claims struct {
//	Email string `json:"email"`
//	ID    int    `json:"id"`
//	jwt.RegisteredClaims
//}
//
//// GenerateJWT creates a new JWT token for a user
//func GenerateJWT(user models.User) (string, error) {
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
//		"id":    user.ID,
//		"email": user.Email,
//		"exp":   time.Now().Add(time.Hour * 24).Unix(),
//	})
//
//	tokenString, err := token.SignedString(config.GetJWTSecret())
//	if err != nil {
//		return "", err
//	}
//
//	return tokenString, nil
//}
//
//// ValidateJWT validates a JWT token and returns the claims if valid
//func ValidateJWT(tokenString string) (*Claims, error) {
//	// Parse the token
//	token, err := jwt.ParseWithClaims(
//		tokenString,
//		&Claims{},
//		func(token *jwt.Token) (interface{}, error) {
//			// Validate the algorithm
//			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
//			}
//			return []byte(config.GetJWTSecret()), nil
//			//return []byte(jwtSecret), nil
//		},
//	)
//
//	//if err != nil {
//	//	return nil, err
//	//}
//	if err != nil {
//		return nil, fmt.Errorf("failed to parse token: %w", err)
//	}
//
//	// Check if the token is valid
//	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
//		// Check if the token is expired
//		if time.Now().Unix() > claims.ExpiresAt.Unix() {
//			return nil, errors.New("token expired")
//		}
//		return claims, nil
//	}
//
//	return nil, errors.New("invalid token")
//}
//
////
////// ExtractUserIDFromRequest extracts the user ID from the request using the JWT token
////func ExtractUserIDFromRequest(r *http.Request) (int, error) {
////	// Get the Authorization header
////	authHeader := r.Header.Get("Authorization")
////	if authHeader == "" {
////		return 0, fmt.Errorf("authorization header is required")
////	}
////
////	// Split the header to get the token
////	splitToken := strings.Split(authHeader, "Bearer ")
////	if len(splitToken) != 2 {
////		return 0, fmt.Errorf("invalid token format")
////	}
////
////	tokenString := splitToken[1]
////
////	fmt.Println("TOKEN", tokenString)
////
////	// Validate the token using the existing utils function
////	claims, err := ValidateJWT(tokenString)
////	if err != nil {
////		return 0, err
////	}
////	fmt.Println("Claims", claims.ID)
////	return claims.ID, nil
////}
//
//// ExtractUserIDFromRequest extracts the user ID from the request using the JWT token
////func ExtractUserIDFromRequest(r *http.Request) (int, error) {
////	// Get the Authorization header
////	authHeader := r.Header.Get("Authorization")
////	fmt.Println(">>>>>>>>>>>>...", authHeader)
////	if authHeader == "" {
////		return 0, fmt.Errorf("authorization header is required")
////	}
////
////	// Split the header to get the token
////	parts := strings.SplitN(authHeader, " ", 2)
////	if len(parts) != 2 || parts[0] != "Bearer" {
////		return 0, fmt.Errorf("invalid token format, expected 'Bearer <token>'")
////	}
////
////	tokenString := parts[1]
////
////	// Validate the token using the existing utils function
////	claims, err := ValidateJWT(tokenString)
////	if err != nil {
////		return 0, err
////	}
////
////	return claims.ID, nil
////}
//
//func ExtractUserIDFromRequest(r *http.Request) (int, error) {
//	// Try to get the token from the Authorization header first
//	authHeader := r.Header.Get("Authorization")
//	if authHeader != "" {
//		fmt.Println("Authorization header found:", authHeader)
//		parts := strings.SplitN(authHeader, " ", 2)
//		if len(parts) == 2 && parts[0] == "Bearer" {
//			tokenString := parts[1]
//			claims, err := ValidateJWT(tokenString)
//			if err != nil {
//				return 0, err
//			}
//			return claims.ID, nil
//		}
//	}
//
//	// If not in header, check the query parameter (for WebSocket)
//	queryToken := r.URL.Query().Get("token")
//	if queryToken != "" {
//		fmt.Println("Token found in query parameter:", queryToken)
//		claims, err := ValidateJWT(queryToken)
//		if err != nil {
//			return 0, err
//		}
//		return claims.ID, nil
//	}
//
//	// If no token was found at all
//	return 0, fmt.Errorf("no valid token provided")
//}

package utils

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"real-time-forum/config"
	"real-time-forum/models"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims defines the structure for JWT claims
type Claims struct {
	Email string `json:"email"`
	ID    int    `json:"id"`
	jwt.RegisteredClaims
}

// TokenBlacklist stores invalidated tokens
var (
	TokenBlacklist     = make(map[string]time.Time)
	TokenBlacklistLock sync.RWMutex
)

// AddTokenToBlacklist adds a token to the blacklist with expiry time
func AddTokenToBlacklist(tokenString string, expiry time.Time) {
	TokenBlacklistLock.Lock()
	defer TokenBlacklistLock.Unlock()
	TokenBlacklist[tokenString] = expiry
}

// IsTokenBlacklisted checks if a token is in the blacklist
func IsTokenBlacklisted(tokenString string) bool {
	TokenBlacklistLock.RLock()
	defer TokenBlacklistLock.RUnlock()
	_, found := TokenBlacklist[tokenString]
	return found
}

// CleanupBlacklist removes expired tokens from the blacklist
func CleanupBlacklist() {
	TokenBlacklistLock.Lock()
	defer TokenBlacklistLock.Unlock()
	now := time.Now()
	for token, expiry := range TokenBlacklist {
		if now.After(expiry) {
			delete(TokenBlacklist, token)
		}
	}
}

// GenerateJWT creates a new JWT token for a user
func GenerateJWT(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(config.GetJWTSecret())
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT validates a JWT token and returns the claims if valid
func ValidateJWT(tokenString string) (*Claims, error) {
	// Check if token is blacklisted
	if IsTokenBlacklisted(tokenString) {
		return nil, errors.New("token revoked")
	}

	// Parse the token
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			// Validate the algorithm
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(config.GetJWTSecret()), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Check if the token is valid
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// Check if the token is expired
		if time.Now().Unix() > claims.ExpiresAt.Unix() {
			return nil, errors.New("token expired")
		}
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func ExtractUserIDFromRequest(r *http.Request) (int, error) {
	var tokenString string

	// Try to get the token from the Authorization header first
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		fmt.Println("Authorization header found:", authHeader)
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			tokenString = parts[1]
		}
	}

	// If not in header, check the query parameter (for WebSocket)
	if tokenString == "" {
		tokenString = r.URL.Query().Get("token")
		if tokenString != "" {
			fmt.Println("Token found in query parameter:", tokenString)
		}
	}

	// If no token was found at all
	if tokenString == "" {
		return 0, fmt.Errorf("no valid token provided")
	}

	// Store the token in the request context for potential blacklisting
	ctx := r.Context()
	ctx = context.WithValue(ctx, "token", tokenString)
	*r = *r.WithContext(ctx)

	claims, err := ValidateJWT(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.ID, nil
}

func ExtractUserIDFromToken(tokenString string) (int, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the algorithm is what you expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return config.GetJWTSecret(), nil
	})

	if err != nil {
		return 0, err
	}

	// Check if token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extract user ID from claims
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			return 0, fmt.Errorf("invalid user ID in token")
		}

		userID := int(userIDFloat)
		return userID, nil
	}

	return 0, fmt.Errorf("invalid token")
}

// BlacklistUserTokens blacklists all tokens for a specific user ID
func BlacklistUserTokens(userID int) {
	// This is a simplified version - in a real implementation,
	// you would store user tokens in a database and invalidate them all
	// For this example, we're just marking that this user's sessions should be invalidated
	// The actual implementation would depend on your token storage strategy

	// Create a special entry in the blacklist for this user ID
	TokenBlacklistLock.Lock()
	defer TokenBlacklistLock.Unlock()

	// Use a special format to indicate this is a user ID blacklist entry
	TokenBlacklist[fmt.Sprintf("user:%d", userID)] = time.Now().Add(time.Hour * 24)
}

// IsUserBlacklisted checks if a user has been blacklisted
func IsUserBlacklisted(userID int) bool {
	TokenBlacklistLock.RLock()
	defer TokenBlacklistLock.RUnlock()
	_, found := TokenBlacklist[fmt.Sprintf("user:%d", userID)]
	return found
}

// SplitPath splits a URL path into segments
func SplitPath(path string) []string {
	// Remove leading slash if present
	if path[0] == '/' {
		path = path[1:]
	}
	// Remove trailing slash if present
	if len(path) > 0 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	return strings.Split(path, "/")
}
