export function connectWebSocket (updateUserListCallback) {
// connection à la websocket -- > new variable qui inclus NewConnectionWebsocket
// src.onmessage 
let socket = null;
const wsurl = 'ws://localhost:8000/ws'; // Replace with your WebSocket URL
socket = new WebSocket(wsurl);

socket.onopen = () => {
    console.log('WebSocket connection established');
};

socket.onmessage = (event) => {
    try {
        const data = JSON.parse(event.data);
        if (Array.isArray(data)) {
            updateUserListCallback(data);
        }else{
            console.log('[Websocket] Received message:', data);
        }
     } catch (e) {
        console.error('[Websocket] Failed to parse message', e);
    }
};
socket.onerror = (error => {
    console.log('[Websocket] Error:', error);
});
socket.onclose = (event => {
console.log('[Websocket] Connection closed:', event);
})
}
