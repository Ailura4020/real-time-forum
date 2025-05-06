package models

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"real-time-forum/repository"
	"real-time-forum/utils"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type PrivateMessage struct {
	ID         int       `json:"id"`
	SenderID   int       `json:"sender_id"`
	ReceiverID int       `json:"receiver_id"`
	Content    string    `json:"content"`
	Timestamp  time.Time `json:"timestamp"`
}

func GetMessagesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := utils.ExtractUserIDFromRequest(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		vars := mux.Vars(r)
		receiverIDStr, ok := vars["receiverId"]
		if !ok {
			http.Error(w, "Missing receiver ID", http.StatusBadRequest)
			return
		}

		receiverID, err := strconv.Atoi(receiverIDStr)
		if err != nil {
			http.Error(w, "Invalid receiver ID", http.StatusBadRequest)
			return
		}

		messages, err := repository.GetConversationMessages(db, userID, receiverID)
		if err != nil {
			http.Error(w, "Failed to fetch messages", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(messages)
	}
}

// func GetMessagesHandler(db *sql.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		userID, err := utils.ExtractUserIDFromRequest(r)
// 		if err != nil {
// 			http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			return
// 		}

// 		vars := mux.Vars(r)
// 		receiverIDStr, ok := vars["receiverId"]
// 		if !ok {
// 			http.Error(w, "Missing receiver ID", http.StatusBadRequest)
// 			return
// 		}

// 		receiverID, err := strconv.Atoi(receiverIDStr)
// 		if err != nil {
// 			http.Error(w, "Invalid receiver ID", http.StatusBadRequest)
// 			return
// 		}

// 		messages, err := repository.GetConversationMessages(db, userID, receiverID)
// 		if err != nil {
// 			http.Error(w, "Failed to fetch messages", http.StatusInternalServerError)
// 			return
// 		}

// 		json.NewEncoder(w).Encode(messages)
// 	}
// }
