package demo

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"real-time-forum/models"
	"real-time-forum/repository"
)

// RegisterUsers reads users from a CSV file and registers them in the database
func RegisterUsers(db *sql.DB, csvFilePath string) {
	userRepo := repository.NewUserRepository(db)

	// Open the CSV file
	file, err := os.Open(csvFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read the CSV file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Iterate over the records and create users
	for i, record := range records {
		if i == 0 { // Skip header
			continue
		}

		// Convert age from string to int
		age, err := strconv.Atoi(record[1])
		if err != nil {
			log.Printf("Invalid age for user %s: %v", record[5], err)
			continue
		}

		gender := record[2]
		firstName := record[3]
		lastName := record[4]
		email := record[5]
		password := record[6]
		dateRegister := record[7]

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Failed to hash password for user %s: %v", email, err)
			continue
		}

		user := models.User{
			Nickname:     record[0],
			Age:          age, // Use the integer age
			Gender:       gender,
			FirstName:    firstName,
			LastName:     lastName,
			Email:        email,
			Password:     string(hashedPassword),
			DateRegister: dateRegister,
		}

		// Create the user in the database
		_, err = userRepo.CreateUser(user)
		if err != nil {
			log.Printf("Failed to create user %s: %v", email, err)
		} else {
			fmt.Printf("User %s registered successfully\n", user.Nickname)
		}
	}
}
