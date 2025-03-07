document.addEventListener('DOMContentLoaded', () => {
    const registerForm = document.querySelector('#register form');
    const loginForm = document.querySelector('#login form');
    const registerMessage = document.getElementById('registerMessage');
    const loginMessage = document.getElementById('loginMessage');
    const messageContainer = document.getElementById('messageContainer');
    const messageInput = document.getElementById('messageInput');
    const sendButton = document.getElementById('sendButton');

    // WebSocket connection
    const socket = new WebSocket('ws://localhost:8080/ws');

    // Handle WebSocket connection open
    socket.addEventListener('open', () => {
        console.log('WebSocket connection established');
    });

    // Handle incoming messages
    socket.addEventListener('message', (event) => {
        const message = document.createElement('div');
        message.textContent = event.data; // Display the incoming message
        messageContainer.appendChild(message);
    });

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
            date_register: new Date().toISOString() // Use current date in ISO format
        };

        try {
            const response = await fetch('http://localhost:8080/api/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(data)
            });

            // Check if the response is OK (status in the range 200-299)
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const result = await response.json();
            registerMessage.textContent = result.message; // Show success or error message
            registerMessage.style.color = result.success ? 'green' : 'red';
        } catch (error) {
            console.error('Error during registration:', error);
            registerMessage.textContent = 'Registration failed. Please try again.';
            registerMessage.style.color = 'red';
        }
    });

    // Handle login
    loginForm.addEventListener('submit', async (event) => {
        event.preventDefault();
        clearMessages(loginMessage);

        const formData = new FormData(loginForm);
        const data = {
            email: formData.get('identifier'), // Assuming identifier is the email
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
            loginMessage.textContent = result.message; // Show success or error message
            loginMessage.style.color = result.success ? 'green' : 'red';
        } catch (error) {
            console.error('Error during login:', error);
            loginMessage.textContent = 'Login failed. Please try again.';
            loginMessage.style.color = 'red';
        }
    });

    // // Handle sending messages via WebSocket
    sendButton.addEventListener('click', () => {
        const message = messageInput.value;
        if (message) {
            socket.send(message); // Send the message to the WebSocket server
            messageInput.value = ''; // Clear the input field
        }else{
            console.error('Websocket is not open. Message not sent.')
        }
    });
    
    // // Optional: Handle pressing Enter to send messages
    messageInput.addEventListener('keypress', (event) => {
        if (event.key === 'Enter') {
            sendButton.click(); // Trigger the send button click
        }
    });
    

    // UI:Clear messages function
    const clearMessages = (messageElement) => {
        messageElement.textContent = '';
    };
});

socket.addEventListener('close', () => {
    console.log('Websocket connection closed');
})
