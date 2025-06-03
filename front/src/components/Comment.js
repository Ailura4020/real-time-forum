// src/components/Comment.js
import { api } from '../main.js';

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
        // Use date_creation for consistency with the rest of the app
        const dateA = new Date(b.date_creation || b.create_date);
        const dateB = new Date(a.date_creation || a.create_date);
        return dateA - dateB;
    });

    // Create comment elements
    comments.forEach(comment => {
        const commentElement = document.createElement('article');
        commentElement.className = 'comment';

        // Format date - handle both date property naming conventions
        const commentDate = new Date(comment.date_creation || comment.create_date);
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
        const token = localStorage.getItem('token');
        if (!token) {
            alert('You must be logged in to comment');
            return;
        }
        
        const response = await fetch('http://localhost:8080/api/comments', {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                text_content: content,
                post_id: parseInt(postId, 10)
            })
        });

        console.log('Sending comment:', {
            text_content: content,
            post_id: postId
        });

        if (!response.ok) {
            const errorResponse = await response.json();
            console.error('Error response:', errorResponse);
            throw new Error('Network response was not ok: ' + (errorResponse.message || response.statusText));
        }

        const result = await response.json();
        console.log('Comment submitted successfully:', result);
        // alert('Comment submitted successfully!');

        // Reset form
        contentTextarea.value = '';

        // Reload comments
        const commentsListContainer = document.getElementById('comments-list');
        await loadComments(postId, commentsListContainer);
    } catch (error) {
        console.error('Error submitting comment:', error);
        alert('Failed to submit comment. Please try again.');
    }
}

