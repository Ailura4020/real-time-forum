export const chatTemplate = (classes) => `
  <div id="user-info" class="${classes.userInfo}"></div>
  <div class="chat">
    <h1>Chat</h1>
    <div class="${classes.chatLayout}">
      <div class="${classes.chatSidebar}">
        <div class="${classes.sidebarSection}">
          <h3>Online Users</h3>
          <div id="connected-users" class="${classes.usersList}"></div>
        </div>
         <div class="${classes.sidebarSection}">
          <h3>Recent Conversations</h3>
          <div id="recent-conversations" class="${classes.conversationsList}">
            <!-- Les conversations seront ajoutées dynamiquement ici -->
            <!-- Exemple de conversation -->
            <div id="conversation-1" class="conversationItem user-notification">
              <span class="nickname">User Nickname</span>
              <span class="lastMessagePreview">Last message preview...</span>
            </div>
          </div>
        </div>
      </div>
      <div class="${classes.chatContent}">
        <div id="chat-recipient" class="${classes.chatRecipient}">
          <span class="recipient-label">Select a user to start a chat</span>
        </div>
        <div class="${classes.chatMessages}" id="messages"></div>
        <div class="${classes.chatInput}">
          <input type="text" id="messageInput" placeholder="Type your message...">
          <button id="send-button" class="${classes.sendButton}">Send</button>
        </div>
      </div>
    </div>
  </div>
`;

export const postFormTemplate = (classes) => `
  <h2>Create New Post</h2>
  <form id="post-form" class="${classes.postForm}">
    <div class="${classes.formGroup}">
      <label for=post-title>Title</label>
      <input type="text" id="post-title" name="title" required>
    </div>
    <div class="${classes.formGroup}">
      <label for="post-category">Category</label>
      <select id="post-category" name="category" required>
        <option value="General">General</option>
        <option value="Question">Question</option>
        <option value="Discussion">Discussion</option>
        <option value="Announcement">Announcement</option>
      </select>
    </div>
    <div class="${classes.formGroup}">
      <label for="post-content">Content</label>
      <textarea id="post-content" name="content" rows="6" required></textarea>
    </div>
    <div class="${classes.formActions}">
      <button type="button" id="cancel-post">Cancel</button>
      <button type="submit">Submit Post</button>
    </div>
  </form>
`;

export const postTemplate = (classes, post, formattedDate) => `
<div class="post-header">
    <h2 class="${classes.postTitle}">
        <a href="/post/${post.id}">${post.title}</a>
    </h2>
    <span class="post-category">${post.category}</span>
</div>
<div class="post-content">
    <p class="post-content">${post.text_content}</p>
</div>
<div class="post-footer">
    <div class="post-meta">
        <span class="post-author">Posted by: ${post.user_nickname}</span>
        <span class="post-date">• ${formattedDate}</span>
    </div>
</div>
`;

// ARCH

// src/templates.js
// This file contains reusable HTML templates for the application

// export const templates = {
//     // Error message template
//     errorMessage: (message) => `
//     <div class="error-message">
//       <span class="error-icon">⚠️</span>
//       <span class="error-text">${message}</span>
//     </div>
//   `,
//
//     // Success message template
//     successMessage: (message) => `
//     <div class="success-message">
//       <span class="success-icon">✅</span>
//       <span class="success-text">${message}</span>
//     </div>
//   `,
//
//     // Loading indicator template
//     loadingIndicator: (message = 'Loading...') => `
//     <div class="loading-indicator">
//       <div class="spinner"></div>
//       <span class="loading-text">${message}</span>
//     </div>
//   `,
//
//     // Empty state template
//     emptyState: (message, actionText = null, actionUrl = null) => `
//     <div class="empty-state">
//       <div class="empty-icon">📭</div>
//       <p class="empty-message">${message}</p>
//       ${actionText && actionUrl ? `<a href="${actionUrl}" class="empty-action">${actionText}</a>` : ''}
//     </div>
//   `,
//
//     // Confirmation dialog template
//     confirmDialog: (title, message, confirmText = 'Confirm', cancelText = 'Cancel') => `
//     <div class="confirm-dialog">
//       <h3 class="dialog-title">${title}</h3>
//       <p class="dialog-message">${message}</p>
//       <div class="dialog-actions">
//         <button class="btn btn-secondary dialog-cancel">${cancelText}</button>
//         <button class="btn btn-primary dialog-confirm">${confirmText}</button>
//       </div>
//     </div>
//   `
// };
//
// // Helper function to show a toast notification
// export function showToast(message, type = 'info', duration = 3000) {
//     // Create toast container if it doesn't exist
//     let toastContainer = document.getElementById('toast-container');
//     if (!toastContainer) {
//         toastContainer = document.createElement('div');
//         toastContainer.id = 'toast-container';
//         document.body.appendChild(toastContainer);
//     }
//
//     // Create toast element
//     const toast = document.createElement('div');
//     toast.className = `toast toast-${type}`;
//     toast.textContent = message;
//
//     // Add to container
//     toastContainer.appendChild(toast);
//
//     // Trigger animation
//     setTimeout(() => {
//         toast.classList.add('show');
//     }, 10);
//
//     // Remove after duration
//     setTimeout(() => {
//         toast.classList.remove('show');
//         setTimeout(() => {
//             toast.remove();
//         }, 300);
//     }, duration);
// }
