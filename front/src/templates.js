// src/templates.js
// This file contains reusable HTML templates for the application

export const templates = {
    // Error message template
    errorMessage: (message) => `
    <div class="error-message">
      <span class="error-icon">⚠️</span>
      <span class="error-text">${message}</span>
    </div>
  `,

    // Success message template
    successMessage: (message) => `
    <div class="success-message">
      <span class="success-icon">✅</span>
      <span class="success-text">${message}</span>
    </div>
  `,

    // Loading indicator template
    loadingIndicator: (message = 'Loading...') => `
    <div class="loading-indicator">
      <div class="spinner"></div>
      <span class="loading-text">${message}</span>
    </div>
  `,

    // Empty state template
    emptyState: (message, actionText = null, actionUrl = null) => `
    <div class="empty-state">
      <div class="empty-icon">📭</div>
      <p class="empty-message">${message}</p>
      ${actionText && actionUrl ? `<a href="${actionUrl}" class="empty-action">${actionText}</a>` : ''}
    </div>
  `,

    // Confirmation dialog template
    confirmDialog: (title, message, confirmText = 'Confirm', cancelText = 'Cancel') => `
    <div class="confirm-dialog">
      <h3 class="dialog-title">${title}</h3>
      <p class="dialog-message">${message}</p>
      <div class="dialog-actions">
        <button class="btn btn-secondary dialog-cancel">${cancelText}</button>
        <button class="btn btn-primary dialog-confirm">${confirmText}</button>
      </div>
    </div>
  `
};

// Helper function to show a toast notification
export function showToast(message, type = 'info', duration = 3000) {
    // Create toast container if it doesn't exist
    let toastContainer = document.getElementById('toast-container');
    if (!toastContainer) {
        toastContainer = document.createElement('div');
        toastContainer.id = 'toast-container';
        document.body.appendChild(toastContainer);
    }

    // Create toast element
    const toast = document.createElement('div');
    toast.className = `toast toast-${type}`;
    toast.textContent = message;

    // Add to container
    toastContainer.appendChild(toast);

    // Trigger animation
    setTimeout(() => {
        toast.classList.add('show');
    }, 10);

    // Remove after duration
    setTimeout(() => {
        toast.classList.remove('show');
        setTimeout(() => {
            toast.remove();
        }, 300);
    }, duration);
}