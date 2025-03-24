package demo

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	//"strconv"
	//"time"

	"database/sql"
	"real-time-forum/models"
	"real-time-forum/repository"
)

// RegisterPosts reads posts from a CSV file and registers them in the database
func RegisterPosts(db *sql.DB, csvFilePath string) {
	postRepo := repository.NewPostRepository(db)

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

	// Iterate over the records and create posts
	for i, record := range records {
		if i == 0 { // Skip header
			continue
		}

		title := record[0]
		category := record[1]
		textContent := record[2]
		userID, err := strconv.Atoi(record[3])
		//userID := record[3]
		if err != nil {
			log.Printf("Invalid user ID for post %s: %v", title, err)
			continue
		}
		dateCreation := record[4]

		post := models.Post{
			Title:        title,
			Category:     category,
			TextContent:  textContent,
			UserID:       userID,
			DateCreation: dateCreation,
		}

		// Create the post in the database
		_, err = postRepo.CreatePost(post)
		if err != nil {
			log.Printf("Failed to create post %s: %v", title, err)
		} else {
			fmt.Printf("Post '%s' registered successfully\n", post.Title)
		}
	}
}
