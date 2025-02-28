package config

import (
	"bufio"
	"os"
	"strings"
)

// JWT secret key (old method, in os ENV's)
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// SetJWTSecret sets the JWT secret from an environment variable
func SetJWTSecret(secret string) {
	jwtSecret = []byte(secret)
}

//
//// GetJWTSecret returns the JWT secret
//func GetJWTSecret() []byte {
//	return jwtSecret
//}

func LoadEnv(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Split the line into key and value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // Invalid line format
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		err := os.Setenv(key, value)
		if err != nil {
			return err
		} // Set the environment variable
	}

	return scanner.Err()
}

func GetEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
