package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"real-time-forum/api"
	"real-time-forum/demo"
	"time"

	"real-time-forum/config"
	"real-time-forum/db"
	middleware "real-time-forum/middleware"
)

// AppConfig holds the application configuration
type appConfig struct {
	ServerAddr string
	Token      string
}

// initConfig initializes the application configuration from the .env file (/root)
func initConfig() (*appConfig, error) {
	// Load environment variables from .env file
	if err := config.LoadEnv("../.env"); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	// Retrieve configuration values
	serverAddr := config.GetEnv("SERVER_ADDR", ":8080")
	token := config.GetEnv("JWT_SECRET", "")
	// ...

	// fill the struct
	return &appConfig{
		ServerAddr: serverAddr,
		Token:      token,
	}, nil
}

func main() {
	// check if first init
	first := false
	if _, err := os.Stat("./db/forum.db"); os.IsNotExist(err) {
		fmt.Println("Database does not exist")
		first = true
	}

	// dedicated logger for the server errors
	errorLogger := log.New(os.Stderr, "ERROR: ", log.LstdFlags)

	// init the config
	appEnv, err := initConfig()
	if err != nil {
		errorLogger.Printf("Error loading .env file: %v", err)
		//log.Fatalf("Failed to initialize configuration: %v", err)
	}

	// fill the environment variables from the .env file
	addr := appEnv.ServerAddr
	secretFromEnv := appEnv.Token

	log.Printf("Starting %v server...\n", addr)

	// set JWTSecret and warn if empty
	if secretFromEnv != "" {
		config.SetJWTSecret(secretFromEnv)
	} else {
		log.Printf("Warning: You should set JWT_SECRET environment variable.")
	}

	// init the database from the sql statements (schema.sql)
	initDB := db.InitDB()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("Error closing db: %v", err)
		}
	}(initDB)

	// DEMO
	if first {
		fmt.Println("DEMO MODE: Registering users and posts in the database...")
		csvFilePathUsers := "./demo/users.csv"
		demo.RegisterUsers(initDB, csvFilePathUsers)
		csvFilePathPosts := "./demo/posts.csv"
		demo.RegisterPosts(initDB, csvFilePathPosts)
	}

	// Define routes
	finalHandler := middleware.Chain(
		//api.DevRoutes(errorLogger), // this returns http.Handler
		api.RegisterRoutes(initDB, errorLogger), // this returns http.Handler
		middleware.SecurityHeaders,
		middleware.CORSMiddleware,
		middleware.RateLimit,
	)

	//// testing routes (not for production)
	//api.DevRoutes(router, errorLogger)
	//// Register routes with the initialized DB connection
	//api.RegisterRoutes(router, initDB, errorLogger)

	// Apply global middleware (new chaining method w/Gorilla)
	//router.Use(middleware.SecurityHeaders) // protect your application from various attacks (like XSS, clickjacking, etc.)
	//router.Use(middleware.CORSMiddleware)  // handling cross-origin requests
	//router.Use(middleware.RateLimit)       // ensure that it can track and limit requests effectively

	// testing routes (not for production)
	//api.DevRoutes(router, errorLogger)
	// Register routes with the initialized DB connection
	//api.RegisterRoutes(router, initDB, errorLogger)

	// Start server
	server := &http.Server{
		Addr: addr,
		// Good practice to set timeouts to avoid "Slow Loris" attacks [wiki](https://en.wikipedia.org/wiki/Slowloris_(cyber_attack)).
		WriteTimeout:   time.Second * 15,
		ReadTimeout:    time.Second * 15,
		IdleTimeout:    time.Second * 60,
		MaxHeaderBytes: 4096, // 4 KB
		ErrorLog:       errorLogger,
		Handler:        finalHandler, // pass instance to Gorilla mux
	}

	fmt.Printf("Server running on http://localhost%s\n", server.Addr) // valid only on local

	log.Fatal(server.ListenAndServe())
}

//MaxHeaderBytes: 1 << 20, // 1048576 bytes, or 1 megabyte (MB)
//MaxHeaderBytes: 16384, // 16 KB
