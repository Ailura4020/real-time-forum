package api

import (
	"database/sql"
	"log"
	"net/http"
	"real-time-forum/handler"
	middleware "real-time-forum/middlware"
	"real-time-forum/utils"
	"strings"
)

func RegisterRoutes(db *sql.DB, errorLogger *log.Logger) http.Handler {
	mux := http.NewServeMux()

	// Static file server (images, icons, etc.)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Serve static landing page
	mux.Handle("/", http.FileServer(http.Dir("./static")))

	// websocket for the CHAT
	hub := utils.NewHub()
	mux.HandleFunc("/ws", handler.HandleWebSocket(hub, db))

	// API routes
	mux.HandleFunc("/api/user", middleware.ErrorHandler(handler.UserHandler(db), errorLogger))

	mux.HandleFunc("/api/logout", middleware.ErrorHandler(handler.LogoutHandler, errorLogger))

	// POST /api/register
	mux.HandleFunc("/api/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			middleware.ErrorHandler(handler.RegisterHandler(db), errorLogger)(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// POST /api/login
	mux.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			middleware.ErrorHandler(handler.LoginHandler(db), errorLogger)(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	//// GET /api/posts/{id} – manual parsing of ID
	//mux.HandleFunc("/api/posts/", func(w http.ResponseWriter, r *http.Request) {
	//	if r.Method != "GET" {
	//		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	//		return
	//	}
	//
	//	path := strings.TrimPrefix(r.URL.Path, "/api/posts/")
	//	if path == "" {
	//		http.Error(w, "Post ID is required", http.StatusBadRequest)
	//		return
	//	}
	//
	//	// Attach ID via context or request URL
	//	r.URL.RawQuery = "id=" + path
	//	middleware.ErrorHandler(handler.GetPostHandler(db), errorLogger)(w, r)
	//})

	// Handle /api/posts
	mux.HandleFunc("/api/posts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			middleware.ErrorHandler(handler.GetPostsHandler(db), errorLogger)(w, r)
		} else if r.Method == "POST" {
			middleware.ErrorHandler(handler.CreatePostHandler(db), errorLogger)(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// Handle /api/posts/{id} for retrieving a specific post
	mux.HandleFunc("/api/posts/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/api/posts/")
		if path == "" {
			http.Error(w, "Post ID is required", http.StatusBadRequest)
			return
		}

		// Call the handler to get the specific post by ID
		middleware.ErrorHandler(handler.GetPostHandler(db), errorLogger)(w, r)
	})

	// POST /api/comments
	mux.HandleFunc("/api/comments", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			middleware.ErrorHandler(handler.AddCommentHandler(db), errorLogger)(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// GET /api/messages/{receiverID}
	//mux.HandleFunc("/api/messages/", middleware.ErrorHandler(handler.GetMessagesHandler(db), errorLogger))

	mux.HandleFunc("/api/messages/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			middleware.ErrorHandler(handler.GetMessagesHandler(db), errorLogger)(w, r)
		case "POST":
			middleware.ErrorHandler(handler.SendMessageHandler(db), errorLogger)(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// Protected route
	mux.HandleFunc("/api/protected", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		protected := middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			handler.SendResponse(w, true, "You accessed a protected route", nil, "")
		})
		middleware.ErrorHandler(protected, errorLogger)(w, r)
	})

	return mux
}
