import { router } from './router.js';
import { renderNavigation } from './components/Navigation.js';
import { connectWebSocket } from './websocket.js';
// import './styles/Main.module.css';

document.addEventListener('DOMContentLoaded', async () => {
    const appElement = document.getElementById('app');

    // Render navigation bar
    const navElement = renderNavigation();
    appElement.appendChild(navElement);

    // Ajoute le conteneur de toast notifications
    if (!document.getElementById('toast-container')) {
        const toastContainer = document.createElement('div');
        toastContainer.id = 'toast-container';
        toastContainer.style.position = 'fixed';
        toastContainer.style.top = '2rem';
        toastContainer.style.right = '2rem';
        toastContainer.style.zIndex = '2000';
        appElement.appendChild(toastContainer);
    }

    // Fonction globale pour afficher un toast
    // unused (for alert boxes)
    // window.showToast = function(message) {
    //     const toastContainer = document.getElementById('toast-container');
    //     if (!toastContainer) return;
    //     const toast = document.createElement('div');
    //     toast.className = 'toast-notif-red';
    //     toast.innerHTML = `<span style="font-weight:bold;">${message}</span>`;
    //     toast.style.background = 'rgba(255,0,0,0.95)';
    //     toast.style.color = '#fff';
    //     toast.style.padding = '1rem 2rem';
    //     toast.style.marginBottom = '1rem';
    //     toast.style.borderRadius = '8px';
    //     toast.style.boxShadow = '0 0 16px 4px #ff0000, 0 0 32px 8px #ff0000';
    //     toast.style.fontSize = '1.1rem';
    //     toast.style.letterSpacing = '1px';
    //     toast.style.display = 'flex';
    //     toast.style.alignItems = 'center';
    //     toast.style.animation = 'toastFadeIn 0.3s';
    //     toastContainer.appendChild(toast);
    //     // Lecture du son
    //     const audio = new Audio('/static/sounds/lightsaber3.mp3');
    //     audio.play();
    //     // Disparition auto
    //     setTimeout(() => {
    //         toast.style.animation = 'toastFadeOut 0.5s';
    //         setTimeout(() => toast.remove(), 500);
    //     }, 3500);
    // };

    // Create main content container
    const mainContent = document.createElement('main');
    mainContent.id = 'main-content';
    appElement.appendChild(mainContent);

    // Check and validate token first
    const token = localStorage.getItem('token');
    if (token) {
        const isValid = await validateToken();

        if (isValid) {
            // Token is valid, update UI with user data
            const userData = await api.get('/user');
            updateAuthUI(userData);
        } else {
            // Token was invalid and has been cleared
            updateAuthUI();
        }
    }

    // Initialize router
    router.init();
});

/**
 * Validates the authentication token stored in localStorage.
 * Attempts to verify the token by making an API request to fetch user data.
 * If the token is invalid or expired, clears the authentication data.
 *
 * @async
 * @returns {Promise<boolean>} Returns true if token is valid and user data is retrieved successfully,
 *                            false if token is invalid, missing, or the API request fails
 */
export async function validateToken() {
    const token = localStorage.getItem('token');

    if (!token) {
        return false;
    }

    try {
        // Try to get user data with the token
        const userData = await api.get('/user');

        // If we got a successful response, the token is valid
        if (userData && userData.data) {
            return true;
        } else {
            // Invalid token response
            console.warn('Token validation failed: Invalid response', userData);
            clearAuthData();
            return false;
        }
    } catch (error) {
        // API error means token is invalid or expired
        console.warn('Token validation failed:', error);
        clearAuthData();
        return false;
    }
}

// Helper function to clear auth data
export function clearAuthData() {
    localStorage.removeItem('token');
    // Any other cached user data you might have
}


// Function to fetch user data
async function fetchUserData() {
    const token = localStorage.getItem('token');
    if (!token) return null;

    try {
        const userData = await api.get('/user'); // Adjust the endpoint as necessary
        // console.log("USERDATA",userData)
        return userData;
    } catch (error) {
        console.error('Error fetching user data:', error);
        return null;
    }
}

