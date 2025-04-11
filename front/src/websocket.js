function updateConnectedUsers(userContainer, users) {
    userContainer.innerHTML = '';
    users.forEach(user => {
      const userElement = document.createElement('div');
      userElement.className = 'user-item';
      userElement.textContent = user.nickname;
      userContainer.appendChild(userElement);
    });
    console.log('Connected users:', users);
  }
  


export function connectWebSocket(updateUserListCallback) {
    // connection à la websocket -- > new variable qui inclus NewConnectionWebsocket
    // src.onmessage 
    const socket = new WebSocket('ws://localhost:8080/ws');

    socket.onopen = () => {
        console.log('[Websocket] Connected');
    };

    socket.onmessage = (event) => {
        const data = JSON.parse(event.data);
        if (data.type === 'users') {
          const container = document.getElementById('connected-users');
          if (container) {
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
 