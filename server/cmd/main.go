package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"real-time-forum/api"
	"time"

	"github.com/gorilla/mux"

	"real-time-forum/config"
	"real-time-forum/db"
	middleware "real-time-forum/middlware"
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

	if err := db.InitDB(); err != nil {
		log.Fatalf("Erreur lors de l'initialisation de la base de données: %v", err)
	}
	defer func() {
		if db.DB != nil {
			db.DB.Close()
		}
	}()

	// Define routes (new method w/ gorilla)
	router := mux.NewRouter()

	//// testing routes (not for production)
	//api.DevRoutes(router, errorLogger)
	//// Register routes with the initialized DB connection
	//api.RegisterRoutes(router, initDB, errorLogger)

	// Apply global middleware (new chaining method w/Gorilla)
	router.Use(middleware.SecurityHeaders) // protect your application from various attacks (like XSS, clickjacking, etc.)
	router.Use(middleware.CORSMiddleware)  // handling cross-origin requests
	router.Use(middleware.RateLimit)       // ensure that it can track and limit requests effectively

	// testing routes (not for production)
	api.DevRoutes(router, errorLogger)
	// Register routes with the initialized DB connection
	api.RegisterRoutes(router, db.DB, errorLogger)

	// Start server
	server := &http.Server{
		Addr: addr,
		// Good practice to set timeouts to avoid "Slow Loris" attacks [wiki](https://en.wikipedia.org/wiki/Slowloris_(cyber_attack)).
		WriteTimeout:   time.Second * 15,
		ReadTimeout:    time.Second * 15,
		IdleTimeout:    time.Second * 60,
		MaxHeaderBytes: 4096, // 4 KB
		ErrorLog:       errorLogger,
		Handler:        router, // pass instance to Gorilla mux
	}

	fmt.Printf("Server running on http://localhost%s\n", server.Addr) // valid only on local

	log.Fatal(server.ListenAndServe())
}

//MaxHeaderBytes: 1 << 20, // 1048576 bytes, or 1 megabyte (MB)
//MaxHeaderBytes: 16384, // 16 KB
