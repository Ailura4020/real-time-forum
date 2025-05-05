package utils

import (
	"errors"
	"fmt"
	"net/http"
	"real-time-forum/config"
	"real-time-forum/models"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims defines the structure for JWT claims
type Claims struct {
	Email string `json:"email"`
	ID    int    `json:"id"`
	jwt.RegisteredClaims
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
			//return []byte(jwtSecret), nil
		},
	)

	//if err != nil {
	//	return nil, err
	//}
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

//
//// ExtractUserIDFromRequest extracts the user ID from the request using the JWT token
//func ExtractUserIDFromRequest(r *http.Request) (int, error) {
//	// Get the Authorization header
//	authHeader := r.Header.Get("Authorization")
//	if authHeader == "" {
//		return 0, fmt.Errorf("authorization header is required")
//	}
//
//	// Split the header to get the token
//	splitToken := strings.Split(authHeader, "Bearer ")
//	if len(splitToken) != 2 {
//		return 0, fmt.Errorf("invalid token format")
//	}
//
//	tokenString := splitToken[1]
//
//	fmt.Println("TOKEN", tokenString)
//
//	// Validate the token using the existing utils function
//	claims, err := ValidateJWT(tokenString)
//	if err != nil {
//		return 0, err
//	}
//	fmt.Println("Claims", claims.ID)
//	return claims.ID, nil
//}

// ExtractUserIDFromRequest extracts the user ID from the request using the JWT token
//func ExtractUserIDFromRequest(r *http.Request) (int, error) {
//	// Get the Authorization header
//	authHeader := r.Header.Get("Authorization")
//	fmt.Println(">>>>>>>>>>>>...", authHeader)
//	if authHeader == "" {
//		return 0, fmt.Errorf("authorization header is required")
//	}
//
//	// Split the header to get the token
//	parts := strings.SplitN(authHeader, " ", 2)
//	if len(parts) != 2 || parts[0] != "Bearer" {
//		return 0, fmt.Errorf("invalid token format, expected 'Bearer <token>'")
//	}
//
//	tokenString := parts[1]
//
//	// Validate the token using the existing utils function
//	claims, err := ValidateJWT(tokenString)
//	if err != nil {
//		return 0, err
//	}
//
//	return claims.ID, nil
//}

func ExtractUserIDFromRequest(r *http.Request) (int, error) {
	// Try to get the token from the Authorization header first
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		fmt.Println("Authorization header found:", authHeader)
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			tokenString := parts[1]
			claims, err := ValidateJWT(tokenString)
			if err != nil {
				return 0, err
			}
			return claims.ID, nil
		}
	}

	// If not in header, check the query parameter (for WebSocket)
	queryToken := r.URL.Query().Get("token")
	if queryToken != "" {
		fmt.Println("Token found in query parameter:", queryToken)
		claims, err := ValidateJWT(queryToken)
		if err != nil {
			return 0, err
		}
		return claims.ID, nil
	}

	// If no token was found at all
	return 0, fmt.Errorf("no valid token provided")
}
