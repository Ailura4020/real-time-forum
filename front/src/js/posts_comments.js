// script.js

// Ajout de styles CSS
const style = document.createElement('style');
style.innerHTML = `
    body {
        font-family: 'Arial', sans-serif;
        background-color: #000;
        color: #fff;
        margin: 20px;
    }
    #post-form {
        margin-bottom: 20px;
        background-color: #222;
        padding: 15px;
        border-radius: 5px;
    }
    input, textarea {
        width: 100%;
        padding: 10px;
        margin: 5px 0;
        border: 1px solid #444;
        border-radius: 5px;
        background-color: #333;
        color: #fff;
    }
    button {
        background-color: #f39c12;
        color: #fff;
        border: none;
        padding: 10px;
        border-radius: 5px;
        cursor: pointer;
    }
    button:hover {
        background-color: #e67e22;
    }
    .post {
        border: 1px solid #444;
        padding: 10px;
        margin-bottom: 10px;
        background-color: #222;
        border-radius: 5px;
    }
    .comments {
        margin-top: 10px;
    }
    .comment {
        border-top: 1px solid #555;
        padding: 5px 0;
    }
`;
document.head.appendChild(style);

// Structure de données pour stocker les publications
let posts = [];

// Fonction pour rendre les publications
function renderPosts() {
    const postsContainer = document.getElementById('posts');
    postsContainer.innerHTML = ''; // Effacer les publications existantes

    posts.forEach((post, index) => {
        const postElement = document.createElement('div');
        postElement.className = 'post';
        postElement.innerHTML = `
            <h3>${post.title}</h3>
            <p>${post.content}</p>
            <button onclick="toggleComments(${index})">Commentaires (${post.comments.length})</button>
            <div class="comments" id="comments-${index}" style="display: none;">
                <textarea id="comment-input-${index}" placeholder="Ajouter un commentaire"></textarea>
                <button onclick="addComment(${index})">Soumettre le commentaire</button>
                <div class="comment-list" id="comment-list-${index}"></div>
            </div>
        `;
        postsContainer.appendChild(postElement);
    });
}

// Fonction pour créer une nouvelle publication
function createPost() {
    const title = document.getElementById('post-title').value;
    const content = document.getElementById('post-content').value;

    if (title && content) {
        const newPost = {
            title: title,
            content: content,
            comments: []
        };
        posts.push(newPost);
        renderPosts();
        document.getElementById('post-title').value = '';
        document.getElementById('post-content').value = '';
    } else {
        alert('Veuillez remplir les deux champs.');
    }
}

// Fonction pour basculer la visibilité des commentaires
function toggleComments(index) {
    const commentsDiv = document.getElementById(`comments-${index}`);
    commentsDiv.style.display = commentsDiv.style.display === 'none' ? 'block' : 'none';
    renderComments(index);
}

// Fonction pour rendre les commentaires d'une publication spécifique
function renderComments(index) {
    const commentList = document.getElementById(`comment-list-${index}`);
    commentList.innerHTML = ''; // Effacer les commentaires existants

    posts[index].comments.forEach(comment => {
        const commentElement = document.createElement('div');
        commentElement.className = '
            commentElement.className = 'comment';
        commentElement.innerText = comment;
        commentList.appendChild(commentElement);
    });
}

// Fonction pour ajouter un commentaire à une publication
function addComment(index) {
    const commentInput = document.getElementById(`comment-input-${index}`);
    const comment = commentInput.value;

    if (comment) {
        posts[index].comments.push(comment);
        renderComments(index);
        commentInput.value = ''; // Réinitialiser le champ de commentaire
    } else {
        alert('Veuillez entrer un commentaire.');
    }
}

// Écouteur d'événements pour créer une publication
document.getElementById('create-post').addEventListener('click', createPost);

// Rendu initial
renderPosts();
