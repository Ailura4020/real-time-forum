import { router } from './router.js';
import { renderNavigation } from './components/Navigation.js';
import { connectWebSocket } from './websocket.js';
// import './styles/Main.module.css';

// document.addEventListener('DOMContentLoaded', async () => {
//     const appElement = document.getElementById('app');
//
//     // Render navigation bar
//     const navElement = renderNavigation();
//     appElement.appendChild(navElement);
//
//     // Create main content container
//     const mainContent = document.createElement('main');
//     mainContent.id = 'main-content';
//     appElement.appendChild(mainContent);
//
//     // Initialize router
//     router.init();
//
//     // Check if user is already logged in
//     const token = localStorage.getItem('token');
//     // const userData = localStorage.getItem('userData');
//
//     // if (token && userData) {
//     //     // Update UI for logged in user
//     //     const userDataObj = JSON.parse(userData);
//     //     updateAuthUI(userDataObj);
//     //
//     // }
//
//     if (token) {
//         // Fetch user data
//         try {
//             const userData = await api.get('/user'); // Assuming you have an endpoint to get user data
//             updateAuthUI(userData);
//         } catch (error) {
//             console.error('Failed to fetch user data:', error);
//         }
//
//     }
//
// });


document.addEventListener('DOMContentLoaded', async () => {
    const appElement = document.getElementById('app');

    // Render navigation bar
    const navElement = renderNavigation();
    appElement.appendChild(navElement);

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
    console.log(token);
    
    let userData = null;

    if (token) {
        userData = await fetchUserData();
        if (userData) {
            // console.log("[SUCCESS]",userData, userData.data.nickname);
            // connectWebSocket(token)
            // User is logged in
            authContainer.innerHTML = `
                <span class="welcome-message">Welcome, ${userData.data.nickname}</span>
                <button id="logout-button" class="nav-button">Logout</button>
            `;

            // Add logout event listener
            document.getElementById('logout-button').addEventListener('click', () => {
                // Close WebSocket connection if it exists
                import('./websocket.js').then(module => {
                    module.closeWebSocket();
                    
                    // Clear user data
                    localStorage.removeItem('token');
                    
                    // Clear any chat-related state
                    if (window.currentUser) {
                        window.currentUser = null;
                    }
                    
                    // Update UI and navigate
                    updateAuthUI();
                    router.navigate('/');
                });
            });
        }
    } else {
        // User is not logged in
        authContainer.innerHTML = `
            <button id="login-button" class="nav-button">Login</button>
            <button id="register-button" class="nav-button">Register</button>
        `;

        // Add login/register event listeners
        document.getElementById('login-button').addEventListener('click', () => {
            router.navigate('/login');
        });

        document.getElementById('register-button').addEventListener('click', () => {
            router.navigate('/register');
        });
    }
}


// Create a simple API client
export const api = {
    baseUrl: 'http://localhost:8080/api',

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
            throw error;
        }
    },

    async post(endpoint, data) {
        const token = localStorage.getItem('token');
        console.log("<<<<<TOKEN",token)
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

