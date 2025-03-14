package api

import (
	"github.com/gorilla/mux"
	"log"
	"real-time-forum/handler"
	middleware "real-time-forum/middlware"
)

func DevRoutes(router *mux.Router, errorLogger *log.Logger) {
	// FOR TESTING PURPOSE ONLY (old method)
	//mux.HandleFunc("/bad-request", middleware.ErrorHandler(handler.BadRequestHandler)) // 400 Bad Request
	//mux.HandleFunc("/not-found", middleware.ErrorHandler(handler.NotFoundHandler))     // 404 Not Found
	//mux.HandleFunc("/panic", middleware.ErrorHandler(handler.PanicHandler))            // 500 Panic

	// FOR TESTING PURPOSE ONLY (new method) /!\ not for production
	router.HandleFunc("/test", middleware.ErrorHandler(handler.TestHandler, errorLogger)).Methods("GET", "POST", "PUT", "DELETE") // test route

	router.HandleFunc("/bad-request", middleware.ErrorHandler(handler.BadRequestHandler, errorLogger)).Methods("GET") // 400 Bad Request
	router.HandleFunc("/not-found", middleware.ErrorHandler(handler.NotFoundHandler, errorLogger)).Methods("GET")     // 404 Not Found
	router.HandleFunc("/panic", middleware.ErrorHandler(handler.PanicHandler, errorLogger)).Methods("GET")            // 500 Panic
}
