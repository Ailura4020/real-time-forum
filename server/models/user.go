package models

// User represents the user model
type User struct {
	ID           int    `json:"id"`
	Nickname     string `json:"nickname"`
	Age          int    `json:"age"`
	Gender       string `json:"gender"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Password     string `json:"-"` // Don't return password in JSON responses
	DateRegister string `json:"date_register"`
}

// RegisterRequest captures the registration data
type RegisterRequest struct {
	Nickname  string `json:"nickname"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// LoginRequest captures login credentials
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Response structure for API responses
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Token   string      `json:"token,omitempty"`
}
