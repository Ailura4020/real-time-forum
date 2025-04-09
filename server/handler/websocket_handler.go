package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"real-time-forum/models"
	"real-time-forum/repository"
	"real-time-forum/service"
	"real-time-forum/utils"
	"strconv"

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

	userID, err := getUserInfo(r)
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Unauthorized"))
		return
	}
	user, err := userService.GetUserByID(userID)
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("User not found"))
		return
	}
	hub.AddClient(conn, &models.Userlist{
		UserID:   user.ID,
		Nickname: user.Nickname,
	})

	var connectedUsers []*models.Userlist
	hub.Mutex.Lock()
	for _, u := range hub.Clients {
		connectedUsers = append(connectedUsers, u)
	}
	hub.Mutex.Unlock()

	msg, _ := json.Marshal(struct {
		Type  string             `json:"type"`
		Users []*models.Userlist `json:"users"`
	}{
		Type:  "users_online",
		Users: connectedUsers,
	})
	hub.BroadcastMessage(msg)

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
