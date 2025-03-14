package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB initializes the database and creates tables if they don't exist
func InitDB() *sql.DB {
	_, err := os.Stat("db/forum.db")
	dbExists := !os.IsNotExist(err)

	db, err := sql.Open("sqlite3", "db/forum.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	if !dbExists {
		createTablesFromSQLFile(db)
	}

	return db
}

// Create necessary tables in the database
func createTablesFromSQLFile(db *sql.DB) {
	sqlFile, err := os.ReadFile("db/schema.sql")
	if err != nil {
		log.Fatalf("Failed to read SQL file: %v", err)
	}

	_, err = db.Exec(string(sqlFile))
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	log.Println("Database tables created successfully")
}
