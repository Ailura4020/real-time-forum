package api

import (
	"database/sql"
	"log"
	"net/http"
	"real-time-forum/handler"
	middleware "real-time-forum/middlware"

	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router, db *sql.DB, errorLogger *log.Logger) {
	dir := "./static"
	// static route for future usage (user thumbnail, images, ico) http://localhost:8000/static/<filename>
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))

	// Serve the static landing page
	router.Handle("/", http.FileServer(http.Dir("./static")))

	// todo: API endpoints (new method)
	router.HandleFunc("/api/user", middleware.ErrorHandler(handler.UserHandler(db), errorLogger)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/register", middleware.ErrorHandler(handler.RegisterHandler(db), errorLogger)).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/login", middleware.ErrorHandler(handler.LoginHandler(db), errorLogger)).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/posts", middleware.ErrorHandler(handler.GetPostsHandler(db), errorLogger)).Methods("GET")
	router.HandleFunc("/api/posts/{id}", middleware.ErrorHandler(handler.GetPostHandler(db), errorLogger)).Methods("GET")
	router.HandleFunc("/api/posts", middleware.ErrorHandler(handler.CreatePostHandler(db), errorLogger)).Methods("POST")
	router.HandleFunc("/api/comments", middleware.ErrorHandler(handler.AddCommentHandler(db), errorLogger)).Methods("POST")

	// websocket for the CHAT
	router.HandleFunc("/ws", handler.HandleWebSocket)

	// todo: add a protected route (endpoint that requires authentication and/or authorization to access > ex: for CRUD operations)
	router.HandleFunc("/api/protected", middleware.ErrorHandler(middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handler.SendResponse(w, true, "You accessed a protected route", nil, "")
	}), errorLogger)).Methods("GET")

}
