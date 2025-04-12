
export function updateConnectedUsers(userContainer, users) {
  console.log('[updateConnectedUsers] Utilisateurs reçus :', users);
    userContainer.innerHTML = '';
    users.forEach(user => {
      console.log('Ajout de', user);
      const userElement = document.createElement('div');
      userElement.className = 'user-item';
      userElement.textContent = user.nickname;
      userContainer.appendChild(userElement);
    });
  }

export function connectWebSocket(token) {
    // connection à la websocket -- > new variable qui inclus NewConnectionWebsocket
    console.log("Token utilisé pour la connexion:", token);
    const socket = new WebSocket(`ws://localhost:8080/ws?token=${token}`);

    socket.onopen = () => {
        console.log('WebSocket connection established');
    };

    socket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      console.log('[WebSocket] Message reçu :', data);
    
      if (data.type === 'users') {
        const container = document.getElementById('connected-users');
        console.log('[WebSocket] Container trouvé :', container);
        if (container) {
          updateConnectedUsers(container, data.users);
        }
      }
    
      if (data.type === 'user_disconnect') {
        console.log('[WebSocket] Déconnexion reçue pour l\'utilisateur:', data.userID);
        const container = document.getElementById('connected-users');
        if (container) {
          removeUserFromList(data.userID);
        }
      }
    
      // Gérer les déconnexions d'utilisateur
      if (data.type === 'user_disconnect') {
        const container = document.getElementById('connected-users');
        console.log('[WebSocket] Container trouvé pour la déconnexion :', container);
        if (container) {
          // Retirer l'utilisateur déconnecté de la liste
          removeUserFromList(data.userID);
        }
      }
    };
    socket.onerror = (error => {
        console.log('[Websocket] Error:', error);
    });
    socket.onclose = (event => {
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