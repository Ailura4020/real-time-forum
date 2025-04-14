import {api} from '../main.js';
import classes from '../styles/Home.module.css';
import {postFormTemplate, postTemplate} from '../templates.js';
// import "../styles/Home.module.css"
// import styles from '../styles/Home.module.css';

export async function renderHomePage(container) {
    // Create a loading indicator
    container.innerHTML = '<div class="${classes.loading}">Loading posts...</div>';

    try {
      // // connexion websockets pour afficher la liste des users connectés
      // connectWebSocket(usersList => {
      //   // console.log(userContainer)
      //   updateConnectedUsers(userContainer, usersList);
      // });

        // Fetch posts from the API
        const postsResponse = await api.get('/posts');
        if (postsResponse.success) {
            const posts = postsResponse.data;

            // Clear loading indicator
            container.innerHTML = '';

            // Create posts container
            const postsContainer = document.createElement('div');
            postsContainer.className = 'posts-container';

            // Add heading
            const heading = document.createElement('h1');
            heading.textContent = 'Recent Posts';
            postsContainer.appendChild(heading);

            // Check if user is logged in to show create post button
            if (localStorage.getItem('token')) {
                const createPostButton = document.createElement('button');
                createPostButton.className = 'create-post-button';
                createPostButton.textContent = 'Create New Post';
                createPostButton.addEventListener('click', showCreatePostForm);
                postsContainer.appendChild(createPostButton);

                // Create post form (hidden initially)
                const createPostForm = document.createElement('div');
                createPostForm.className = 'create-post-form hidden';
                createPostForm.id = 'create-post-form';
                createPostForm.innerHTML = postFormTemplate(classes);

                postsContainer.appendChild(createPostForm);

                // Add event listeners for post form
                setTimeout(() => {
                    document.getElementById('cancel-post')?.addEventListener('click', () => {
                        document.getElementById('create-post-form').classList.add('hidden');
                    });

                    document.getElementById('post-form')?.addEventListener('submit', async (e) => {
                        e.preventDefault();
                        await createPost();
                    });
                }, 100);
            }

            // Add posts list
            const postsList = document.createElement('div');
            postsList.className = 'posts-list';

            if (posts.length === 0) {
                postsList.innerHTML = '<p class="no-posts">No posts found.</p>';
            } else {
                // Sort posts by date (newest first)
                posts.sort((a, b) => {
                    return new Date(b.date_creation) - new Date(a.date_creation);
                });

                posts.forEach(post => {
                    const postElement = createPostElement(post);
                    postsList.appendChild(postElement);
                });
            }

            postsContainer.appendChild(postsList);
            container.appendChild(postsContainer);
        } else {
            container.innerHTML = `<div class="error-container">
        <h2>Error Loading Posts</h2>
        <p>${postsResponse.message || 'Failed to load posts'}</p>
      </div>`;
        }
    } catch (error) {
        container.innerHTML = `<div class="error-container">
      <h2>Error Loading Posts</h2>
      <p>Failed to connect to the server. Please try again later.</p>
    </div>`;
        console.error('Error fetching posts:', error);
    }
}

function createPostElement(post) {
    const postElement = document.createElement('article');
    postElement.className = 'post-card';
    postElement.dataset.postId = post.id;

    // Format date
    const postDate = new Date(post.date_creation);
    const formattedDate = postDate.toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    });

    // template
    postElement.innerHTML = postTemplate(classes, post, formattedDate);

    return postElement;
}

function showCreatePostForm() {
    document.getElementById('create-post-form').classList.remove('hidden');
}

async function createPost() {
    const titleInput = document.getElementById('post-title');
    const categorySelect = document.getElementById('post-category');
    const contentTextarea = document.getElementById('post-content');

    const postData = {
        title: titleInput.value,
        category: categorySelect.value,
        text_content: contentTextarea.value // Ensure this matches your API's expected field name
    };

    if (!postData.title.trim() || !postData.category.trim() || !postData.text_content.trim()) {
        alert('Please fill in all fields');
        return;
    }

    try {
        const token = localStorage.getItem('token'); // Get the JWT token from local storage

        console.log("[DATA]",postData)
        console.log("[JSON]",JSON.stringify(postData))

        const response = await fetch('http://localhost:8080/api/posts', {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(postData) // Send the post data in the request body
        });

        if (!response.ok) {
            throw new Error('Network response was not ok ' + response.statusText);
        }

        const result = await response.json();
        console.log('Post created successfully:', result);
        // alert('Post created successfully!');

        // Reset form and hide it
        document.getElementById('post-form').reset();
        document.getElementById('create-post-form').classList.add('hidden');

        // Reload posts
        const container = document.getElementById('main-content');
        await renderHomePage(container);
    } catch (error) {
        console.error('Error creating post:', error);
        alert('Failed to create post. Please try again.');
    }
}
