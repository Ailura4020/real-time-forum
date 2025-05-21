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
    window.showToast = function(message) {
        const toastContainer = document.getElementById('toast-container');
        if (!toastContainer) return;
        const toast = document.createElement('div');
        toast.className = 'toast-notif-red';
        toast.innerHTML = `<span style="font-weight:bold;">${message}</span>`;
        toast.style.background = 'rgba(255,0,0,0.95)';
        toast.style.color = '#fff';
        toast.style.padding = '1rem 2rem';
        toast.style.marginBottom = '1rem';
        toast.style.borderRadius = '8px';
        toast.style.boxShadow = '0 0 16px 4px #ff0000, 0 0 32px 8px #ff0000';
        toast.style.fontSize = '1.1rem';
        toast.style.letterSpacing = '1px';
        toast.style.display = 'flex';
        toast.style.alignItems = 'center';
        toast.style.animation = 'toastFadeIn 0.3s';
        toastContainer.appendChild(toast);
        // Lecture du son
        const audio = new Audio('/static/sounds/lightsaber3.mp3');
        audio.play();
        // Disparition auto
        setTimeout(() => {
            toast.style.animation = 'toastFadeOut 0.5s';
            setTimeout(() => toast.remove(), 500);
        }, 3500);
    };

    // Create main content container
    const mainContent = document.createElement('main');
    mainContent.id = 'main-content';
    appElement.appendChild(mainContent);

    // Initialize router
    router.init();

    // Check if user is already logged in
    const token = localStorage.getItem('token');
    // const userData = localStorage.getItem('userData');

    // if (token && userData) {
    //     // Update UI for logged in user
    //     const userDataObj = JSON.parse(userData);
    //     updateAuthUI(userDataObj);
    //
    // }

    if (token) {
        // Fetch user data
        try {
            const userData = await api.get('/user'); // Assuming you have an endpoint to get user data
            updateAuthUI(userData);
        } catch (error) {
            console.error('Failed to fetch user data:', error);
        }

    }

});

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

// Ajoute l'animation CSS pour le toast
const style = document.createElement('style');
style.innerHTML = `
@keyframes toastFadeIn { from { opacity: 0; transform: translateY(-20px);} to { opacity: 1; transform: translateY(0);} }
@keyframes toastFadeOut { from { opacity: 1; } to { opacity: 0; transform: translateY(-20px);} }
`;
document.head.appendChild(style);