// Update the UI based on authentication state
export async function updateAuthUI() {
    const authContainer = document.getElementById('auth-container');
    if (!authContainer) return;

    const token = localStorage.getItem('token');
    let userData = null;

    if (token) {
        userData = await fetchUserData();
        if (userData) {
            // User is logged in
            // On prépare le HTML avec Welcome, cloche, logout
            authContainer.innerHTML = `
                <span class="welcome-user-starwars">Welcome, ${userData.data.nickname}</span>
                <button id="notification-button" class="nav-button notification-bell">🔔</button>
                <button id="logout-button" class="nav-button">Logout <img src="/static/images/star-wars-rebels.svg" alt="logout" style="height:1.2em;width:1.2em;margin-left:8px;vertical-align:middle;filter:invert(1) brightness(2);"></button>
            `;

            // Glow rouge si notification non lue
            if (window.hasNotification) {
                document.getElementById('notification-button').classList.add('glow-red');
            }

            // Ajout du listener logout
            document.getElementById('logout-button').addEventListener('click', () => {
                import('./websocket.js').then(module => {
                    module.closeWebSocket();
                    localStorage.removeItem('token');
                    if (window.currentUser) window.currentUser = null;
                    updateAuthUI();
                    router.navigate('/');
                });
            });

            // Listener notification (redirige vers chat)
            document.getElementById('notification-button').addEventListener('click', () => {
                window.hasNotification = false;
                document.getElementById('notification-button').classList.remove('glow-red');
                window.hideNotificationBadge && window.hideNotificationBadge();
                router.navigate('/chat');
            });

            // Méthode globale pour activer le glow
            window.showNotificationGlow = function () {
                window.hasNotification = true;
                document.getElementById('notification-button')?.classList.add('glow-red');
            };
            window.hideNotificationGlow = function () {
                window.hasNotification = false;
                document.getElementById('notification-button')?.classList.remove('glow-red');
            };
        }
    } else {
        // User is not logged in
        authContainer.innerHTML = `
            <button id="login-button" class="nav-button">Login <img src="/static/images/star-wars-rebels.svg" alt="login" style="height:1.2em;width:1.2em;margin-left:8px;vertical-align:middle;filter:invert(1) brightness(2);"></button>
            <button id="register-button" class="nav-button">Register</button>
        `;
        document.getElementById('login-button').addEventListener('click', () => {
            router.navigate('/login');
        });
        document.getElementById('register-button').addEventListener('click', () => {
            router.navigate('/register');
        });
    }
}

/**
 * API module for making HTTP requests.
 * @module api
 */
export const api = {
    /**
     * The base URL for the API.
     * @type {string}
     */
    baseUrl: 'http://localhost:8080/api',

    /**
     * Makes a GET request to the specified endpoint.
     *
     * @async
     * @function get
     * @param {string} endpoint - The API endpoint to send the GET request to.
     * @returns {Promise<Object>} - A promise that resolves to the JSON response from the API.
     * @throws {Error} - Throws an error if the request fails.
     *
     * @example
     * api.get('/users')
     *   .then(data => console.log(data))
     *   .catch(error => console.error(error));
     */
    async get(endpoint) {
        const token = localStorage.getItem('token');
        const headers = {
            'Content-Type': 'application/json'
        };

        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }

        try {
            const response = await fetch(`${this.baseUrl}${endpoint}`, {
                method: 'GET',
                headers
            });

            return await response.json();
        } catch (error) {
            console.error('API Error:', error);
            console.log("RESPONSE: ", response);
            throw error;
        }
    },

    /**
     * Makes a POST request to the specified endpoint with the provided data.
     *
     * @async
     * @function post
     * @param {string} endpoint - The API endpoint to send the POST request to.
     * @param {Object} data - The data to be sent in the body of the POST request.
     * @returns {Promise<Object>} - A promise that resolves to the JSON response from the API.
     * @throws {Error} - Throws an error if the request fails.
     *
     * @example
     * api.post('/users', { name: 'John Doe' })
     *   .then(data => console.log(data))
     *   .catch(error => console.error(error));
     */
    async post(endpoint, data) {
        const token = localStorage.getItem('token');
        // console.log("<<<<<TOKEN", token);
        const headers = {
            'Content-Type': 'application/json'
        };

        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }

        try {
            const response = await fetch(`${this.baseUrl}${endpoint}`, {
                method: 'POST',
                headers,
                body: JSON.stringify(data)
            });

            return await response.json();
        } catch (error) {
            console.error('API Error:', error);
            throw error;
        }
    }
};

// Ajoute l'animation CSS pour le toast
// const style = document.createElement('style');
// style.innerHTML = `
// @keyframes toastFadeIn { from { opacity: 0; transform: translateY(-20px);} to { opacity: 1; transform: translateY(0);} }
// @keyframes toastFadeOut { from { opacity: 1; } to { opacity: 0; transform: translateY(-20px);} }
// `;
// document.head.appendChild(style);
