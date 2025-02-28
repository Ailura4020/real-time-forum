package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"real-time-forum/repository"

	"real-time-forum/models"
	"real-time-forum/service"
	"real-time-forum/utils" // Import the utils package
)

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

// SendResponse sends a JSON response
func SendResponse(w http.ResponseWriter, success bool, message string, data interface{}, token string) {
	response := models.Response{
		Success: success,
		Message: message,
		Data:    data,
		Token:   token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
