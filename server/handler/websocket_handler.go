package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"real-time-forum/models"
	"real-time-forum/repository"
	"real-time-forum/service"
	"real-time-forum/utils"
	"real-time-forum/middleware"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

// Config de l'upgrader WebSocket pour accepter toutes les origines
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocket(hub *utils.Hub, db *sql.DB) http.HandlerFunc {
	log.Println("HandleWebSocket")
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID and token from context (set by middleware)
		ctx := r.Context()
		userIDValue := ctx.Value(middleware.ContextUserID)
		tokenValue := ctx.Value(middleware.ContextToken)

		userID, ok := userIDValue.(int)
		tokenString, okToken := tokenValue.(string)
		if !ok || !okToken || userID == 0 || tokenString == "" {
			http.Error(w, "Unauthorized: No user info in context", http.StatusUnauthorized)
			return
		}

		fmt.Println("[HandleWebSocket] USER INFO", userID)

		// Check if user already has an active connection and disconnect it
		existingConnection := false
		hub.Mutex.Lock()
		for conn, user := range hub.Clients {
			if user.UserID == userID {
				existingConnection = true
				// Send a logout message to the existing connection
				conn.WriteJSON(map[string]interface{}{
					"type":    "system_message",
					"content": "You have been logged out because your account was accessed from another location",
					"action":  "force_logout",
				})
				// Close the existing connection
				conn.Close()
				// Remove the client from the hub
				delete(hub.Clients, conn)
				fmt.Printf("Disconnected existing session for user ID: %d\n", userID)
				break // Only need to disconnect one connection
			}
		}
		hub.Mutex.Unlock()

		// If there was an existing connection, blacklist the previous token
		if existingConnection {
			claims, _ := utils.ValidateJWT(tokenString)
			if claims != nil {
				// Get token expiry time from claims
				expiryTime := time.Unix(claims.ExpiresAt.Unix(), 0)
				utils.AddTokenToBlacklist(tokenString, expiryTime)
				fmt.Printf("Blacklisted token for user ID: %d\n", userID)
			}
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("Error during WebSocket upgrade:", err)
			http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
			return
		}
		defer conn.Close()

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

		// Start a goroutine to periodically clean up the token blacklist
		go func() {
			ticker := time.NewTicker(1 * time.Hour)
			defer ticker.Stop()
			for range ticker.C {
				utils.CleanupBlacklist()
			}
		}()

		// todo, add more messageType (keyboard_in, keyboard_out for )

		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				hub.RemoveClient(conn)
				hub.BroadcastUser()
				fmt.Println("WebSocket connection error:", err)
				break
			}

			fmt.Printf("Received message - Type: %d, Time: %s\n", messageType, time.Now().Format(time.RFC3339))
			fmt.Printf("Raw message content: %s\n", string(message))

			var msg map[string]interface{}
			if err := json.Unmarshal(message, &msg); err != nil {
				fmt.Println("Error parsing JSON message:", err)
				continue
			}
			fmt.Printf("Parsed message content: %+v\n", msg)

			if msg["type"] == "private_message" {
				fmt.Printf("Private message received at %s\n", time.Now().Format(time.RFC3339))
				fmt.Printf("Message details - From: %d, Content: %s\n", userID, msg["content"])

				senderID := userID

				// Récupérer 'to' et gérer les types possibles (string ou float64)
				var receiverID int
				switch v := msg["to"].(type) {
				case float64:
					receiverID = int(v) // Si "to" est un nombre
				case string:
					// Si "to" est une chaîne, tente de le convertir en nombre
					id, err := strconv.Atoi(v)
					if err != nil {
						fmt.Println("Erreur de conversion pour 'to' :", err)
						continue
					}
					receiverID = id
				default:
					fmt.Println("Type de 'to' inconnu :", v)
					continue
				}

				content := msg["content"].(string)

				err := repository.SavePrivateMessage(db, senderID, receiverID, content)
				if err != nil {
					fmt.Println("erreur insertion message privé :", err)
				}

				hub.Mutex.Lock()
				for clientConn, user := range hub.Clients {
					if user.UserID == receiverID {
						err := clientConn.WriteJSON(map[string]interface{}{
							"type":     "private_message",
							"from":     senderID,
							"content":  content,
							"datetime": time.Now().Format("2006-01-02 15:04:05"),
						})
						if err != nil {
							fmt.Println("Erreur envoi message au destinataire :", err)
						}
						notification := map[string]interface{}{
							"type":          "notification",
							"from":          senderID,
							"from_nickname": userInfo,
							"content":       "Vous avez reçu un nouveau message privé.",
							"timestamp":     time.Now().Format("2006-01-02 15:04:05"),
						}

						err = clientConn.WriteJSON(notification)
						if err != nil {
							fmt.Println("Erreur lors de l'envoi de la notification :", err)
						}

					}
				}
				hub.Mutex.Unlock()
			}
		}
	}
}
