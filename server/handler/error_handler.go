package handler

import (
	"encoding/json"
	"net/http"
	"real-time-forum/models"
)

// BadRequestHandler simulates a 400 Bad Request error.
func BadRequestHandler(w http.ResponseWriter, r *http.Request) {
	response := models.ErrorResponse{
		Code:    http.StatusBadRequest,
		Message: "Bad Request: The request could not be understood or was missing required parameters.",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(response)
}

// NotFoundHandler simulates a 404 Not Found error.
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	response := models.ErrorResponse{
		Code:    http.StatusNotFound,
		Message: "Not Found: The requested resource could not be found.",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(response)
}

// PanicHandler simulates a panic.
func PanicHandler(w http.ResponseWriter, r *http.Request) {
	panic("This is a test panic!")
}
