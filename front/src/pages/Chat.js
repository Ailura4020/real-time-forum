import {api} from '../main.js';
// src/pages/About.js
// import '../styles/About.module.css';

export async function renderChat(container) {
    // Clear the container before rendering
    container.innerHTML = `
     <div class="chat">
      <h1>Chat</h1>
     <div class="chat-layout">
  <div class="chat-sidebar">
    <h3>Utilisateurs connectés</h3>
    <div id="connected-users" class="user-list"></div>
  </div>
  <div class="chat-content">
          <div class="chat-messages" id="messages"></div>
          <div class="chat-input">
            <input type="text" id="message-input" placeholder="Votre message...">
            <button id="send-button">Envoyer</button>
          </div>
        </div>
      </div>
      <div id="user-info"></div>
    </div>
  `;

    // Quelques styles inline pour que ce soit visible rapidement
    const style = document.createElement('style');
    style.textContent = `
      .chat-layout {
        display: flex;
        border: 1px solid #ccc;
        height: 400px;
      }
      .chat-sidebar {
        width: 200px;
        border-right: 1px solid #ccc;
        padding: 10px;
        overflow-y: auto;
        background: #f0f0f0;
      }
      #connected-users {
        border: 1px solid #aaa;
        padding: 10px;
        background: white;
        height: 300px;
        overflow-y: auto;
      }
      .user-item {
        padding: 5px;
        margin-bottom: 5px;
        background: #e8e8e8;
        border-radius: 4px;
        cursor: pointer;
      }
    `;
    document.head.appendChild(style);

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
