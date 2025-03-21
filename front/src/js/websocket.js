import { fetchChats } from "../api/usersChatsAPI.js";
import {
  appendChatMessage,
  changeChat,
  createChats,
  getLastMessageID,
  messagesCount,
  updateUserStatus,
} from "../components/chat.js";
import { isLoggedIn } from "../utils/auth.js";
import { config } from "../config/config.js";

let ws;
let reconnectAttempts = 0;
const MAX_RECONNECT_ATTEMPTS = 5;
const RECONNECT_INTERVAL = 3000;

export let currentUser;

export async function connectWebSocket() {
  if (!window["WebSocket"]) {
    alert("❌ WebSockets ne sont pas supportés dans ce navigateur.");
    return;
  }

  try {
    const protocol = window.location.protocol === "https:" ? "wss" : "ws";
const ws = new WebSocket(`${protocol}://${window.location.hostname}:8080/ws`);

    const url = `${protocol}://${window.location.hostname}:${config.wsPort}/ws`;
    ws = new WebSocket(url);

    ws.onopen = async function () {
        console.log("✅ WebSocket connecté !");
        currentUser = await isLoggedIn();
        if (currentUser) {
          console.log("📤 Envoi des infos utilisateur :", JSON.stringify({
            username: currentUser?.username || "inconnu",
            userId: currentUser?.id || 0
        }));
        
          ws.send(
            JSON.stringify({
              username: currentUser.username,
              userId: currentUser.id,
            })
          );
          console.log("📤 Infos utilisateur envoyées:", currentUser);
        }
      };

    ws.onmessage = function (message) {
      try {
        const data = JSON.parse(message.data);
        routeEvent(data);
      } catch (error) {
        console.error("❌ Erreur lors du parsing du message WebSocket:", error);
      }
    };

    ws.onerror = function (error) {
      console.error("❌ WebSocket error:", error);
    };

    ws.onclose = function (event) {
      console.warn("⚠️ WebSocket fermé :", event);
      attemptReconnect();
    };
  } catch (error) {
    console.error("❌ Impossible de se connecter à WebSocket.");
    alert("Connexion WebSocket perdue. Rechargez la page.");
    attemptReconnect();
  }
}

function attemptReconnect() {
  if (reconnectAttempts < MAX_RECONNECT_ATTEMPTS) {
    reconnectAttempts++;
    console.log(`🔄 Tentative de reconnexion #${reconnectAttempts}...`);
    setTimeout(connectWebSocket, RECONNECT_INTERVAL);
  } else {
    console.error("❌ Impossible de reconnecter WebSocket après plusieurs tentatives.");
  }
}

export function sendEvent(eventType, payload) {
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify({ type: eventType, payload }));
  } else {
    console.warn("⚠️ Impossible d'envoyer le message, WebSocket fermé.");
  }
}

async function routeEvent(msg) {
  switch (msg.type) {
    case "new_message":
      getLastMessageID(msg.payload.id);
      appendChatMessage(msg.payload);
      break;
    case "past_messages":
      const messages = msg.payload;
      if (messages.messages) {
        messages.messages.reverse().forEach((message) => appendChatMessage(message, true));
      }
      if (messages.count_of_messages > 0) {
        messagesCount(messages.count_of_messages);
      }
      break;
    case "chat_list_update":
      createChats(await fetchChats());
      break;
    case "change_chat":
      changeChat(msg.payload);
      break;
    case "status_update":
      if (msg.payload) {
        msg.payload.forEach(({ username, status, last_seen }) => updateUserStatus(username, status, last_seen));
      }
      break;
    default:
      console.warn("⚠️ Type de message WebSocket non supporté :", msg.type);
  }
}

export function closeWebSocket() {
  if (ws) {
    ws.close();
  }
}

window.addEventListener("load", connectWebSocket);
window.addEventListener("beforeunload", () => {
  sendEvent("logout", { userId: currentUser?.id });
});
