import { api, updateAuthUI, setCookie, deleteCookie } from '../main.js';
import { router } from '../router.js';
import { connectWebSocket } from '../websocket.js';
import '../styles/Login.module.css';

export function renderLoginPage(container) {
    // Check if user is already logged in
    if (localStorage.getItem('token')) {
        container.innerHTML = `
      <div class="auth-container">
        <h1>Already Logged In</h1>
        <p>You are already logged in.</p>
        <button id="go-home" class="btn btn-primary">Go to Home</button>
        <button id="logout" class="btn btn-secondary">Logout</button>
      </div>
    `;

        // Add event listeners
        setTimeout(() => {
            document.getElementById('go-home').addEventListener('click', () => {
                router.navigate('/');
            });

            document.getElementById('logout').addEventListener('click', () => {
                localStorage.removeItem('token');
                localStorage.removeItem('userData');
                deleteCookie('jwt_token');
                updateAuthUI();
                renderLoginPage(container);
            });
        }, 0);

        return;
    }

    // Render login form
    container.innerHTML = `
    <div class="auth-form-container">
      <h1>Login</h1>
      <form id="login-form" class="auth-form">
        <div class="form-group">
          <label for="email">Email</label>
          <input type="email" id="email" name="email" required>
        </div>
        <div class="form-group">
          <label for="password">Password</label>
          <input type="password" id="password" name="password" required>
        </div>
        <div id="login-error" class="error-message hidden"></div>
        <div class="form-actions">
          <button type="submit" class="btn btn-primary">Login</button>
        </div>
      </form>
      <p class="auth-link">Don't have an account? <a href="/register" id="register-link">Register</a></p>
    </div>
  `;

    // Add event listeners
    setTimeout(() => {
        document.getElementById('login-form').addEventListener('submit', handleLogin);

        document.getElementById('register-link').addEventListener('click', (e) => {
            e.preventDefault();
            router.navigate('/register');
        });
    }, 0);
}

async function handleLogin(e) {
    e.preventDefault();

    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    const errorElement = document.getElementById('login-error');

    // Reset error message
    errorElement.classList.add('hidden');
    errorElement.textContent = '';

    // Validate inputs
    if (!email || !password) {
        errorElement.textContent = 'Please fill in all fields';
        errorElement.classList.remove('hidden');
        return;
    }

    try {
        // Send login request
        const response = await api.post('/login', { email, password });

        if (response.success) {
            // Store token and user data
            localStorage.setItem('token', response.token);
            setCookie('jwt_token', response.token, 7); // Store JWT in cookie for 7 days
            // localStorage.setItem('userData', JSON.stringify(response.data)); // no need

            // Initiate WebSocket connection
            connectWebSocket(response.token);

            // Update UI
            updateAuthUI(response.data);

            // Redirect to home page
            // router.navigate('/');
            // Check if there's a redirect path stored
            const redirectPath = sessionStorage.getItem('redirectAfterLogin');
            if (redirectPath) {
                sessionStorage.removeItem('redirectAfterLogin'); // Clear it after use
                router.navigate(redirectPath);
            } else {
                // Default redirect to home page
                router.navigate('/');
            }

        } else {
            // Show error message
            errorElement.textContent = response.message || 'Login failed';
            errorElement.classList.remove('hidden');
        }
    } catch (error) {
        console.error('Login error:', error);
        errorElement.textContent = 'Connection error. Please try again later.';
        errorElement.classList.remove('hidden');
    }
}
