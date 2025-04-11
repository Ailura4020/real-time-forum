
export function updateConnectedUsers(userContainer, users) {
    userContainer.innerHTML = '';
    users.forEach(user => {
      const userElement = document.createElement('div');
      userElement.className = 'user-item';
      userElement.textContent = user.nickname;
      userContainer.appendChild(userElement);
    });
  }

export function connectWebSocket(token) {
    // connection à la websocket -- > new variable qui inclus NewConnectionWebsocket
    // src.onmessage
    // const token = localStorage.getItem('token');
    // const socket = new WebSocket('ws://localhost:8080/ws', null, {
    //     headers: {
    //         'Authorization': `Bearer ${token}`,
    //         'Content-Type': 'application/json'
    //     }
    // });
    console.log("Token utilisé pour la connexion:", token);
    const socket = new WebSocket(`ws://localhost:8080/ws?token=${token}`);

    socket.onopen = () => {
        console.log('WebSocket connection established');
    };

    socket.onmessage = (event) => {
      // console.log("event", event.data);    
      const data = JSON.parse(event.data);
      if (data.type === 'users') {
        const container = document.getElementById('connected-users');
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
