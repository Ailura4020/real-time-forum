I have a web app for handling an api written in golang here is the js code

```js
document.addEventListener('DOMContentLoaded', () => {
const registerForm = document.querySelector('#register form');
const loginForm = document.querySelector('#login form');
const registerMessage = document.getElementById('registerMessage');
const loginMessage = document.getElementById('loginMessage');
const messageContainer = document.getElementById('messageContainer');
const messageInput = document.getElementById('messageInput');
const sendButton = document.getElementById('sendButton');
const userInfoContainer = document.getElementById('userInfoContainer');

    // Chat container elements - assuming you want to hide/show these based on auth status    const chatSection = document.getElementById('chatSection');
    const loginSection = document.getElementById('login');
    const registerSection = document.getElementById('register');

    // Cache user data in memory during the session, not in localStorage    let currentUser = null;
    let socket = null;

    // Initialize WebSocket connection with authentication    const initializeWebSocket = (token) => {
        // Close existing connection if any        if (socket) {
            socket.close();
        }

        socket = new WebSocket(`ws://localhost:8080/ws?token=${token}`);

        socket.addEventListener('open', () => {
            console.log('WebSocket connection established');
            // Show chat UI when connection is established            showChatInterface();
        });

        socket.addEventListener('message', (event) => {
            try {
                // Try to parse as JSON first                const data = JSON.parse(event.data);
                displayFormattedMessage(data);
            } catch (e) {
                // If not JSON, display as plain text                const message = document.createElement('div');
                message.textContent = event.data;
                messageContainer.appendChild(message);
                // Auto-scroll to bottom                messageContainer.scrollTop = messageContainer.scrollHeight;
            }
        });

        socket.addEventListener('close', () => {
            console.log('WebSocket connection closed');
        });

        socket.addEventListener('error', (error) => {
            console.error('WebSocket error:', error);
        });
    };

    // Display formatted messages    const displayFormattedMessage = (data) => {
        const message = document.createElement('div');
        message.className = 'message';

        // Format based on message type        if (data.user) {
            message.innerHTML = `<strong>${data.user}</strong>: ${data.content}`;
            message.className += currentUser && data.user === currentUser.nickname ? ' own-message' : ' other-message';
        } else {
            message.textContent = data.content || data;
        }

        messageContainer.appendChild(message);
        // Auto-scroll to bottom        messageContainer.scrollTop = messageContainer.scrollHeight;
    };

    // Handle registration    registerForm.addEventListener('submit', async (event) => {
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
                // Store token and minimal user data                storeUserSession(result.token);
                // Set current user in memory                currentUser = result.data;
                // Display user info                displayUserInfo(result.data);
                // Initialize WebSocket with the token                initializeWebSocket(result.token);
                // Reset form                registerForm.reset();
            }
        } catch (error) {
            console.error('Error during registration:', error);
            registerMessage.textContent = error.message || 'Registration failed. Please try again.';
            registerMessage.style.color = 'red';
        }
    });

    // Handle login    loginForm.addEventListener('submit', async (event) => {
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
                    'Content-Type': 'application/json'                },
                body: JSON.stringify(data)
            });

            const result = await response.json();

            if (!response.ok) {
                throw new Error(result.message || `HTTP error! status: ${response.status}`);
            }

            loginMessage.textContent = result.message;
            loginMessage.style.color = 'green';

            if (result.success) {
                // Store token only                storeUserSession(result.token);
                // Set current user in memory                currentUser = result.data;
                // Display user info                displayUserInfo(result.data);
                // Initialize WebSocket with the token                initializeWebSocket(result.token);
                // Reset form                loginForm.reset();
            }
        } catch (error) {
            console.error('Error during login:', error);
            loginMessage.textContent = error.message || 'Login failed. Please try again.';
            loginMessage.style.color = 'red';
        }
    });

    // Handle sending messages via WebSocket    sendButton.addEventListener('click', () => {
        sendMessage();
    });

    // Handle pressing Enter to send messages    messageInput.addEventListener('keypress', (event) => {
        if (event.key === 'Enter') {
            event.preventDefault(); // Prevent default to avoid form submission            sendMessage();
        }
    });

    // Extract sending logic to a separate function    const sendMessage = () => {
        const message = messageInput.value.trim();
        if (message && socket && socket.readyState === WebSocket.OPEN && currentUser) {
            // Send as JSON with user info            const messageData = {
                content: message,
                user: currentUser.nickname,
                userId: currentUser.id            };

            socket.send(JSON.stringify(messageData));
            messageInput.value = '';
        } else if (!socket || socket.readyState !== WebSocket.OPEN) {
            console.error('WebSocket is not open. Message not sent.');
            // Try to reconnect            const token = getToken();
            if (token) {
                initializeWebSocket(token);
            }
        } else if (!currentUser) {
            console.error('User information not available. Message not sent.');
            checkUserStatus(getToken());
        }
    };

    // UI: Clear messages function    const clearMessages = (messageElement) => {
        messageElement.textContent = '';
    };

    // Store only the token in local storage    const storeUserSession = (token) => {
        localStorage.setItem('token', token);
    };

    // Helper to get token    const getToken = () => {
        return localStorage.getItem('token');
    };

    // Display user information    const displayUserInfo = (userData) => {
        userInfoContainer.innerHTML = `            <p>Welcome, ${userData.nickname}!</p>            <p>Email: ${userData.email}</p>            <p>Age: ${userData.age}</p>            <p>Gender: ${userData.gender}</p>            <button id="logoutButton">Logout</button>        `;

        // Add logout button functionality        document.getElementById('logoutButton').addEventListener('click', logout);
    };

    // Handle logout    const logout = () => {
        // Close the WebSocket connection        if (socket) {
            socket.close();
        }

        // Clear local storage and memory        localStorage.removeItem('token');
        currentUser = null;

        // Reset UI        userInfoContainer.innerHTML = '';
        messageContainer.innerHTML = '';

        // Show login/register forms        showAuthInterface();
    };

    // Show chat interface, hide login/register    const showChatInterface = () => {
        if (chatSection) chatSection.style.display = 'block';
        if (loginSection) loginSection.style.display = 'none';
        if (registerSection) registerSection.style.display = 'none';
    };

    // Show login/register interface, hide chat    const showAuthInterface = () => {
        if (chatSection) chatSection.style.display = 'none';
        if (loginSection) loginSection.style.display = 'block';
        if (registerSection) registerSection.style.display = 'block';
    };

    // Fetch user information using the token and check user status    const checkUserStatus = async (token) => {
        if (!token) {
            showAuthInterface();
            return;
        }

        try {
            const response = await fetch('http://localhost:8080/api/user', {
                method: 'GET',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'                }
            });
            document.addEventListener('DOMContentLoaded', () => {
                const registerForm = document.querySelector('#register form');
                const loginForm = document.querySelector('#login form');
                const registerMessage = document.getElementById('registerMessage');
                const loginMessage = document.getElementById('loginMessage');
                const messageContainer = document.getElementById('messageContainer');
                const messageInput = document.getElementById('messageInput');
                const sendButton = document.getElementById('sendButton');
                const userInfoContainer = document.getElementById('userInfoContainer');

                // Chat container elements - assuming you want to hide/show these based on auth status                const chatSection = document.getElementById('chatSection');
                const loginSection = document.getElementById('login');
                const registerSection = document.getElementById('register');

                // Cache user data in memory during the session, not in localStorage                let currentUser = null;
                let socket = null;

                // Initialize WebSocket connection with authentication                const initializeWebSocket = (token) => {
                    // Close existing connection if any                    if (socket) {
                        socket.close();
                    }

                    socket = new WebSocket(`ws://localhost:8080/ws?token=${token}`);

                    socket.addEventListener('open', () => {
                        console.log('WebSocket connection established');
                        // Show chat UI when connection is established                        showChatInterface();
                    });

                    socket.addEventListener('message', (event) => {
                        try {
                            // Try to parse as JSON first                            const data = JSON.parse(event.data);
                            displayFormattedMessage(data);
                        } catch (e) {
                            // If not JSON, display as plain text                            const message = document.createElement('div');
                            message.textContent = event.data;
                            messageContainer.appendChild(message);
                            // Auto-scroll to bottom                            messageContainer.scrollTop = messageContainer.scrollHeight;
                        }
                    });

                    socket.addEventListener('close', () => {
                        console.log('WebSocket connection closed');
                    });

                    socket.addEventListener('error', (error) => {
                        console.error('WebSocket error:', error);
                    });
                };

                // Display formatted messages                const displayFormattedMessage = (data) => {
                    const message = document.createElement('div');
                    message.className = 'message';

                    // Format based on message type                    if (data.user) {
                        message.innerHTML = `<strong>${data.user}</strong>: ${data.content}`;
                        message.className += currentUser && data.user === currentUser.nickname ? ' own-message' : ' other-message';
                    } else {
                        message.textContent = data.content || data;
                    }

                    messageContainer.appendChild(message);
                    // Auto-scroll to bottom                    messageContainer.scrollTop = messageContainer.scrollHeight;
                };

                // Handle registration                registerForm.addEventListener('submit', async (event) => {
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
                            // Store token and minimal user data                            storeUserSession(result.token);
                            // Set current user in memory                            currentUser = result.data;
                            // Display user info                            displayUserInfo(result.data);
                            // Initialize WebSocket with the token                            initializeWebSocket(result.token);
                            // Reset form                            registerForm.reset();
                        }
                    } catch (error) {
                        console.error('Error during registration:', error);
                        registerMessage.textContent = error.message || 'Registration failed. Please try again.';
                        registerMessage.style.color = 'red';
                    }
                });

                // Handle login                loginForm.addEventListener('submit', async (event) => {
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
                                'Content-Type': 'application/json'                            },
                            body: JSON.stringify(data)
                        });

                        const result = await response.json();

                        if (!response.ok) {
                            throw new Error(result.message || `HTTP error! status: ${response.status}`);
                        }

                        loginMessage.textContent = result.message;
                        loginMessage.style.color = 'green';

                        if (result.success) {
                            // Store token only                            storeUserSession(result.token);
                            // Set current user in memory                            currentUser = result.data;
                            // Display user info                            displayUserInfo(result.data);
                            // Initialize WebSocket with the token                            initializeWebSocket(result.token);
                            // Reset form                            loginForm.reset();
                        }
                    } catch (error) {
                        console.error('Error during login:', error);
                        loginMessage.textContent = error.message || 'Login failed. Please try again.';
                        loginMessage.style.color = 'red';
                    }
                });

                // Handle sending messages via WebSocket                sendButton.addEventListener('click', () => {
                    sendMessage();
                });

                // Handle pressing Enter to send messages                messageInput.addEventListener('keypress', (event) => {
                    if (event.key === 'Enter') {
                        event.preventDefault(); // Prevent default to avoid form submission                        sendMessage();
                    }
                });

                // Extract sending logic to a separate function                const sendMessage = () => {
                    const message = messageInput.value.trim();
                    if (message && socket && socket.readyState === WebSocket.OPEN && currentUser) {
                        // Send as JSON with user info                        const messageData = {
                            content: message,
                            user: currentUser.nickname,
                            userId: currentUser.id                        };

                        socket.send(JSON.stringify(messageData));
                        messageInput.value = '';
                    } else if (!socket || socket.readyState !== WebSocket.OPEN) {
                        console.error('WebSocket is not open. Message not sent.');
                        // Try to reconnect                        const token = getToken();
                        if (token) {
                            initializeWebSocket(token);
                        }
                    } else if (!currentUser) {
                        console.error('User information not available. Message not sent.');
                        checkUserStatus(getToken());
                    }
                };

                // UI: Clear messages function                const clearMessages = (messageElement) => {
                    messageElement.textContent = '';
                };

                // Store only the token in local storage                const storeUserSession = (token) => {
                    localStorage.setItem('token', token);
                };

                // Helper to get token                const getToken = () => {
                    return localStorage.getItem('token');
                };

                // Display user information                const displayUserInfo = (userData) => {
                    userInfoContainer.innerHTML = `            <p>Welcome, ${userData.nickname}!</p>            <p>Email: ${userData.email}</p>            <p>Age: ${userData.age}</p>            <p>Gender: ${userData.gender}</p>            <button id="logoutButton">Logout</button>        `;

                    // Add logout button functionality                    document.getElementById('logoutButton').addEventListener('click', logout);
                };

                // Handle logout                const logout = () => {
                    // Close the WebSocket connection                    if (socket) {
                        socket.close();
                    }

                    // Clear local storage and memory                    localStorage.removeItem('token');
                    currentUser = null;

                    // Reset UI                    userInfoContainer.innerHTML = '';
                    messageContainer.innerHTML = '';

                    // Show login/register forms                    showAuthInterface();
                };

                // Show chat interface, hide login/register                const showChatInterface = () => {
                    if (chatSection) chatSection.style.display = 'block';
                    if (loginSection) loginSection.style.display = 'none';
                    if (registerSection) registerSection.style.display = 'none';
                };

                // Show login/register interface, hide chat                const showAuthInterface = () => {
                    if (chatSection) chatSection.style.display = 'none';
                    if (loginSection) loginSection.style.display = 'block';
                    if (registerSection) registerSection.style.display = 'block';
                };

                // Fetch user information using the token and check user status                const checkUserStatus = async (token) => {
                    if (!token) {
                        showAuthInterface();
                        return;
                    }

                    try {
                        const response = await fetch('http://localhost:8080/api/user', {
                            method: 'GET',
                            headers: {
                                'Authorization': `Bearer ${token}`,
                                'Content-Type': 'application/json'                            }
                        });

                        console.log('[CURRENT USER]',currentUser);

                        const result = await response.json();

                        if (!response.ok) {
                            if (response.status === 401 || response.status === 403) {
                                // Token is invalid or expired                                localStorage.removeItem('token');
                                currentUser = null;
                                userInfoContainer.textContent = 'Your session has expired. Please log in again.';
                                userInfoContainer.style.color = 'red';
                                showAuthInterface();
                            } else {
                                throw new Error(result.message || `HTTP error! status: ${response.status}`);
                            }
                        } else {
                            // Store user data in memory, not localStorage                            currentUser = result.data;
                            displayUserInfo(result.data);
                            // console.log('CURRENT USER:\n',currentUser);                            // Initialize WebSocket with the token                            initializeWebSocket(token);
                        }
                    } catch (error) {
                        console.error('Error fetching user info:', error);
                        showAuthInterface();
                    }
                };

                // Initialize: Check token and set up the app accordingly                const initialize = () => {
                    const token = getToken();
                    if (token) {
                        // If we have a token, validate it                        checkUserStatus(token);
                    } else {
                        // No token, show login/register                        showAuthInterface();
                    }
                };

                // Start the app                initialize();
            });
            console.log('[CURRENT USER]',currentUser);

            const result = await response.json();

            if (!response.ok) {
                if (response.status === 401 || response.status === 403) {
                    // Token is invalid or expired                    localStorage.removeItem('token');
                    currentUser = null;
                    userInfoContainer.textContent = 'Your session has expired. Please log in again.';
                    userInfoContainer.style.color = 'red';
                    showAuthInterface();
                } else {
                    throw new Error(result.message || `HTTP error! status: ${response.status}`);
                }
            } else {
                // Store user data in memory, not localStorage                currentUser = result.data;
                displayUserInfo(result.data);
                // console.log('CURRENT USER:\n',currentUser);                // Initialize WebSocket with the token                initializeWebSocket(token);
            }
        } catch (error) {
            console.error('Error fetching user info:', error);
            showAuthInterface();
        }
    };

    // Initialize: Check token and set up the app accordingly    const initialize = () => {
        const token = getToken();
        if (token) {
            // If we have a token, validate it            checkUserStatus(token);
        } else {
            // No token, show login/register            showAuthInterface();
        }
    };

    // Start the app    initialize();
});
```

## objective: i want to implement a better system based on modules (classes)

here is an example of architecture i want to implement:

```shell
├── index.html
├── package.json
├── public
│   ├── css
│   │   └── style.css
│   ├── images
│   └── js
│       └── main.js
├── src
│   ├── api
│   │   ├── apiClient.js
│   │   ├── authApi.js
│   │   └── messageApi.js
│   ├── components
│   │   ├── auth
│   │   ├── chat
│   │   │   ├── chatContainer.js
│   │   │   ├── messageInput.js
│   │   │   └── messageList.js
│   │   └── shared
│   │       ├── errorMessage.js
│   │       └── userInfo.js
│   ├── config
│   │   └── config.js
│   ├── services
│   │   └── websocket.js
│   ├── store
│   │   ├── events.js
│   │   └── state.js
│   └── utils
│       ├── dom.js
│       ├── storage.js
│       └── validation.js
```

## Rules

- i do not want you write the entire project, just the login, registration and authentification with the token
- do not take the real time chat into account for the moment - i also want to use a publisher-subscriber pattern
- i want to test the new component system before anything else
- if you change the architecture, give me a shell command to do so


## Restriction

- vanilla js only, no framework


Do you have question?