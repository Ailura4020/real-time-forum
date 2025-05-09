package repository

import (
	"database/sql"
	"fmt"
	"real-time-forum/models"
	"time"
)

type PrivateMessageRepository struct {
	DB *sql.DB
}

// NewUserRepository creates a new UserRepository
func NewPrivateMessageRepository(db *sql.DB) *PrivateMessageRepository {
	return &PrivateMessageRepository{DB: db}
}

// insertion d'un msg privé
func (r *PrivateMessageRepository) InsetPrivateMessage(msg *models.PrivateMessage) error {
	query := `
	INSERT INTO private_messages (sender_id,receiver_id,content,timestamp)
	VALUES (?, ?, ?, ?)
	`
	_, err := r.DB.Exec(query, msg.SenderID, msg.ReceiverID, msg.Content, msg.Timestamp)
	return err
}

func (r *PrivateMessageRepository) GetPrivateMessages(user1, user2, limit, offset int) ([]models.PrivateMessage, error) {
	query := `
	SELECT id,sender_id,receiver_id,content,timestamp
	From private_messages
	WHERE 
	(sender_id = ? AND receiver_id =?)
	OR (sender_id =? AND receiver_id = ?)
	ORDER BY timestamp DESC
	LIMIT ? OFFSET ?
	`
	rows, err := r.DB.Query(query, user1, user2, user2, user1, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.PrivateMessage
	for rows.Next() {
		var msg models.PrivateMessage
		var timestampStr string

		err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.Timestamp)
		if err != nil {
			return nil, err
		}

		timestamp, err := time.Parse("2006-01-02 15:04:05", timestampStr)
		if err != nil {
			return nil, err
		}

		msg.Timestamp = timestamp

		messages = append(messages, msg)
	}
	return messages, nil

}

//func GetConversationMessages(db *sql.DB, userID int, receiverID int) ([]models.PrivateMessage, error) {
//	query := `
//SELECT id, sender_id, receiver_id, content, created_at
//FROM private_messages
//WHERE (sender_id = ? AND receiver_id = ?)
//   OR (sender_id = ? AND receiver_id = ?)
//ORDER BY created_at ASC
//`
//	rows, err := db.Query(query, userID, receiverID, receiverID, userID)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//
//	var messages []models.PrivateMessage
//
//	for rows.Next() {
//		var msg models.PrivateMessage
//		err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.Timestamp)
//		if err != nil {
//			return nil, err
//		}
//		messages = append(messages, msg)
//	}
//	return messages, nil
//}

// GetConversationMessages retrieves messages between two users
func GetConversationMessages(db *sql.DB, userID, otherUserID int) ([]models.PrivateMessage, error) {
	query := `
		SELECT id, sender_id, receiver_id, content, timestamp 
		FROM private_messages 
		WHERE (sender_id = ? AND receiver_id = ?) 
		   OR (sender_id = ? AND receiver_id = ?) 
		ORDER BY timestamp ASC
	`

	rows, err := db.Query(query, userID, otherUserID, otherUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.PrivateMessage
	for rows.Next() {
		var msg models.PrivateMessage
		err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.Timestamp)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

// SavePrivateMessage saves a private message to the database
func SavePrivateMessage(db *sql.DB, senderID, receiverID int, content string) error {
	query := `INSERT INTO private_messages (sender_id, receiver_id, content, timestamp) VALUES (?, ?, ?, ?)`
	_, err := db.Exec(query, senderID, receiverID, content, time.Now())
	return err
}

// GetUserMessages retrieves all messages sent to a specific user ID
func GetUserMessages(db *sql.DB, userID int) ([]models.PrivateMessage, error) {
	query := `
		SELECT id, sender_id, receiver_id, content, timestamp 
		FROM private_messages 
		WHERE receiver_id = ? 
		ORDER BY timestamp ASC
	`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.PrivateMessage
	for rows.Next() {
		var msg models.PrivateMessage
		//err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.Timestamp)
		//err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.Timestamp)
		err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.Timestamp)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

// GetAllConversations retrieves all conversations for the current user
func GetAllConversations(db *sql.DB, currentUserID int) ([]models.ConversationData, error) {
	// First, find all users that the current user has messaged with
	query := `
		SELECT DISTINCT 
			CASE 
				WHEN sender_id = ? THEN receiver_id 
				ELSE sender_id 
			END AS other_user_id
		FROM private_messages
		WHERE sender_id = ? OR receiver_id = ?
	`

	rows, err := db.Query(query, currentUserID, currentUserID, currentUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find conversation partners: %w", err)
	}
	defer rows.Close()

	// Collect all user IDs the current user has conversations with
	var otherUserIDs []int
	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			return nil, fmt.Errorf("failed to scan user ID: %w", err)
		}
		otherUserIDs = append(otherUserIDs, userID)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through users: %w", err)
	}

	// For each user, get their nickname and all messages
	var conversations []models.ConversationData
	for _, otherUserID := range otherUserIDs {
		// Get the other user's nickname
		var nickname string
		// Fixed: Use user_id instead of id in the query
		nicknameQuery := "SELECT nickname FROM users WHERE user_id = ?"
		err := db.QueryRow(nicknameQuery, otherUserID).Scan(&nickname)
		if err != nil {
			// If no user is found, use "Unknown User" instead of failing
			nickname = fmt.Sprintf("User #%d", otherUserID)
			fmt.Printf("Could not find nickname for user %d: %v\n", otherUserID, err)
			// Continue with the conversation anyway
		}

		// Get all messages between current user and other user
		messagesQuery := `
			SELECT id, sender_id, content, timestamp
			FROM private_messages
			WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)
			ORDER BY timestamp ASC
		`

		messageRows, err := db.Query(messagesQuery, currentUserID, otherUserID, otherUserID, currentUserID)
		if err != nil {
			return nil, fmt.Errorf("failed to get messages for conversation with user %d: %w", otherUserID, err)
		}

		var messages []models.ConversationMessage
		var lastUpdated time.Time

		for messageRows.Next() {
			var msg models.ConversationMessage
			var senderID int

			err := messageRows.Scan(&msg.ID, &senderID, &msg.Content, &msg.Timestamp)
			if err != nil {
				messageRows.Close()
				return nil, fmt.Errorf("failed to scan message: %w", err)
			}

			// Set the flag based on whether current user sent the message
			msg.IsFromCurrentUser = (senderID == currentUserID)

			messages = append(messages, msg)

			// Update the last timestamp if this message is newer
			if msg.Timestamp.After(lastUpdated) {
				lastUpdated = msg.Timestamp
			}
		}
		messageRows.Close()

		if err = messageRows.Err(); err != nil {
			return nil, fmt.Errorf("error iterating through messages: %w", err)
		}

		// Create the conversation entry
		conversation := models.ConversationData{
			UserID:      otherUserID,
			Nickname:    nickname,
			Messages:    messages,
			LastUpdated: lastUpdated,
		}

		conversations = append(conversations, conversation)
	}

	return conversations, nil
}
