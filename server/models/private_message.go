package models

import "time"

type PrivateMessage struct {
	ID         int       `json:"id"`
	SenderID   int       `json:"sender_id"`
	ReceiverID int       `json:"receiver_id"`
	Content    string    `json:"content"`
	Timestamp  time.Time `json:"timestamp"`
}

// ConversationData represents a single conversation with a user
type ConversationData struct {
	UserID      int                   `json:"user_id"`
	Nickname    string                `json:"nickname"`
	Messages    []ConversationMessage `json:"messages"`
	LastUpdated time.Time             `json:"last_updated"`
}

// ConversationMessage represents a message in a conversation
type ConversationMessage struct {
	ID                int       `json:"id"`
	Content           string    `json:"content"`
	Timestamp         time.Time `json:"timestamp"`
	IsFromCurrentUser bool      `json:"is_from_current_user"`
}
