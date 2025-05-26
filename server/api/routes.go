package api

import (
	"database/sql"
	"log"
	"net/http"
	"real-time-forum/handler"
	middleware "real-time-forum/middleware"
	"real-time-forum/utils"
	"strings"
	"fmt"
)

func RegisterRoutes(db *sql.DB, errorLogger *log.Logger) http.Handler {
	mux := http.NewServeMux()

	// Static file server for frontend assets
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Public API endpoints
	mux.HandleFunc("/api/login", methodHandler("POST", handler.LoginHandler(db), errorLogger))
	mux.HandleFunc("/api/register", methodHandler("POST", handler.RegisterHandler(db), errorLogger))

	// Authentication check endpoint
	mux.HandleFunc("/api/auth/check", func(w http.ResponseWriter, r *http.Request) {
		userID, err := utils.ExtractUserIDFromRequest(r)
		if err != nil {
			handler.SendResponse(w, false, "Not authenticated", nil, "")
		} else {
			handler.SendResponse(w, true, "Authenticated", map[string]interface{}{"user_id": userID}, "")
		}
	})

	// Protected routes using existing AuthMiddleware
	mux.HandleFunc("/api/user", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		middleware.ErrorHandler(handler.UserHandler(db), errorLogger)(w, r)
	}))

	mux.HandleFunc("/api/logout", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		middleware.ErrorHandler(handler.LogoutHandler, errorLogger)(w, r)
	}))

	// Websocket for chat
	hub := utils.NewHub()
	mux.HandleFunc("/ws", middleware.WebSocketAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("WebSocketAuthMiddleware: route")
		handler.HandleWebSocket(hub, db)(w, r)
	}))
	
	// Posts
	mux.HandleFunc("/api/posts", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			middleware.ErrorHandler(handler.GetPostsHandler(db), errorLogger)(w, r)
		} else if r.Method == "POST" {
			middleware.ErrorHandler(handler.CreatePostHandler(db), errorLogger)(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}))

	// Individual post
	mux.HandleFunc("/api/posts/", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/posts/")
		if path == "" {
			http.Error(w, "Post ID is required", http.StatusBadRequest)
			return
		}
		middleware.ErrorHandler(handler.GetPostHandler(db), errorLogger)(w, r)
	}))

	// Comments
	mux.HandleFunc("/api/comments", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			middleware.ErrorHandler(handler.AddCommentHandler(db), errorLogger)(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}))

	// Messages
	mux.HandleFunc("/api/messages/", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
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
	}))

	// All conversations
	mux.HandleFunc("/api/messages/all", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			middleware.ErrorHandler(handler.GetAllConversationsHandler(db), errorLogger)(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}))

	return mux
}

func methodHandler(method string, handlerFunc http.HandlerFunc, errorLogger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		middleware.ErrorHandler(handlerFunc, errorLogger)(w, r)
	}
}
