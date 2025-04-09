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
	router.HandleFunc("/ws", middleware.ErrorHandler(handler.HandleWebsocket(db, hub), errorLogger)).Methods("GET", "OPTIONS")

	// Posts (GET /api/posts or POST /api/posts)
	mux.HandleFunc("/api/posts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			middleware.ErrorHandler(handler.GetPostsHandler(db), errorLogger)(w, r)
		} else if r.Method == "POST" {
			middleware.ErrorHandler(handler.CreatePostHandler(db), errorLogger)(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// GET /api/posts/{id} – manual parsing of ID
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

		// Attach ID via context or request URL
		r.URL.RawQuery = "id=" + path
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

	// WebSocket route
	mux.HandleFunc("/ws", handler.HandleWebSocket)

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
