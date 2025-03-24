// src/pages/About.js
// import '../styles/About.module.css';

export function renderAboutPage(container) {
    container.innerHTML = `
    <div class="about-container">
      <h1>About Our Forum</h1>
      <div class="about-content">
        <section class="about-section">
          <h2>Welcome to Our Community</h2>
          <p>This forum is a place for users to share ideas, ask questions, and engage in meaningful discussions.</p>
          <p>Our platform is built with vanilla JavaScript and a Go backend to provide a fast and responsive experience.</p>
        </section>
        
        <section class="about-section">
          <h2>Features</h2>
          <ul>
            <li>Create and browse posts on various topics</li>
            <li>Comment on posts and engage with other users</li>
            <li>Like or dislike content to show your opinion</li>
            <li>User accounts with secure authentication</li>
          </ul>
        </section>
        
        <section class="about-section">
          <h2>Rules & Guidelines</h2>
          <ol>
            <li>Be respectful to other members</li>
            <li>No spamming or self-promotion</li>
            <li>Stay on topic in discussions</li>
            <li>No offensive or inappropriate content</li>
            <li>Have fun and enjoy the community!</li>
          </ol>
        </section>
      </div>
    </div>
  `;
}