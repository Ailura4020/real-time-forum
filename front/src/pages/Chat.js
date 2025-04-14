import {api} from '../main.js';
// import classes from '../styles/Chat.module.css';

export async function renderChat(container) {
    // Clear the container before rendering
    container.innerHTML = `
<div id="user-info" class="${classes.userInfo}"></div>
 <div class="chat">
   <h1>Chat</h1>
   <div class="${classes.chatLayout}">
      <div class="${classes.chatSidebar}">
         <h3>Utilisateurs connectés</h3>
         <div id="connected-users" class="${classes.usersList}"></div>
      </div>
      <div class="${classes.chatContent}">
         <div class="${classes.chatMessages}" id="messages"></div>
         <div class="${classes.chatInput}">
            <input type="text" id="messageInput" placeholder="Votre message...">
            <button id="send-button" class="${classes.sendButton}">Envoyer</button>
         </div>
      </div>
   </div>
</div>
  `;

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
