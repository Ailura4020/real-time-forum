package utils

import (
	"encoding/json"
	"fmt"
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
			delete(hub.Clients, conn)
		}
	}
}
