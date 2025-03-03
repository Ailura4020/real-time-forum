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
	router.HandleFunc("/api/register", middleware.ErrorHandler(handler.RegisterHandler(db), errorLogger)).Methods("POST")
	router.HandleFunc("/api/login", middleware.ErrorHandler(handler.LoginHandler(db), errorLogger)).Methods("POST")
	router.HandleFunc("/ws", handler.HandleWebSocket)
	//router.HandleFunc("/api/posts", middleware.ErrorHandler(handler.GetPostsHandler(initDB))).Methods("GET")
	//router.HandleFunc("/api/posts/{id}", middleware.ErrorHandler(handler.GetPostHandler(initDB))).Methods("GET")
	//router.HandleFunc("/api/posts", middleware.ErrorHandler(handler.CreatePostHandler(initDB))).Methods("POST")
	//router.HandleFunc("/api/comments", middleware.ErrorHandler(handler.AddCommentHandler(initDB))).Methods("POST")

	// todo: websocket for the CHAT?

	// todo: add a protected route (endpoint that requires authentication and/or authorization to access > ex: for CRUD operations)
	router.HandleFunc("/api/protected", middleware.ErrorHandler(middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handler.SendResponse(w, true, "You accessed a protected route", nil, "")
	}), errorLogger)).Methods("GET")

}
