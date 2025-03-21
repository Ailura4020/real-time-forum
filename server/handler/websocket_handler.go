package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

// Config WebSocket sécurisé
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Autoriser temporairement toutes les origines (à restreindre en production)
	},
}

// Map des clients WebSocket
var (
	clients      = make(map[int]*Client)
	clientsMutex = sync.Mutex{}
)

// Structure d'un client connecté via WebSocket
type Client struct {
	Conn     *websocket.Conn
	Username string
	UserId   int
}

// Fonction de gestion des connexions WebSocket
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	log.Println("📡 Tentative de connexion WebSocket reçue...")

	log.Println("🔍 Headers reçus pour WebSocket :")
	for name, values := range r.Header {
		for _, value := range values {
			log.Printf("%s: %s", name, value)
		}
	}

	// Vérifier que la requête est bien une WebSocket
	if r.Header.Get("Upgrade") != "websocket" {
		log.Println("❌ Erreur: Le header Upgrade n'est pas 'websocket'")
		http.Error(w, "Invalid WebSocket request", http.StatusBadRequest)
		return
	}
	if !strings.Contains(strings.ToLower(r.Header.Get("Connection")), "upgrade") {
		log.Println("❌ Erreur: Le header Connection ne contient pas 'Upgrade'")
		http.Error(w, "Invalid WebSocket request", http.StatusBadRequest)
		return
	}

	// Connexion WebSocket établie
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("❌ Erreur lors de l'upgrade WebSocket :", err)
		return
	}
	defer conn.Close()

	log.Println("✅ Connexion WebSocket réussie !")

	// Lire le premier message pour obtenir les infos utilisateur
	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println("❌ Erreur lors de la lecture du message initial :", err)
		return
	}
	log.Println("📩 Message reçu à la connexion:", string(msg))

	var userInfo struct {
		Username string `json:"username"`
		UserId   int    `json:"userId"`
	}

	log.Println("🔍 Tentative de parsing des infos utilisateur...")
	if err := json.Unmarshal(msg, &userInfo); err != nil {
		log.Printf("🔍 Données utilisateur reçues: %+v\n", userInfo)
		log.Println("❌ Erreur lors du parsing du message initial :", err)
		log.Println("🔴 Contenu du message reçu:", string(msg))
		return
	}

	log.Printf("✅ Utilisateur connecté: %s (ID: %d)\n", userInfo.Username, userInfo.UserId)

	client := &Client{
		Conn:     conn,
		Username: userInfo.Username,
		UserId:   userInfo.UserId,
	}

	// Ajouter le client en verrouillant l'accès
	clientsMutex.Lock()
	clients[client.UserId] = client
	clientsMutex.Unlock()

	log.Printf("✅ Client connecté: %s (ID: %d)\n", client.Username, client.UserId)

	// Rester en écoute des messages
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("❌ Erreur de lecture du message :", err)
			break
		}
		log.Printf("📩 Message reçu : %s\n", msg)
	}
}
