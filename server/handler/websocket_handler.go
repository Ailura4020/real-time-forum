package handler

import (
	"database/sql"
	"net/http"
	"real-time-forum/repository"
	"real-time-forum/service"
	"real-time-forum/utils"

	"github.com/gorilla/websocket"
)

// Config WebSocket sécurisé
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Autoriser temporairement toutes les origines (à restreindre en production)
	},
}

func HandleWebsocket(db *sql.DB, hub *utils.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	user, err := getUserInfo(r)
	// Ici, vous pouvez gérer la connexion WebSocket
	// Par exemple, lire et écrire des messages
}

func getUserInfo() {

}
