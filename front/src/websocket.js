
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
      // console.log("event", event.data);    
      const data = JSON.parse(event.data);
      console.log('[WebSocket] Message reçu :', data);
      if (data.type === 'users') {
        const container = document.getElementById('connected-users');
        console.log('[WebSocket] Container found:', container);
        if (container) {
            // console.log("ou est la listes ",data);
            updateConnectedUsers(container, data.users);
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
