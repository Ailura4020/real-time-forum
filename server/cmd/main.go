package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	"real-time-forum/config"
	"real-time-forum/db"
	"real-time-forum/handler"
	"real-time-forum/middlware"
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

	// like it said, init the config
	appEnv, err := initConfig()
	if err != nil {
		errorLogger.Printf("Error loading .env file: %v", err)
		//log.Fatalf("Failed to initialize configuration: %v", err)
	}

	// fill the environment variables
	addr := appEnv.ServerAddr
	secretFromEnv := appEnv.Token

	log.Printf("Starting %v server...\n", addr)

	// set JWTSecret and warn if empty
	if secretFromEnv != "" {
		config.SetJWTSecret(secretFromEnv)
	} else {
		log.Printf("Warning: You should set JWT_SECRET environment variable:\n\t\t    export JWT_SECRET=\"your secret\"\n")
	}

	// init the database from the sql statements (schema.sql)
	initDB := db.InitDB()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(initDB)

	// define mux (old method)
	//mux := http.NewServeMux()
	// Define routes (new method w/ gorilla)
	router := mux.NewRouter()

	// static (for user thuumbnail?) http://localhost:8000/static/<filename>
	//router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))

	// API endpoints (old method)
	//mux.HandleFunc("/api/register", middleware.ErrorHandler(handler.RegisterHandler(initDB)))
	//mux.HandleFunc("/api/login", middleware.ErrorHandler(handler.LoginHandler(initDB)))

	// todo: API endpoints (new method)
	//router.HandleFunc("/api/register", middleware.ErrorHandler(handler.RegisterHandler(initDB))).Methods("POST")
	//router.HandleFunc("/api/login", middleware.ErrorHandler(handler.LoginHandler(initDB))).Methods("POST")
	//router.HandleFunc("/api/posts", middleware.ErrorHandler(handler.GetPostsHandler(initDB))).Methods("GET")
	//router.HandleFunc("/api/posts/{id}", middleware.ErrorHandler(handler.GetPostHandler(initDB))).Methods("GET")
	//router.HandleFunc("/api/posts", middleware.ErrorHandler(handler.CreatePostHandler(initDB))).Methods("POST")
	//router.HandleFunc("/api/comments", middleware.ErrorHandler(handler.AddCommentHandler(initDB))).Methods("POST")

	// todo: add more endpoints
	//POST /api/login - User login @done
	//POST /api/register - User registration @done
	//GET /api/posts - Retrieve all posts
	//POST /api/posts - Create a new post
	//GET /api/posts/:id - Retrieve a specific post
	//POST /api/comments - Add a comment to a post
	// websocket for the CHAT?

	// FOR TESTING PURPOSE ONLY
	//mux.HandleFunc("/bad-request", middleware.ErrorHandler(handler.BadRequestHandler)) // 400 Bad Request
	//mux.HandleFunc("/not-found", middleware.ErrorHandler(handler.NotFoundHandler))     // 404 Not Found
	//mux.HandleFunc("/panic", middleware.ErrorHandler(handler.PanicHandler))            // 500 Panic

	// FOR TESTING PURPOSE ONLY (new method)
	//router.HandleFunc("/bad-request", middleware.ErrorHandler(handler.BadRequestHandler)).Methods("GET") // 400 Bad Request
	router.HandleFunc("/api/register", middleware.ErrorHandler(handler.BadRequestHandler, errorLogger)).Methods("POST")
	router.HandleFunc("/not-found", middleware.ErrorHandler(handler.NotFoundHandler, errorLogger)).Methods("GET") // 404 Not Found
	router.HandleFunc("/panic", middleware.ErrorHandler(handler.PanicHandler, errorLogger)).Methods("GET")        // 500 Panic
	// todo: add a protected route (endpoint that requires authentication and/or authorization to access > ex: for CRUD operations)
	//mux.HandleFunc("/api/protected", middleware.ErrorHandler(middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
	//	handler.SendResponse(w, true, "You accessed a protected route", nil, "")
	//})))

	// Apply CORS and RateLimit middleware (old method)
	//handler := middleware.CORSMiddleware(mux)
	//handler := middleware.CORSMiddleware(middleware.RateLimit(mux)) // Chaining Middleware

	// Apply global middleware (new method)
	router.Use(middleware.SecurityHeaders) // protect your application from various attacks (like XSS, clickjacking, etc.)
	router.Use(middleware.CORSMiddleware)  // handling cross-origin requests
	router.Use(middleware.RateLimit)       // ensure that it can track and limit requests effectively

	// Start server
	srv := &http.Server{
		Handler:        router,
		Addr:           addr,
		WriteTimeout:   15 * time.Second,
		ReadTimeout:    15 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,
		ErrorLog:       errorLogger,
	}
	fmt.Printf("Server running on http://localhost%s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
