package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"real-time-forum/repository"
	"real-time-forum/utils"
	"strconv"
)

func savePrivateMessage(db *sql.DB, senderID int, receiverID int, content string) error {
	query := `
        INSERT INTO private_messages (sender_id, receiver_id, content)
        VALUES (?, ?, ?)
    `
	_, err := db.Exec(query, senderID, receiverID, content)
	if err != nil {
		return fmt.Errorf("failed to insert private message: %w", err)
	}
	return nil
}

func GetMessagesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := utils.ExtractUserIDFromRequest(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// receiverIDStr := strings.TrimPrefix(r.URL.Path, "/api/messages/")
		// receiverID, err := strconv.Atoi(receiverIDStr)
		// if err != nil {
		// 	http.Error(w, "Invalid receiver ID", http.StatusBadRequest)
		// 	return
		// }

		receiverIDStr := r.URL.Query().Get("receiverId")
		if receiverIDStr == "" {
			http.Error(w, "Missing receiver ID", http.StatusBadRequest)
			return
		}
		receiverID, err := strconv.Atoi(receiverIDStr)
		if err != nil {
			http.Error(w, "Invalid receiver ID", http.StatusBadRequest)
			return
		}
		if receiverID == userID {
			http.Error(w, "Cannot fetch messages with yourself", http.StatusBadRequest)
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
