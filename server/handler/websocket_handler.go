package handler

import (
	"fmt"
	"net/http"

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
}

var clients = make(map[*Client]bool)

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Erreur lors de l'upgrade :", err)
		return
	}
	defer conn.Close()

	client := &Client{Conn: conn}
	clients[client] = true
	fmt.Println("client connecté")

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Erreur lors de la lecture du message", err)
			delete(clients, client)
			break
		}
		fmt.Printf("Message reçu : %s\n", msg)

		for c := range clients {
			if err := c.Conn.WriteMessage(messageType, msg); err != nil {
				fmt.Println("Erreur lors de l'envoi du message", err)
				c.Conn.Close()
				delete(clients, c)
			}
		}
	}
}
