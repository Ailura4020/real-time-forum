package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"real-time-forum/repository"
	"strings"

	"real-time-forum/models"
	"real-time-forum/service"
	"real-time-forum/utils" // Import the utils package
)

// UserHandler handles requests for retrieving user information
// using the authorization token
func UserHandler(db *sql.DB) http.HandlerFunc {
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received %s request for %s", r.Method, r.URL.Path)
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			SendResponse(w, false, "Authorization header required", nil, "")
			return
		}

		// Extract the token from the Authorization header
		// Expected format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			SendResponse(w, false, "Invalid authorization format", nil, "")
			return
		}
		token := parts[1]

		// Validate the token and get user ID
		claims, err := utils.ValidateJWT(token)
		if err != nil {
			log.Printf("Token validation error: %v", err)
			SendResponse(w, false, "Invalid or expired token", nil, "")
			return
		}

		// Get user details using the ID from the token
		user, err := userService.GetUserByID(claims.ID)
		if err != nil {
			log.Printf("Error fetching user: %v", err)
			SendResponse(w, false, "User not found", nil, "")
			return
		}

		// User is authenticated and found
		SendResponse(w, true, "User authorized", user, token)
	}
}

// RegisterHandler handles user registration
func RegisterHandler(db *sql.DB) http.HandlerFunc {
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received %s request for %s", r.Method, r.URL.Path)
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req models.RegisterRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// Basic validation
		if req.Email == "" || req.Password == "" || req.Nickname == "" {
			SendResponse(w, false, "Missing required fields", nil, "")
			return
		}

		// Register user
		user, err := userService.RegisterUser(req)
		if err != nil {
			SendResponse(w, false, err.Error(), nil, "")
			return
		}

		token, err := utils.GenerateJWT(user)
		if err != nil {
			SendResponse(w, true, "User registered successfully", user, "")
			return
		}

		SendResponse(w, true, "User registered successfully", user, token)
	}
}

// LoginHandler handles user login
func LoginHandler(db *sql.DB) http.HandlerFunc {
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received %s request for %s", r.Method, r.URL.Path)
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req models.LoginRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// Login user
		user, err := userService.LoginUser(req)
		if err != nil {
			log.Printf("Login failed for user %s: %v", req.Email, err)
			SendResponse(w, false, err.Error(), nil, "")
			return
		}

		token, err := utils.GenerateJWT(user)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		SendResponse(w, true, "Login successful", user, token)
	}
}

//
//// SendResponse sends a JSON response
//func SendResponse(w http.ResponseWriter, success bool, message string, data interface{}, token string) {
//	response := models.Response{
//		Success: success,
//		Message: message,
//		Data:    data,
//		Token:   token,
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	json.NewEncoder(w).Encode(response)
//}
