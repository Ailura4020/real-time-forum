package repository

import (
	"database/sql"
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

		err := rows.Scan(&msg.ID, &msg.SenderID, msg.ReceiverID, &msg.Content, &msg.Timestamp)
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
		err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.Timestamp)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
