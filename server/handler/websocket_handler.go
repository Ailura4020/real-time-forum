package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"real-time-forum/models"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	//"github.com/gorilla/websocket"
)

// Config de l'upgrader WebSocket pour accepter toutes les origines
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Structure d'un client connecté via WebSocket
type Client struct {
	Conn     *websocket.Conn
	Username string
	UserId   int
}

// Map pour stocker les clients connectés
// var clients = make(map[*Client]bool)
var clients = make(map[int]*Client)

// var pour la connexion la db
var db *sql.DB

func updateUserStatusInDB(userId int, status string) error {
	query := `UPDATE USERS SET Status = ? WHERE UserId = ?`
	_, err := db.Exec(query, status, userId)
	return err
}

// function gestion des connexions WebSocket
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Erreur lors de l'upgrade :", err)
		return
	}
	defer conn.Close()

	// Read the first message to get user information
	_, msg, err := conn.ReadMessage()
	if err != nil {
		fmt.Println("Erreur lors de la lecture du message", err)
		return
	}
	var userInfo struct {
		Username string `json:"username"`
		UserId   int    `json:"userId"`
	}
	if err := json.Unmarshal(msg, &userInfo); err != nil {
		fmt.Println("Erreur lors de la désérialisation du message", err)
		return
	}

	// Mise à jour de l'état de l'utilisateur dans la base de données
	if err := updateUserStatusInDB(userInfo.UserId, "online"); err != nil {
		fmt.Println("Erreur lors de la mise à jour du statut :", err)
	}

	fmt.Printf("User Info: %+v\n", userInfo)

	// création d'un nouveau client
	client := &Client{
		Conn:     conn,
		Username: userInfo.Username,
		UserId:   userInfo.UserId,
	}

	// ajout du client à la map
	clients[client.UserId] = client
	fmt.Println("client connecté", client.Username)
	fmt.Println(client.UserId)

	defer func() {
		if err := updateUserStatusInDB(client.UserId, "offline"); err != nil {
			fmt.Println("Erreur lors de la mise à jour du statut :", err)
		}
		delete(clients, client.UserId) // Retirer le client de la map
	}()

	// boucle lecture d'un message client
	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Erreur lors de la lecture du message", err)
			delete(clients, client.UserId)
			break
		}

		var message models.PrivateMessage
		if err := json.Unmarshal(msg, &message); err != nil {
			fmt.Println("Erreur lors de la désérialisation du message", err)
			continue
		}
		// enregistre le message dans la db
		if err := savePrivateMessage(db, client.UserId, message.ReceiverID, message.Content); err != nil {
			fmt.Println("Erreur lors de l'enregistrement du message", err)
		}
		// fmt.Printf("Message reçu : %s\n", msg)

		// envoie le message au destinataire
		if recipient, ok := clients[message.ReceiverID]; ok {
			err := recipient.Conn.WriteMessage(messageType, msg)
			if err != nil {
				fmt.Println("Erreur lors de l'envoi du message", err)
				recipient.Conn.Close()
				delete(clients, recipient.UserId)
			}
		}
	}
}

// function pour enregistrer un msg privé dans la DB
func savePrivateMessage(db *sql.DB, senderID int, receiverId int, content string) error {
	dateSent := time.Now().Format(time.RFC3339)
	query := `INSERT INTO PRIVATEMESSAGE (TextContent, DateSent, SenderId, ReceiverId) VALUE ( ?, ?, ?, ?)`
	_, err := db.Exec(query, content, dateSent, senderID, receiverId)
	return err
}

func getUserInfo(r *http.Request) (int, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return 0, err
	}
	userID, err := strconv.Atoi(cookie.Value)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
