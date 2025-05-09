export function renderNavigation() {
    const nav = document.createElement('nav');
    nav.classList.add('main-nav');

    nav.innerHTML = `
    <div class="nav-content">
      <div class="logo-container">
        <a href="/" class="logo">Forum App</a>
      </div>
      <div class="nav-links">
        <a href="/" class="nav-link">Home</a>
        <a href="/chat" class="nav-link">Chat</a>
        <a href="/about" class="nav-link">About</a>
      </div>
     <div class="right-section">
  <button id="notification-button" class="nav-button">
    🔔 <span id="notification-badge" style="display: none; color: red;">●</span>
  </button>
  <div id="auth-container" class="auth-container">

        <button id="login-button" class="nav-button">Login</button>
        <button id="register-button" class="nav-button">Register</button>
      </div>
    </div>
  `;
  const notificationButton = nav.querySelector('#notification-button');
  const notificationBadge = nav.querySelector('#notification-badge');
  
  window.showNotificationBadge = function () {
    if (notificationBadge) {
      notificationBadge.style.display = 'inline';
    }
  };
  
  window.hideNotificationBadge = function () {
    if (notificationBadge) {
      notificationBadge.style.display = 'none';
    }
  };
  
  notificationButton.addEventListener('click', () => {
    console.log('[Notification] Cloche cliquée');
    window.hideNotificationBadge();
    window.history.pushState({}, '', '/chat');
    window.dispatchEvent(new PopStateEvent('popstate'));
  });
    // Add event listeners
    const loginButton = nav.querySelector('#login-button');
    const registerButton = nav.querySelector('#register-button');

    loginButton.addEventListener('click', () => {
        window.history.pushState({}, '', '/login');
        const event = new PopStateEvent('popstate');
        window.dispatchEvent(event);
    });

    registerButton.addEventListener('click', () => {
        window.history.pushState({}, '', '/register');
        const event = new PopStateEvent('popstate');
        window.dispatchEvent(event);
    });

    return nav;
}

