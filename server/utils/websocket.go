package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"real-time-forum/models"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	Clients map[*websocket.Conn]*models.Userlist
	Mutex   sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Clients: make(map[*websocket.Conn]*models.Userlist),
	}
}

func (hub *Hub) AddClient(conn *websocket.Conn, user *models.Userlist) {
	hub.Mutex.Lock()
	defer hub.Mutex.Unlock()
	// Log avant d'ajouter
	fmt.Println("Before adding new client, current clients:", len(hub.Clients))
	for c, u := range hub.Clients {
		fmt.Printf("  - Client %p: %s (ID: %d)\n", c, u.Nickname, u.UserID)
	}
	hub.Clients[conn] = user
	// Log après avoir ajouté
	fmt.Println("After adding new client, current clients:", len(hub.Clients))
	for c, u := range hub.Clients {
		fmt.Printf("  - Client %p: %s (ID: %d)\n", c, u.Nickname, u.UserID)
	}
}

func (hub *Hub) RemoveClient(conn *websocket.Conn) {
	hub.Mutex.Lock()
	defer hub.Mutex.Unlock()
	delete(hub.Clients, conn)
	// Récupérer l'ID de l'utilisateur (ou un autre identifiant) associé à cette connexion
	userID := conn.RemoteAddr().String() // ou l'identifiant spécifique de l'utilisateur

	// Affichage de la nouvelle liste des clients restants
	fmt.Println("Clients restants connectés :")
	for c, u := range hub.Clients {
		fmt.Printf("  - Client %p : %s (ID: %d)\n", c, u.Nickname, u.UserID)
	}

	fmt.Printf("Nombre total de clients connectés : %d\n", len(hub.Clients))

	// Créer un message indiquant que l'utilisateur s'est déconnecté
	message := map[string]interface{}{
		"type":   "user_disconnect",
		"userID": userID,
	}

	// Envoyer un message de déconnexion à tous les autres clients
	for client := range hub.Clients {
		if client != conn {
			err := client.WriteJSON(message) // Envoi du message à chaque client
			if err != nil {
				log.Println("Erreur lors de l'envoi du message de déconnexion:", err)
			}
		}
	}
}

func (hub *Hub) BroadcastUser() {
	hub.Mutex.Lock()
	defer hub.Mutex.Unlock()

	var onlineUser []*models.Userlist
	for _, user := range hub.Clients {
		onlineUser = append(onlineUser, user)
	}
	// fmt.Println("Online users:", onlineUser)
	for i, user := range onlineUser {
		fmt.Printf("User %d: ID=%d, Nickname=%s\n", i, user.UserID, user.Nickname)
	}

	message, err := json.Marshal(struct {
		Type  string             `json:"type"`
		Users []*models.Userlist `json:"users"`
	}{
		Type:  "users",
		Users: onlineUser,
	})
	if err != nil {
		return
	}

	fmt.Println("Broadcasting user list:", string(message))

	for conn := range hub.Clients {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			conn.Close()
			fmt.Print("Users disconnected")
			delete(hub.Clients, conn)
		}
	}
}

// ajouter dans boadcastUser fonction d'envoyer msg
