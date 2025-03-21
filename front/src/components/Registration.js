import { api, updateAuthUI } from '../main.js';
import { router } from '../router.js';
import '../styles/Registration.module.css';

export function renderRegistrationPage(container) {
    // Check if user is already logged in
    if (localStorage.getItem('token')) {
        container.innerHTML = `
      <div class="auth-container">
        <h1>Already Logged In</h1>
        <p>You are already registered and logged in.</p>
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
                updateAuthUI();
                renderRegistrationPage(container);
            });
        }, 0);

        return;
    }

    // Render registration form
    container.innerHTML = `
    <div class="auth-form-container">
      <h1>Register</h1>
      <form id="register-form" class="auth-form">
        <div class="form-row">
          <div class="form-group">
            <label for="first-name">First Name</label>
            <input type="text" id="first-name" name="first_name" required>
          </div>
          <div class="form-group">
            <label for="last-name">Last Name</label>
            <input type="text" id="last-name" name="last_name" required>
          </div>
        </div>
        <div class="form-group">
          <label for="nickname">Nickname</label>
          <input type="text" id="nickname" name="nickname" required>
        </div>
        <div class="form-row">
          <div class="form-group">
            <label for="age">Age</label>
            <input type="number" id="age" name="age" min="13" max="120" required>
          </div>
          <div class="form-group">
            <label for="gender">Gender</label>
            <select id="gender" name="gender" required>
              <option value="">Select...</option>
              <option value="male">Male</option>
              <option value="female">Female</option>
              <option value="other">Other</option>
              <option value="prefer-not-to-say">Prefer not to say</option>
            </select>
          </div>
        </div>
        <div class="form-group">
          <label for="email">Email</label>
          <input type="email" id="email" name="email" required>
        </div>
        <div class="form-group">
          <label for="password">Password</label>
          <input type="password" id="password" name="password" minlength="8" required>
        </div>
        <div class="form-group">
          <label for="confirm-password">Confirm Password</label>
          <input type="password" id="confirm-password" name="confirm_password" minlength="8" required>
        </div>
        <div id="register-error" class="error-message hidden"></div>
        <div class="form-actions">
          <button type="submit" class="btn btn-primary">Register</button>
        </div>
      </form>
      <p class="auth-link">Already have an account? <a href="/login" id="login-link">Login</a></p>
    </div>
  `;

    // Add event listeners
    setTimeout(() => {
        document.getElementById('register-form').addEventListener('submit', handleRegistration);

        document.getElementById('login-link').addEventListener('click', (e) => {
            e.preventDefault();
            router.navigate('/login');
        });
    }, 0);
}

async function handleRegistration(e) {
    e.preventDefault();

    const nickname = document.getElementById('nickname').value;
    const age = parseInt(document.getElementById('age').value);
    const gender = document.getElementById('gender').value;
    const firstName = document.getElementById('first-name').value;
    const lastName = document.getElementById('last-name').value;
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    const confirmPassword = document.getElementById('confirm-password').value;
    const errorElement = document.getElementById('register-error');

    // Reset error message
    errorElement.classList.add('hidden');
    errorElement.textContent = '';

    // Validate inputs
    if (!nickname || !age || !gender || !firstName || !lastName || !email || !password) {
        errorElement.textContent = 'Please fill in all fields';
        errorElement.classList.remove('hidden');
        return;
    }

    if (password !== confirmPassword) {
        errorElement.textContent = 'Passwords do not match';
        errorElement.classList.remove('hidden');
        return;
    }

    // Create current date in UTC format
    const dateRegister = new Date().toISOString();

    // Prepare registration data
    const registrationData = {
        nickname,
        age,
        gender,
        first_name: firstName,
        last_name: lastName,
        email,
        password,
        date_register: dateRegister
    };

    try {
        // Send registration request
        const response = await api.post('/register', registrationData);

        if (response.success) {
            // Store token and user data
            localStorage.setItem('token', response.token);
            localStorage.setItem('userData', JSON.stringify(response.data));

            // Update UI
            updateAuthUI(response.data);

            // Redirect to home page
            router.navigate('/');
        } else {
            // Show error message
            errorElement.textContent = response.message || 'Registration failed';
            errorElement.classList.remove('hidden');
        }
    } catch (error) {
        console.error('Registration error:', error);
        errorElement.textContent = 'Connection error. Please try again later.';
        errorElement.classList.remove('hidden');
    }
}