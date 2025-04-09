package utils

import (
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
	hub.Clients[conn] = user
}

func (hub *Hub) RemoveClient(conn *websocket.Conn) {
	hub.Mutex.Lock()
	defer hub.Mutex.Unlock()
	delete(hub.Clients, conn)
}

func (hub *Hub) BroadcastMessage(message []byte) {
	hub.Mutex.Lock()
	defer hub.Mutex.Unlock()

	var onlineUser []*models.Userlist
	for _, user := range hub.Clients {
		onlineUser = append(onlineUser, user)
	}

	for conn := range hub.Clients {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			conn.Close()
			delete(hub.Clients, conn)
		}
	}
}
