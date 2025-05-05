import { setupUserClickListener } from "./pages/Chat";


export let socket = null;

export function updateConnectedUsers(userContainer, users) {
  console.log('[updateConnectedUsers] Utilisateurs reçus :', users);
    userContainer.innerHTML = '';
    users.forEach(user => {
      console.log('Ajout de', user);
      const userElement = document.createElement('div');
      userElement.id = user.id;
      userElement.className = 'user-item';
      // userElement;addEventListener('click', () => {
      //   openPrivateChat(user);
      // })
      userElement.textContent = user.nickname;
      userContainer.appendChild(userElement);
    });
    setupUserClickListener();
  }

export function connectWebSocket(token) {
    console.log("TOKEN",token)
    // connection à la websocket -- > new variable qui inclus NewConnectionWebsocket
    console.log("Token utilisé pour la connexion:", token);
     socket = new WebSocket(`ws://localhost:8080/ws?token=${token}`);

    socket.onopen = () => {
        console.log('WebSocket connection established');
    };

    socket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      console.log('[WebSocket] Message reçu :', data);
    
      if (data.type === 'users') {
        const container = document.getElementById('connected-users');
        if (container) {
          updateConnectedUsers(container, data.users);
        }
      }
    
      if (data.type === 'private_message') {
        console.log('[WebSocket] Message privé reçu :', data);
        displayPrivateMessage({
          senderId: data.from,
          content: data.content,
          datetime: data.datetime,
        });
      }
    };
    
    socket.onerror = (error => {
        console.log('[Websocket] Error:', error);
    });
    socket.onclose = (event => {
    socket.close();
        console.log('[Websocket] Connection closed:', event);
    })
    // createWebSocketConnection();
}

// Fonction pour retirer un utilisateur de la liste
function removeUserFromList(userID) {
  const userElement = document.getElementById(userID);
  if (userElement) {
    userElement.remove();
  }
}

export function displayPrivateMessage(message) {
  const messagesContainer = document.getElementById('messages');
  if (messagesContainer) {
      const messageElement = document.createElement('div');
      messageElement.className = 'message';
      messageElement.innerHTML = `
          <strong>${message.fromNickname}</strong> : ${message.content} <br>
          <small>${message.datetime}</small>
      `;
      messagesContainer.appendChild(messageElement);
      
      // Scroll vers le bas à chaque nouveau message
      messagesContainer.scrollTop = messagesContainer.scrollHeight;
  }
}

export function getSocket() {
  return socket;
}