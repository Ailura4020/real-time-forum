document.addEventListener('DOMContentLoaded', () => {
    const registerForm = document.querySelector('#register form');
    const loginForm = document.querySelector('#login form');
    const registerMessage = document.getElementById('registerMessage');
    const loginMessage = document.getElementById('loginMessage');
    const messageContainer = document.getElementById('messageContainer');
    const messageInput = document.getElementById('messageInput');
    const sendButton = document.getElementById('sendButton');
    const userInfoContainer = document.getElementById('userInfoContainer');

    // Chat container elements - assuming you want to hide/show these based on auth status
    const chatSection = document.getElementById('chatSection');
    const loginSection = document.getElementById('login');
    const registerSection = document.getElementById('register');

    // Cache user data in memory during the session, not in localStorage
    let currentUser = null;

    // Handle registration
    registerForm.addEventListener('submit', async (event) => {
        event.preventDefault();
        clearMessages(registerMessage);

        const formData = new FormData(registerForm);
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

            registerMessage.textContent = result.message;
            registerMessage.style.color = 'green';

            if (result.success) {
                // Store token and minimal user data
                storeUserSession(result.token);
                // Set current user in memory
                currentUser = result.data;
                // Display user info
                displayUserInfo(result.data);
                // Initialize WebSocket with the token
                initializeWebSocket(result.token);
                // Reset form
                registerForm.reset();
            }
        } catch (error) {
            console.error('Error during registration:', error);
            registerMessage.textContent = error.message || 'Registration failed. Please try again.';
            registerMessage.style.color = 'red';
        }
    });

    // Handle login
    loginForm.addEventListener('submit', async (event) => {
        event.preventDefault();
        clearMessages(loginMessage);

        const formData = new FormData(loginForm);
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

            const result = await response.json();

            if (!response.ok) {
                throw new Error(result.message || `HTTP error! status: ${response.status}`);
            }

            loginMessage.textContent = result.message;
            loginMessage.style.color = 'green';

            if (result.success) {
                // Store token only
                storeUserSession(result.token);
                // Set current user in memory
                currentUser = result.data;
                // Display user info
                displayUserInfo(result.data);
                // Initialize WebSocket with the token
                initializeWebSocket(result.token);
                // Reset form
                loginForm.reset();
            }
        } catch (error) {
            console.error('Error during login:', error);
            loginMessage.textContent = error.message || 'Login failed. Please try again.';
            loginMessage.style.color = 'red';
        }
    });

    // Store only the token in local storage
    const storeUserSession = (token) => {
        localStorage.setItem('token', token);
    };

    // Helper to get token
    const getToken = () => {
        return localStorage.getItem('token');
    };

    // Display user information
    const displayUserInfo = (userData) => {
        userInfoContainer.innerHTML = `
            <p>Welcome, ${userData.nickname}!</p>
            <p>Email: ${userData.email}</p>
            <p>Age: ${userData.age}</p>
            <p>Gender: ${userData.gender}</p>
            <button id="logoutButton">Logout</button>
        `;

        // Add logout button functionality
        document.getElementById('logoutButton').addEventListener('click', logout);
    };

    // Handle logout
    const logout = () => {
        // Close the WebSocket connection
        if (socket) {
            socket.close();
        }

        // Clear local storage and memory
        localStorage.removeItem('token');
        currentUser = null;

        // Reset UI
        userInfoContainer.innerHTML = '';
        messageContainer.innerHTML = '';

        // Show login/register forms
        showAuthInterface();
    };

    // Show chat interface, hide login/register
    const showChatInterface = () => {
        if (loginSection) loginSection.style.display = 'none';
        if (registerSection) registerSection.style.display = 'none';
    };

    // Show login/register interface, hide chat
    const showAuthInterface = () => {
        if (loginSection) loginSection.style.display = 'block';
        if (registerSection) registerSection.style.display = 'block';
    };

    // Fetch user information using the token and check user status
    const checkUserStatus = async (token) => {
        if (!token) {
            showAuthInterface();
            return;
        }

        try {
            const response = await fetch('http://localhost:8080/api/user', {
                method: 'GET',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                }
            });

            console.log('[CURRENT USER]',currentUser);

            const result = await response.json();

            if (!response.ok) {
                if (response.status === 401 || response.status === 403) {
                    // Token is invalid or expired
                    localStorage.removeItem('token');
                    currentUser = null;
                    userInfoContainer.textContent = 'Your session has expired. Please log in again.';
                    userInfoContainer.style.color = 'red';
                    showAuthInterface();
                } else {
                    throw new Error(result.message || `HTTP error! status: ${response.status}`);
                }
            } else {
                // Store user data in memory, not localStorage
                currentUser = result.data;
                displayUserInfo(result.data);
                // console.log('CURRENT USER:\n',currentUser);
            }
        } catch (error) {
            console.error('Error fetching user info:', error);
            showAuthInterface();
        }
    };

    // Initialize: Check token and set up the app accordingly
    const initialize = () => {
        const token = getToken();
        if (token) {
            // If we have a token, validate it
            checkUserStatus(token);
        } else {
            // No token, show login/register
            showAuthInterface();
        }
    };

    // Start the app
    initialize();
});
