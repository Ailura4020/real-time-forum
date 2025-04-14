package models

import "time"

type PrivateMessage struct {
	ID         int       `json:"id"`
	SenderID   int       `json:"sender_id"`
	ReceiverID int       `json:"receiver_id"`
	Content    int       `json:"content"`
	Timestamp  time.Time `json:"timestamp"`
}

// type PrivateMessage struct {
// 	ID         int    `json:"id"`
// 	SenderID   int    `json:"sender_id"`
// 	ReceiverID int    `json:"receiver_id"`
// 	Content    string `json:"content"`
// 	DateSent   string `json:"date_sent"`
// }
