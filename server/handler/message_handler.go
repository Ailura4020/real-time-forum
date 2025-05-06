package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"real-time-forum/models"
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

//	func GetMessagesHandler(db *sql.DB) http.HandlerFunc {
//		return func(w http.ResponseWriter, r *http.Request) {
//			userID, err := utils.ExtractUserIDFromRequest(r)
//			if err != nil {
//				http.Error(w, "Unauthorized", http.StatusUnauthorized)
//				return
//			}
//
//			// receiverIDStr := strings.TrimPrefix(r.URL.Path, "/api/messages/")
//			// receiverID, err := strconv.Atoi(receiverIDStr)
//			// if err != nil {
//			// 	http.Error(w, "Invalid receiver ID", http.StatusBadRequest)
//			// 	return
//			// }
//
//			receiverIDStr := r.URL.Query().Get("receiverId")
//			if receiverIDStr == "" {
//				http.Error(w, "Missing receiver ID", http.StatusBadRequest)
//				return
//			}
//			receiverID, err := strconv.Atoi(receiverIDStr)
//			if err != nil {
//				http.Error(w, "Invalid receiver ID", http.StatusBadRequest)
//				return
//			}
//			if receiverID == userID {
//				http.Error(w, "Cannot fetch messages with yourself", http.StatusBadRequest)
//				return
//			}
//
//			messages, err := repository.GetConversationMessages(db, userID, receiverID)
//			if err != nil {
//				http.Error(w, "Failed to fetch messages", http.StatusInternalServerError)
//				return
//			}
//
//			w.Header().Set("Content-Type", "application/json")
//			json.NewEncoder(w).Encode(messages)
//		}
//	}

// GetMessagesHandler handles retrieving conversation history
func GetMessagesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the authenticated user's ID
		userID, err := utils.ExtractUserIDFromRequest(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Get the receiver ID from URL path
		pathParts := utils.SplitPath(r.URL.Path)
		if len(pathParts) < 3 {
			http.Error(w, "Invalid URL path", http.StatusBadRequest)
			return
		}

		receiverIDStr := pathParts[len(pathParts)-1]
		receiverID, err := strconv.Atoi(receiverIDStr)
		if err != nil {
			http.Error(w, "Invalid receiver ID", http.StatusBadRequest)
			return
		}

		// Get conversation messages
		messages, err := repository.GetConversationMessages(db, userID, receiverID)
		if err != nil {
			http.Error(w, "Failed to fetch messages", http.StatusInternalServerError)
			return
		}

		// Return as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(messages)
	}
}

// SendMessageHandler handles sending a new private message
func SendMessageHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the authenticated user's ID
		userID, err := utils.ExtractUserIDFromRequest(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Parse the message from request body
		var messageRequest struct {
			ReceiverID int    `json:"receiver_id"`
			Content    string `json:"content"`
		}

		if err := json.NewDecoder(r.Body).Decode(&messageRequest); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate message content
		if messageRequest.Content == "" {
			http.Error(w, "Message content cannot be empty", http.StatusBadRequest)
			return
		}

		// Save the message
		err = repository.SavePrivateMessage(db, userID, messageRequest.ReceiverID, messageRequest.Content)
		if err != nil {
			http.Error(w, "Failed to save message", http.StatusInternalServerError)
			return
		}

		// Return success response
		response := models.Response{
			Success: true,
			Message: "Message sent successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
