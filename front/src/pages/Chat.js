import {api} from '../main.js';
import classes from '../styles/Chat.module.css';
import { chatTemplate } from "../templates.js"
import { connectWebSocket } from '../websocket.js';
import { getSocket } from '../websocket.js';


export async function renderChat(container) {
    // console.log("CONTAINER",container);
    console.log("CLASSES",classes)
    // Clear the container before rendering
    container.innerHTML = ''; // Clear the container
    const tempDiv = document.createElement('div');
    tempDiv.innerHTML = chatTemplate(classes); // Use the template with the classes
    container.appendChild(tempDiv);

    const token = localStorage.getItem('token');
    if (token) {
      connectWebSocket(token);
    }

    // Fetch user data
    const userData = await fetchUserData();

    // Display user ID and nickname if available
    const userInfoElement = container.querySelector('#user-info');
    if (userData) {
        userInfoElement.innerHTML = `
            <p>User ID: ${userData.data.id}</p>
            <p>Nickname: ${userData.data.nickname}</p>
        `;
    } else {
        userInfoElement.innerHTML = '<p>User not logged in or data could not be retrieved.</p>';
    }

    const sendButton = document.getElementById('send-button');
    sendButton.addEventListener('click', () => {
    console.log('[Chat] Tentative d\'envoi de message à', currentReceiver);

    const messageContent = document.getElementById('messageInput').value;
    if (currentReceiver && messageContent) {
        const message = {
            type: 'private_message',
            to: currentReceiver.id,
            content: messageContent
        };

        // Envoi du message via WebSocket
        console.log('[Chat] Message envoyé :', message);

        // socket.send(JSON.stringify(message));
       const socket = getSocket();
        if (socket && socket.readyState === WebSocket.OPEN) {
            socket.send(JSON.stringify(message));
        } else {
            console.log('[Chat] WebSocket n\'est pas ouvert.');
        }


        // Réinitialiser l'input
        document.getElementById('messageInput').value = '';
    }
});
}

// Function to fetch user data
async function fetchUserData() {
    const token = localStorage.getItem('token');
    if (!token) return null;

    try {
         // Adjust the endpoint as necessary
        return await api.get('/user');
    } catch (error) {
        console.error('Error fetching user data:', error);
        return null;
    }
}

export let currentReceiver = null;

export function setCurrentReceiver(user){
    currentReceiver = user;
    if (!currentReceiver) {
        console.warn('[Chat] Aucun destinataire sélectionné.');
        return;
      }

    const recipientElement = document.getElementById('chat-recipient');
    if (recipientElement){
        recipientElement.innerHTML = `Conversation avec : <strong>${user.nickname}</strong>`;
    }
    console.log("[Chat] conversation ouverte avec",user);
}


export function setupUserClickListener(){

    const userElements = document.querySelectorAll('.user-item');
    if (userElements.length === 0) {
        console.log("Aucun utilisateur trouvé.");
        return;  // Si aucun utilisateur, on s'arrête là
    }

    userElements.forEach(el => {
        el.addEventListener('click',() => {
            const userId = el.id
            const nickname = el.textContent;
            console.log('[Chat] Utilisateur sélectionné :', userId, nickname);
            setCurrentReceiver({id: userId, nickname})
        })
    });
}




