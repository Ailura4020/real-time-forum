package repository

import (
	"database/sql"
	"log"

	"real-time-forum/models"
)

// UserRepository provides methods to interact with the user data
type UserRepository struct {
	DB *sql.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// CreateUser inserts a new user into the database
func (r *UserRepository) CreateUser(user models.User) (int64, error) {
	query := `
		INSERT INTO USERS (nickname, age, gender, first_name, last_name, email, password, date_register)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := r.DB.Exec(query, user.Nickname, user.Age, user.Gender, user.FirstName, user.LastName, user.Email, user.Password, user.DateRegister)
	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		return 0, err
	}
	return result.LastInsertId()
}

// GetUserByEmail retrieves a user by email
func (r *UserRepository) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	query := `SELECT id, nickname, age, gender, first_name, last_name, email, password, date_register FROM USERS WHERE email = ?`
	err := r.DB.QueryRow(query, email).Scan(&user.ID, &user.Nickname, &user.Age, &user.Gender, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.DateRegister)
	if err != nil {
		return user, err
	}
	return user, nil
}
