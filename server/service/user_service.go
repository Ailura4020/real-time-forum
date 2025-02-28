package service

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"real-time-forum/models"
	"real-time-forum/repository"
	"time"
)

// UserService provides methods for user-related operations
type UserService struct {
	UserRepo *repository.UserRepository
}

// NewUserService creates a new UserService
func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{UserRepo: userRepo}
}

// RegisterUser handles user registration
func (s *UserService) RegisterUser(req models.RegisterRequest) (models.User, error) {
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return models.User{}, err
	}

	user := models.User{
		Nickname:     req.Nickname,
		Age:          req.Age,
		Gender:       req.Gender,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		Password:     hashedPassword,
		DateRegister: time.Now().Format(time.RFC3339),
	}

	id, err := s.UserRepo.CreateUser(user)
	if err != nil {
		return models.User{}, err
	}
	user.ID = int(id)
	return user, nil
}

// LoginUser handles user login
func (s *UserService) LoginUser(req models.LoginRequest) (models.User, error) {
	user, err := s.UserRepo.GetUserByEmail(req.Email)
	if err != nil {
		return models.User{}, err
	}
	if !CheckPasswordHash(req.Password, user.Password) {
		return models.User{}, fmt.Errorf("invalid credentials")
	}
	return user, nil
}

// HashPassword creates a bcrypt hash of the password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash compares a password with a hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
