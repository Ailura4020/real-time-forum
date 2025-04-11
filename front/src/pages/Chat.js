import {api} from '../main.js';
// src/pages/About.js
// import '../styles/About.module.css';

export async function renderChat(container) {
    // Clear the container before rendering
    container.innerHTML = `
    <div class="chat">
      <h1>Chat</h1>
      <div class="chat-content"></div>
      <div id="user-info"></div> <!-- Placeholder for user info -->
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
function updateConnectedUsers(userContainer, usersList){
  userContainer.innerHTML='';
  usersList.forEach(user => {
    const userElement = document.createElement('div');
    userElement.className = 'user-item';
    userElement.textContent = user.nickname;
    userContainer.appendChild(userElement);
});
console.log('Connected users:', usersList);
}