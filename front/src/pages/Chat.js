import {api} from '../main.js';
import classes from '../styles/Chat.module.css';
import { chatTemplate } from "../templates.js"
import { connectWebSocket } from '../websocket.js';


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

let currentReceiver = null;

function setCurrentReceiver(user){
    currentReceiver = user;

    const recipientElement = document.getElementById('chat-recipient');
    if (recipientElement){
        recipientElement.innerHTML = `Conversation avec : <strong>${user.nickname}</strong>`;
    }
    console.log("[Chat] conversation ouverte avec",user);
}

export function setupUserClickListener(){
    const userElements = document.querySelectorAll('.user-item');
    userElements.forEach(el => {
        el.addEventListener('click',() => {
            const userId = el.id
            const nickname = el.textContent;

            setCurrentReceiver({id: userId, nickname})
        })
    })
}