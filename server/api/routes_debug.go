package api

import (
	"log"
	"net/http"
	"real-time-forum/handler"
	middleware "real-time-forum/middleware"
)

//func DevRoutes(router *mux.Router, errorLogger *log.Logger) {
//	// FOR TESTING PURPOSE ONLY (old method)
//	//mux.HandleFunc("/bad-request", middleware.ErrorHandler(handler.BadRequestHandler)) // 400 Bad Request
//	//mux.HandleFunc("/not-found", middleware.ErrorHandler(handler.NotFoundHandler))     // 404 Not Found
//	//mux.HandleFunc("/panic", middleware.ErrorHandler(handler.PanicHandler))            // 500 Panic
//
//	// FOR TESTING PURPOSE ONLY (new method) /!\ not for production
//	router.HandleFunc("/test", middleware.ErrorHandler(handler.TestHandler, errorLogger)).Methods("GET", "POST", "PUT", "DELETE") // test route
//
//	router.HandleFunc("/bad-request", middleware.ErrorHandler(handler.BadRequestHandler, errorLogger)).Methods("GET") // 400 Bad Request
//	router.HandleFunc("/not-found", middleware.ErrorHandler(handler.NotFoundHandler, errorLogger)).Methods("GET")     // 404 Not Found
//	router.HandleFunc("/panic", middleware.ErrorHandler(handler.PanicHandler, errorLogger)).Methods("GET")            // 500 Panic
//}

func DevRoutes(errorLogger *log.Logger) http.Handler {
	mux := http.NewServeMux()

	// FOR TESTING PURPOSE ONLY (do not use in production)

	// /test supports all major methods
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		middleware.ErrorHandler(handler.TestHandler, errorLogger)(w, r)
	})

	// /bad-request
	mux.HandleFunc("/bad-request", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			middleware.ErrorHandler(handler.BadRequestHandler, errorLogger)(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// /not-found
	mux.HandleFunc("/not-found", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			middleware.ErrorHandler(handler.NotFoundHandler, errorLogger)(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// /panic
	mux.HandleFunc("/panic", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			middleware.ErrorHandler(handler.PanicHandler, errorLogger)(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}
