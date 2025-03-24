// src/components/Post.js
import { api } from '../main.js';
import { renderComments } from './Comment.js';
// import '../styles/Post.module.css';

export async function renderPostPage(container, params) {
    const postId = params.id;

    // Create a loading indicator
    container.innerHTML = '<div class="loading">Loading post...</div>';

    try {
        // Fetch the specific post with its comments
        const postResponse = await api.get(`/posts/${postId}`);

        if (postResponse.success) {
            const { post, comments } = postResponse.data;

            // Clear loading indicator
            container.innerHTML = '';

            // Create post container
            const postContainer = document.createElement('div');
            postContainer.className = 'single-post-container';

            // Format date
            const postDate = new Date(post.date_creation);
            const formattedDate = postDate.toLocaleDateString('en-US', {
                year: 'numeric',
                month: 'short',
                day: 'numeric',
                hour: '2-digit',
                minute: '2-digit'
            });

            // Render post details
            postContainer.innerHTML = `
        <div class="post-header">
          <h1 class="post-title">${post.title}</h1>
          <div class="post-meta">
            <span class="post-category">${post.category}</span>
            <span class="post-author">Posted by: ${post.user_nickname}</span>
            <span class="post-date">${formattedDate}</span>
          </div>
        </div>
        <div class="post-content">
          <p>${post.text_content}</p>
        </div>
        <!--<div class="post-actions">
          <button class="like-button" data-id="${post.id}">
            <span class="like-icon">👍</span> <span class="like-count">${post.likes}</span>
          </button>
          <button class="dislike-button" data-id="${post.id}">
            <span class="dislike-icon">👎</span> <span class="dislike-count">${post.dislikes}</span>
          </button>-->
        </div>
        <div class="post-navigation">
          <a href="/" class="back-link">← Back to Posts</a>
        </div>
      `;

            container.appendChild(postContainer);

            // Add comments section
            const commentsContainer = document.createElement('div');
            commentsContainer.className = 'comments-container';
            commentsContainer.innerHTML = '<h2>Comments</h2>';
            container.appendChild(commentsContainer);

            // Render comments
            renderComments(commentsContainer, postId, comments);

            // Add event listeners
            setTimeout(() => {
                // Like button
                container.querySelector('.like-button').addEventListener('click', async () => {
                    if (!localStorage.getItem('token')) {
                        alert('Please log in to like posts');
                        return;
                    }

                    try {
                        // This is a placeholder - you would need to implement an API endpoint for this
                        console.log('Like post:', post.id);
                        // Update UI
                        const likeCount = container.querySelector('.like-count');
                        likeCount.textContent = parseInt(likeCount.textContent) + 1;
                    } catch (error) {
                        console.error('Error liking post:', error);
                    }
                });

                // Dislike button
                container.querySelector('.dislike-button').addEventListener('click', async () => {
                    if (!localStorage.getItem('token')) {
                        alert('Please log in to dislike posts');
                        return;
                    }

                    try {
                        // This is a placeholder - you would need to implement an API endpoint for this
                        console.log('Dislike post:', post.id);
                        // Update UI
                        const dislikeCount = container.querySelector('.dislike-count');
                        dislikeCount.textContent = parseInt(dislikeCount.textContent) + 1;
                    } catch (error) {
                        console.error('Error disliking post:', error);
                    }
                });
            }, 100);
        } else {
            container.innerHTML = `<div class="error-container">
        <h2>Error Loading Post</h2>
        <p>${postResponse.message || 'Failed to load post'}</p>
        <a href="/" class="back-link">← Back to Posts</a>
      </div>`;
        }
    } catch (error) {
        container.innerHTML = `<div class="error-container">
      <h2>Error Loading Post</h2>
      <p>Failed to connect to the server. Please try again later.</p>
      <a href="/" class="back-link">← Back to Posts</a>
    </div>`;
        console.error('Error fetching post:', error);
    }
}