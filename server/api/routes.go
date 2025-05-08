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

	// Static file server for frontend assets - accessible to everyone
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Public API endpoints
	// Login
	mux.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			middleware.ErrorHandler(handler.LoginHandler(db), errorLogger)(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// Register
	mux.HandleFunc("/api/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			middleware.ErrorHandler(handler.RegisterHandler(db), errorLogger)(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// Authentication check endpoint - useful for frontend to verify auth status
	mux.HandleFunc("/api/auth/check", func(w http.ResponseWriter, r *http.Request) {
		userID, err := utils.ExtractUserIDFromRequest(r)
		if err != nil {
			// Not authenticated
			handler.SendResponse(w, false, "Not authenticated", nil, "")
		} else {
			// Authenticated
			handler.SendResponse(w, true, "Authenticated", map[string]interface{}{"user_id": userID}, "")
		}
	})

	// Allow public access to the main app entry points for the frontend
	// Frontend will handle redirecting unauthenticated users
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// For root path or specific routes we want to be public (login, register, about)
		if r.URL.Path == "/" || r.URL.Path == "/login" || r.URL.Path == "/register" || r.URL.Path == "/about" {
			http.ServeFile(w, r, "./static/index.html")
			return
		}

		// For all other frontend routes, check authentication
		_, err := utils.ExtractUserIDFromRequest(r)
		if err != nil {
			// Return a 401 so frontend can redirect
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// User is authenticated, serve the main app
		http.ServeFile(w, r, "./static/index.html")
	})

	// Protected API routes
	// User info
	mux.HandleFunc("/api/user", func(w http.ResponseWriter, r *http.Request) {
		_, err := utils.ExtractUserIDFromRequest(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		middleware.ErrorHandler(handler.UserHandler(db), errorLogger)(w, r)
	})

	// Logout
	mux.HandleFunc("/api/logout", func(w http.ResponseWriter, r *http.Request) {
		_, err := utils.ExtractUserIDFromRequest(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		middleware.ErrorHandler(handler.LogoutHandler, errorLogger)(w, r)
	})

	// Websocket for chat
	hub := utils.NewHub()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		_, err := utils.ExtractUserIDFromRequest(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		handler.HandleWebSocket(hub, db)(w, r)
	})

	// Posts
	mux.HandleFunc("/api/posts", func(w http.ResponseWriter, r *http.Request) {
		_, err := utils.ExtractUserIDFromRequest(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if r.Method == "GET" {
			middleware.ErrorHandler(handler.GetPostsHandler(db), errorLogger)(w, r)
		} else if r.Method == "POST" {
			middleware.ErrorHandler(handler.CreatePostHandler(db), errorLogger)(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// Individual post
	mux.HandleFunc("/api/posts/", func(w http.ResponseWriter, r *http.Request) {
		_, err := utils.ExtractUserIDFromRequest(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if r.Method != "GET" {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/api/posts/")
		if path == "" {
			http.Error(w, "Post ID is required", http.StatusBadRequest)
			return
		}

		middleware.ErrorHandler(handler.GetPostHandler(db), errorLogger)(w, r)
	})

	// Comments
	mux.HandleFunc("/api/comments", func(w http.ResponseWriter, r *http.Request) {
		_, err := utils.ExtractUserIDFromRequest(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if r.Method == "POST" {
			middleware.ErrorHandler(handler.AddCommentHandler(db), errorLogger)(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// Messages
	mux.HandleFunc("/api/messages/", func(w http.ResponseWriter, r *http.Request) {
		_, err := utils.ExtractUserIDFromRequest(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if r.Method == "GET" {
			path := r.URL.Path
			parts := strings.Split(path, "/")
			if len(parts) >= 3 && parts[2] != "" {
				middleware.ErrorHandler(handler.GetUserMessagesHandler(db), errorLogger)(w, r)
				return
			}
			http.Error(w, "Not Found", http.StatusNotFound)
		} else if r.Method == "POST" {
			middleware.ErrorHandler(handler.SendMessageHandler(db), errorLogger)(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// All conversations
	mux.HandleFunc("/api/messages/all", func(w http.ResponseWriter, r *http.Request) {
		_, err := utils.ExtractUserIDFromRequest(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if r.Method == "GET" {
			middleware.ErrorHandler(handler.GetAllConversationsHandler(db), errorLogger)(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}
