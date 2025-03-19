// src/main.js
import { router } from './router.js';
import { renderTemplate } from './templates.js';
import LoginModal from './components/Login.js';
import RegisterModal from './components/Registration.js';

let currentUser = null;
let socket = null;

function render() {
    const app = document.getElementById('app');
    const page = router();
    app.innerHTML = renderTemplate(page);
}

function showModal(modalContent) {
    const app = document.getElementById('app');
    app.innerHTML += modalContent;
}

function closeModal() {
    const modal = document.querySelector('.modal');
    if (modal) {
        modal.remove();
    }
}

window.addEventListener('popstate', render);
document.addEventListener('click', (event) => {
    if (event.target.tagName === 'A') {
        event.preventDefault();
        const path = event.target.getAttribute('href');
        window.history.pushState({}, '', path);
        render();
    } else if (event.target.id === 'openLoginModal') {
        showModal(LoginModal({ onLogin: handleLogin, onClose: closeModal }));
    } else if (event.target.id === 'openRegisterModal') {
        showModal(RegisterModal({ onRegister: handleRegister, onClose: closeModal }));
    }
});

// Handle login
async function handleLogin(event) {
    event.preventDefault();
    const formData = new FormData(event.target);
    const data = {
        email: formData.get('identifier'),
        password: formData.get('password')
    };

    try {
        const response = await fetch('http://localhost:8080/api/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        });

        const result = await response.json
        if (!response.ok) {
            throw new Error(result.message || `HTTP error! status: ${response.status}`);
        }

        // Handle successful login
        closeModal(); // Close the modal
        displayUserInfo(result.data); // Display user info
        currentUser = result.data; // Set current user in memory
        initializeWebSocket(result.token); // Initialize WebSocket with the token
    } catch (error) {
        console.error('Error during login:', error);
        const loginMessage = document.getElementById('loginMessage');
        loginMessage.textContent = error.message || 'Login failed. Please try again.';
        loginMessage.style.color = 'red';
    }
}

// Handle registration
async function handleRegister(event) {
    event.preventDefault();
    const formData = new FormData(event.target);
    const data = {
        nickname: formData.get('nickname'),
        age: parseInt(formData.get('age')),
        gender: formData.get('gender'),
        first_name: formData.get('first_name'),
        last_name: formData.get('last_name'),
        email: formData.get('email'),
        password: formData.get('password'),
    };

    try {
        const response = await fetch('http://localhost:8080/api/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data)
        });

        const result = await response.json();

        if (!response.ok) {
            throw new Error(result.message || `HTTP error! status: ${response.status}`);
        }

        // Handle successful registration
        closeModal(); // Close the modal
        displayUserInfo(result.data); // Display user info
        currentUser = result.data; // Set current user in memory
        initializeWebSocket(result.token); // Initialize WebSocket with the token
    } catch (error) {
        console.error('Error during registration:', error);
        const registerMessage = document.getElementById('registerMessage');
        registerMessage.textContent = error.message || 'Registration failed. Please try again.';
        registerMessage.style.color = 'red';
    }
}

// Display user information
function displayUserInfo(userData) {
    const userInfoContainer = document.getElementById('userInfoContainer');
    userInfoContainer.innerHTML = `
        <p>Welcome, ${userData.nickname}!</p>
        <p>Email: ${userData.email}</p>
        <p>Age: ${userData.age}</p>
        <p>Gender: ${userData.gender}</p>
        <button id="logoutButton">Logout</button>
    `;

    // Add logout button functionality
    document.getElementById('logoutButton').addEventListener('click', logout);
}

// Handle logout
function logout() {
    if (socket) {
        socket.close();
    }
    localStorage.removeItem('token');
    currentUser = null;
    document.getElementById('userInfoContainer').innerHTML = '';
    showAuthInterface(); // Show login/register forms
}

// Show authentication interface
function showAuthInterface() {
    // Logic to show login/register forms
}

// Initialize WebSocket connection
function initializeWebSocket(token) {
    // Logic to initialize WebSocket connection
}

// Helper to get token
const getToken = () => {
    return localStorage.getItem('token');
};

// Initial render
window.addEventListener('popstate', render);
document.addEventListener('DOMContentLoaded', () => {
    render();
    const token = getToken();
    if (token) {
        checkUserStatus(token); // Check user status if token exists
    } else {
        showAuthInterface(); // Show login/register if no token
    }
});


// Initial render
render();