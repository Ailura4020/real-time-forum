package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"real-time-forum/models"
	"real-time-forum/repository"
	"real-time-forum/service"
	"real-time-forum/utils"
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

// var pour la connexion la db
// var db *sql.DB

// function gestion des connexions WebSocket
//
//	func HandleWebSocket(hub *utils.Hub) http.HandlerFunc {
//		return func(w http.ResponseWriter, r *http.Request) {
//			conn, err := upgrader.Upgrade(w, r, nil)
//			if err != nil {
//				fmt.Println("Erreur lors de l'upgrade :", err)
//				return
//			}
//			defer conn.Close()
//
//			// Read the first message to get user information
//			UserRepo := repository.NewUserRepository(db)
//			UserService := service.NewUserService(UserRepo)
//			// Mise à jour de l'état de l'utilisateur dans la base de données
//			userID, err := GetUserInfo(r)
//			if err != nil {
//				fmt.Println("Erreur lors de la récupération des informations utilisateur :", err)
//			}
//			userInfo, err := UserService.GetNickname(userID)
//			if err != nil {
//				fmt.Println("Erreur lors de la mise à jour du statut :", err)
//			}
//
//			var userList models.Userlist
//			userList.UserID = userID
//			userList.Nickname = userInfo
//			fmt.Println("USERS:", userList)
//			hub.AddClient(conn, &userList)
//			hub.BroadcastUser()
//
//			for {
//				_, _, err := conn.ReadMessage()
//				if err != nil {
//					break
//				}
//			}
//			hub.RemoveClient(conn)
//			hub.BroadcastUser()
//		}
//	}

func HandleWebSocket(hub *utils.Hub, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract user ID from JWT token
		userID, err := utils.ExtractUserIDFromRequest(r)
		if err != nil {
			fmt.Println("Erreur lors de la récupération des informations utilisateur :", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		fmt.Println("USER INFOR----------------------------------", userID)

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("Erreur lors de l'upgrade :", err)
			return
		}
		defer conn.Close()
		// hub := utils.NewHub()
		UserRepo := repository.NewUserRepository(db)
		UserService := service.NewUserService(UserRepo)

		userInfo, err := UserService.GetNickname(userID)
		if err != nil {
			fmt.Println("Erreur lors de la mise à jour du statut :", err)
		}

		var userList models.Userlist
		userList.UserID = userID
		userList.Nickname = userInfo
		fmt.Println("USERS:", userList)
		hub.AddClient(conn, &userList)
		hub.BroadcastUser()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
		hub.RemoveClient(conn)
		hub.BroadcastUser()
	}
}

// function pour enregistrer un msg privé dans la DB
func savePrivateMessage(db *sql.DB, senderID int, receiverId int, content string) error {
	dateSent := time.Now().Format(time.RFC3339)
	query := `INSERT INTO PRIVATEMESSAGE (TextContent, DateSent, SenderId, ReceiverId) VALUE ( ?, ?, ?, ?)`
	_, err := db.Exec(query, content, dateSent, senderID, receiverId)
	return err
}

//func GetUserInfo(r *http.Request) (int, error) {
//	cookie, err := r.Cookie("session_id")
//	if err != nil {
//		return 0, err
//	}
//	userID, err := strconv.Atoi(cookie.Value)
//	if err != nil {
//		return 0, err
//	}
//	return userID, nil
//}
