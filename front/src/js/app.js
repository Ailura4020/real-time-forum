// Main JavaScript file for the forum application
document.addEventListener('DOMContentLoaded', function() {
    // API URL - change this to match your server
    const API_URL = 'http://localhost:8080';

    console.log(API_URL);

    // DOM Elements
    const registrationForm = document.getElementById('registration-form');
    const loginForm = document.getElementById('login-form');
    // const socialButtons = document.querySelectorAll('.social-login button');

    // Check if we have a token from social login
    // const checkForToken = () => {
    //     // Look for token in URL hash
    //     if (window.location.hash.includes('login-success')) {
    //         const urlParams = new URLSearchParams(window.location.hash.replace('#/', ''));
    //         const token = urlParams.get('token');
    //
    //         if (token) {
    //             // Store the token
    //             localStorage.setItem('auth_token', token);
    //             // Remove the token from the URL
    //             window.history.replaceState({}, document.title, '/');
    //
    //             // Show success message and redirect to a protected page
    //             alert('Login successful!');
    //             // Optionally redirect to a dashboard here
    //             // window.location.href = '/dashboard.html';
    //         }
    //     }
    // };

    // Handle Registration
    if (registrationForm) {
        registrationForm.addEventListener('submit', function(e) {
            e.preventDefault();

            // Get form values
            const nickname = this.querySelector('input[placeholder="Nickname"]').value;
            const age = parseInt(this.querySelector('input[placeholder="Age"]').value, 10);
            const gender = this.querySelector('select').value;
            const firstName = this.querySelector('input[placeholder="First Name"]').value;
            const lastName = this.querySelector('input[placeholder="Last Name"]').value;
            const email = this.querySelector('input[placeholder="E-mail"]').value;
            const password = this.querySelector('input[placeholder="Password"]').value;

            // Validate form
            if (!nickname || !age || !gender || !firstName || !lastName || !email || !password) {
                alert('Please fill in all fields');
                return;
            }

            // Create request data
            const userData = {
                nickname,
                age,
                gender,
                first_name: firstName,
                last_name: lastName,
                email,
                password,
                date_register: new Date().toISOString() // Add current date
            };

            // Send registration request
            fetch(`${API_URL}/api/register`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(userData)
            })
                .then(response => response.json())
                .then(data => {
                    if (data.success) {
                        // Store token
                        if (data.token) {
                            localStorage.setItem('auth_token', data.token);
                        }

                        // Clear form
                        registrationForm.reset();

                        // Show success message
                        alert('Registration successful! You are now logged in.');

                        // Optionally redirect to another page
                        // window.location.href = '/dashboard.html';
                    } else {
                        alert(data.message || 'Registration failed. Please try again.');
                    }
                })
                .catch(error => {
                    console.error('Registration error:', error);
                    alert('An error occurred during registration. Please try again.');
                });
        });
    }

    // Handle Login
    if (loginForm) {
        loginForm.addEventListener('submit', function(e) {
            e.preventDefault();

            // Get form values
            const email = this.querySelector('input[placeholder="E-mail"]').value;
            const password = this.querySelector('input[placeholder="Password"]').value;

            // Validate form
            if (!email || !password) {
                alert('Please enter both email and password');
                return;
            }

            // Create request data
            const loginData = {
                email,
                password
            };

            // Send login request
            fetch(`${API_URL}/api/login`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(loginData)
            })
                .then(response => response.json())
                .then(data => {
                    if (data.success) {
                        // Store token
                        if (data.token) {
                            localStorage.setItem('auth_token', data.token);
                        }

                        // Clear form
                        loginForm.reset();

                        // Show success message
                        alert('Login successful!');

                        // Optionally redirect to another page
                        // window.location.href = '/dashboard.html';
                    } else {
                        alert(data.message || 'Login failed. Please check your credentials.');
                    }
                })
                .catch(error => {
                    console.error('Login error:', error);
                    alert('An error occurred during login. Please try again.');
                });
        });
    }

    // Handle Social Login Buttons
    // Uncomment and implement if you have social login functionality
    // if (socialButtons) {
    //     socialButtons.forEach(button => {
    //         button.addEventListener('click', function() {
    //             const provider = this.textContent.trim().toLowerCase();
    //             window.location.href = `${API_URL}/auth/${provider}`;
    //         });
    //     });
    // }

    // Function to make authenticated requests
    const makeAuthenticatedRequest = (url, method = 'GET', data = null) => {
        const token = localStorage.getItem('auth_token');

        if (!token) {
            console.error('No authentication token found');
            return Promise.reject('Not authenticated');
        }

        const options = {
            method,
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json'
            }
        };

        if (data && (method === 'POST' || method === 'PUT')) {
            options.body = JSON.stringify(data);
        }

        return fetch(url, options)
            .then(response => {
                if (response.status === 401) {
                    // Token expired or invalid
                    localStorage.removeItem('auth_token');
                    throw new Error('Authentication expired. Please log in again.');
                }
                return response.json();
            });
    };

    // Example of accessing a protected route
    const accessProtectedRoute = () => {
        makeAuthenticatedRequest(`${API_URL}/api/protected`)
            .then(data => {
                console.log('Protected data:', data);
            })
            .catch(error => {
                console.error('Error accessing protected route:', error);
                // Redirect to login if needed
                // window.location.href = '/login.html';
            });
    };

    // Check for token from social login redirect
    // checkForToken();

    // Check if user is logged in
    const isLoggedIn = () => {
        return localStorage.getItem('auth_token') !== null;
    };

    // Logout function
    const logout = () => {
        localStorage.removeItem('auth_token');
        // Redirect to login page
        // window.location.href = '/login.html';
        alert('You have been logged out.');
    };

    // Optional: Add logout button event listener
    const logoutButton = document.getElementById('logout-button');
    if (logoutButton) {
        logoutButton.addEventListener('click', logout);
    }

    // Optional: Update UI based on login status
    const updateUI = () => {
        const loggedIn = isLoggedIn();

        // Elements to show when logged in
        document.querySelectorAll('.logged-in').forEach(el => {
            el.style.display = loggedIn ? 'block' : 'none';
        });

        // Elements to show when logged out
        document.querySelectorAll('.logged-out').forEach(el => {
            el.style.display = loggedIn ? 'none' : 'block';
        });
    };

    // Run initial UI update
    updateUI();
});
