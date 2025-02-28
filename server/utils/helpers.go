package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"real-time-forum/config"
	"real-time-forum/models"
	"time"
)

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
