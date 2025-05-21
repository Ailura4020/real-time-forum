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
        <div id="auth-container" class="auth-container"></div>
      </div>
    </div>
  `;

  return nav;
}
