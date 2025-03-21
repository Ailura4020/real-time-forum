// src/components/Comment.js
import { api } from '../main.js';
import '../styles/Comment.module.css';

export function renderComments(container, postId, existingComments = null) {
    // Add comment form if user is logged in
    if (localStorage.getItem('token')) {
        const commentFormContainer = document.createElement('div');
        commentFormContainer.className = 'comment-form-container';
        commentFormContainer.innerHTML = `
      <form id="comment-form" class="comment-form">
        <div class="form-group">
          <label for="comment-content">Add a comment</label>
          <textarea id="comment-content" name="content" rows="3" required></textarea>
        </div>
        <div class="form-actions">
          <button type="submit" class="btn btn-primary">Submit Comment</button>
        </div>
      </form>
    `;
        container.appendChild(commentFormContainer);

        // Add event listener for comment form
        setTimeout(() => {
            document.getElementById('comment-form').addEventListener('submit', async (e) => {
                e.preventDefault();
                await submitComment(postId, container);
            });
        }, 100);
    } else {
        const loginPrompt = document.createElement('div');
        loginPrompt.className = 'login-prompt';
        loginPrompt.innerHTML = `
      <p>Please <a href="/login" class="login-link">log in</a> to add a comment.</p>
    `;
        container.appendChild(loginPrompt);
    }

    // Create comments list container
    const commentsListContainer = document.createElement('div');
    commentsListContainer.className = 'comments-list';
    commentsListContainer.id = 'comments-list';
    container.appendChild(commentsListContainer);

    // Display comments or load them if not provided
    if (existingComments !== null) {
        displayComments(commentsListContainer, existingComments);
    } else {
        loadComments(postId, commentsListContainer);
    }
}

async function loadComments(postId, container) {
    container.innerHTML = '<p class="loading-comments">Loading comments...</p>';

    try {
        // Fetch post details which include comments
        const response = await api.get(`/posts/${postId}`);

        if (response.success && response.data.comments) {
            displayComments(container, response.data.comments);
        } else {
            container.innerHTML = '<p class="no-comments">No comments yet.</p>';
        }
    } catch (error) {
        console.error('Error loading comments:', error);
        container.innerHTML = '<p class="error">Failed to load comments. Please try again later.</p>';
    }
}

function displayComments(container, comments) {
    if (!comments || comments.length === 0) {
        container.innerHTML = '<p class="no-comments">No comments yet.</p>';
        return;
    }

    // Clear container
    container.innerHTML = '';

    // Sort comments by date (newest first)
    comments.sort((a, b) => {
        return new Date(b.date_creation) - new Date(a.date_creation);
    });

    // Create comment elements
    comments.forEach(comment => {
        const commentElement = document.createElement('div');
        commentElement.className = 'comment';

        // Format date
        const commentDate = new Date(comment.date_creation);
        const formattedDate = commentDate.toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });

        commentElement.innerHTML = `
      <div class="comment-header">
        <span class="comment-author">${comment.user_nickname}</span>
        <span class="comment-date">${formattedDate}</span>
      </div>
      <div class="comment-content">
        <p>${comment.text_content}</p>
      </div>
    `;

        container.appendChild(commentElement);
    });
}

async function submitComment(postId, parentContainer) {
    const contentTextarea = document.getElementById('comment-content');
    const content = contentTextarea.value;

    if (!content.trim()) {
        alert('Please enter a comment');
        return;
    }

    try {
        // This is a placeholder - you would need to implement an API endpoint for this
        console.log('Submitting comment for post:', postId, 'content:', content);
        alert('Comment submission not implemented in this demo. Would send: ' + content);

        // Reset form
        contentTextarea.value = '';

        // Reload comments
        const commentsListContainer = document.getElementById('comments-list');
        loadComments(postId, commentsListContainer);
    } catch (error) {
        console.error('Error submitting comment:', error);
        alert('Failed to submit comment. Please try again.');
    }
}