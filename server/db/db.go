package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// Variable globale pour stocker la connexion DB
var DB *sql.DB

// InitDB initializes the database and creates tables if they don't exist
func InitDB() error {
	_, err := os.Stat("db/forum.db")
	dbExists := !os.IsNotExist(err)

	var errDB error
	DB, err = sql.Open("sqlite3", "db/forum.db") // Stocker la connexion dans DB
	if errDB != nil {
		log.Fatalf("Failed to open database: %v", err)
		return errDB
	}

	if !dbExists {
		createTablesFromSQLFile(DB)
	}
	return nil
}

// Create necessary tables in the database
func createTablesFromSQLFile(db *sql.DB) {
	sqlFile, err := os.ReadFile("db/schema.sql")
	if err != nil {
		log.Fatalf("Failed to read SQL file: %v", err)
	}

	if _, err = db.Exec(string(sqlFile)); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	log.Println("Database tables created successfully")
}
