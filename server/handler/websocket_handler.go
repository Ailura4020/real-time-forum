// package handler

// import (
// 	"database/sql"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"real-time-forum/models"
// 	"time"

// 	"github.com/gorilla/websocket"
// )

// // Config de l'upgrader WebSocket pour accepter toutes les origines
// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

// // Structure d'un client connecté via WebSocket
// type Client struct {
// 	Conn     *websocket.Conn
// 	Username string
// 	UserId   int
// }

// // Map pour stocker les clients connectés
// // var clients = make(map[*Client]bool)
// var clients = make(map[int]*Client)

// // var pour la connexion la db
// var db *sql.DB

// func updateUserStatusInDB(userId int, status string) error {
// 	query := `UPDATE USERS SET Status = ? WHERE UserId = ?`
// 	_, err := db.Exec(query, status, userId)
// 	return err
// }

// // function gestion des connexions WebSocket
// func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		fmt.Println("Erreur lors de l'upgrade :", err)
// 		return
// 	}
// 	defer conn.Close()

// 	// Read the first message to get user information
// 	_, msg, err := conn.ReadMessage()
// 	if err != nil {
// 		fmt.Println("Erreur lors de la lecture du message", err)
// 		return
// 	}
// 	var userInfo struct {
// 		Username string `json:"username"`
// 		UserId   int    `json:"userId"`
// 	}
// 	if err := json.Unmarshal(msg, &userInfo); err != nil {
// 		fmt.Println("Erreur lors de la désérialisation du message", err)
// 		return
// 	}

// 	// Mise à jour de l'état de l'utilisateur dans la base de données
// 	if err := updateUserStatusInDB(userInfo.UserId, "online"); err != nil {
// 		fmt.Println("Erreur lors de la mise à jour du statut :", err)
// 	}

// 	fmt.Printf("User Info: %+v\n", userInfo)

// 	// création d'un nouveau client
// 	client := &Client{
// 		Conn:     conn,
// 		Username: userInfo.Username,
// 		UserId:   userInfo.UserId,
// 	}

// 	// ajout du client à la map
// 	clients[client.UserId] = client
// 	fmt.Println("client connecté", client.Username)
// 	fmt.Println()

// 	defer func() {
// 		if err := updateUserStatusInDB(client.UserId, "offline"); err != nil {
// 			fmt.Println("Erreur lors de la mise à jour du statut :", err)
// 		}
// 		delete(clients, client.UserId) // Retirer le client de la map
// 	}()

// 	// boucle lecture d'un message client
// 	for {
// 		messageType, msg, err := conn.ReadMessage()
// 		if err != nil {
// 			fmt.Println("Erreur lors de la lecture du message", err)
// 			delete(clients, client.UserId)
// 			break
// 		}

// 		var message models.PrivateMessage
// 		if err := json.Unmarshal(msg, &message); err != nil {
// 			fmt.Println("Erreur lors de la désérialisation du message", err)
// 			continue
// 		}
// 		// enregistre le message dans la db
// 		if err := savePrivateMessage(db, client.UserId, message.ReceiverID, message.Content); err != nil {
// 			fmt.Println("Erreur lors de l'enregistrement du message", err)
// 		}
// 		// fmt.Printf("Message reçu : %s\n", msg)

// 		// envoie le message au destinataire
// 		if recipient, ok := clients[message.ReceiverID]; ok {
// 			err := recipient.Conn.WriteMessage(messageType, msg)
// 			if err != nil {
// 				fmt.Println("Erreur lors de l'envoi du message", err)
// 				recipient.Conn.Close()
// 				delete(clients, recipient.UserId)
// 			}
// 		}
// 		// for c := range clients {
// 		// 	if c.UserId == message.ReceiverID {
// 		// 		err := c.Conn.WriteMessage(messageType, msg)
// 		// 		if err != nil {
// 		// 			fmt.Println("Erreur lors de l'envoi du message", err)
// 		// 			c.Conn.Close()
// 		// 			delete(clients, c)
// 		// 		}
// 		// 		break
// 		// 	}
// 		// }
// 	}
// }

// // function pour enregistrer un msg privé dans la DB
// func savePrivateMessage(db *sql.DB, senderID int, receiverId int, content string) error {
// 	dateSent := time.Now().Format(time.RFC3339)
// 	query := `INSERT INTO PRIVATEMESSAGE (TextContent, DateSent, SenderId, ReceiverId) VALUE ( ?, ?, ?, ?)`
// 	_, err := db.Exec(query, content, dateSent, senderID, receiverId)
// 	return err
// }

package handler

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Allow all connections
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close()

	for {
		// Read message from browser
		_, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("read error:", err)
			break
		}
		fmt.Printf("Received: %s\n", msg)

		// Write message back to browser
		if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
			fmt.Println("write error:", err)
			break
		}
	}
}

var clients = make(map[*websocket.Conn]bool) // Track active clients

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close()

	clients[ws] = true

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("read error:", err)
			delete(clients, ws)
			break
		}

		// Broadcast message to all clients
		for client := range clients {
			if err := client.WriteMessage(websocket.TextMessage, msg); err != nil {
				fmt.Println("broadcast error:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
	defer func() {
		delete(clients, ws)
		ws.Close()
	}()

	ws.SetPongHandler(func(appData string) error {
		fmt.Println("pong received")
		return nil
	})
}
