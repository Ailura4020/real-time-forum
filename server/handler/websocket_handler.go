package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"real-time-forum/models"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	Conn     *websocket.Conn
	Username string
	UserId   int
}

var clients = make(map[*Client]bool)
var db *sql.DB

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Erreur lors de l'upgrade :", err)
		return
	}
	defer conn.Close()

	client := &Client{
		Conn:     conn,
		Username: "user1",
		UserId:   1,
	}

	clients[client] = true
	fmt.Println("client connecté", client.Username)

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Erreur lors de la lecture du message", err)
			delete(clients, client)
			break
		}

		var message models.PrivateMessage
		if err := json.Unmarshal(msg, &message); err != nil {
			fmt.Println("Erreur lors de la désérialisation du message", err)
			continue
		}

		if err := savePrivateMessage(db, client.UserId, message.ReceiverID, message.Content); err != nil {
			fmt.Println("Erreur lors de l'enregistrement du message", err)
		}

		// fmt.Printf("Message reçu : %s\n", msg)

		for c := range clients {
			if c.UserId == message.ReceiverID {
				err := c.Conn.WriteMessage(messageType, msg)
				if err != nil {
					fmt.Println("Erreur lors de l'envoi du message", err)
					c.Conn.Close()
					delete(clients, c)
				}
				break
			}
		}
	}
}

func savePrivateMessage(db *sql.DB, senderID int, receiverId int, content string) error {
	dateSent := time.Now().Format(time.RFC3339)
	query := `INSERT INTO PRIVATEMESSAGE (TextContent, DateSent, SenderId, ReceiverId) VALUE ( ?, ?, ?, ?)`
	_, err := db.Exec(query, content, dateSent, senderID, receiverId)
	return err
}
