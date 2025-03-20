package handler

import (
	"encoding/json"
	"net/http"
	"real-time-forum/models"
)

// SendErrorResponse sends an error response to the client
func SendErrorResponse(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := models.Response{
		Success: false,
		Message: err.Error(),
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

// SendResponse sends a JSON response to the client
func SendResponse(w http.ResponseWriter, success bool, message string, data interface{}, token string) {
	w.Header().Set("Content-Type", "application/json")

	response := models.Response{
		Success: success,
		Message: message,
		Data:    data,
		Token:   token,
	}

	json.NewEncoder(w).Encode(response)
}
