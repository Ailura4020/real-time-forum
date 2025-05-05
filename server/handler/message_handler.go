package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"real-time-forum/repository"
	"real-time-forum/utils"
	"strconv"
	"strings"
	"time"
)

func savePrivateMessage(db *sql.DB, senderID int, receiverId int, content string) error {
	dateSent := time.Now().Format(time.RFC3339)
	query := `INSERT INTO PRIVATEMESSAGE (TextContent, DateSent, SenderId, ReceiverId) VALUES ( ?, ?, ?, ?)`
	_, err := db.Exec(query, content, dateSent, senderID, receiverId)
	return err
}

func GetMessagesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := utils.ExtractUserIDFromRequest(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		receiverIDStr := strings.TrimPrefix(r.URL.Path, "/api/messages/")
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

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(messages)
	}
}
