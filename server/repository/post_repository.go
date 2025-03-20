package repository

import (
	"database/sql"
	"real-time-forum/models"
)

// PostRepository provides methods to interact with post data
type PostRepository struct {
	DB *sql.DB
}

// NewPostRepository creates a new PostRepository
func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{DB: db}
}

// CreatePost inserts a new post into the database
func (r *PostRepository) CreatePost(post models.Post) (int64, error) {
	query := `
        INSERT INTO posts (title, category, text_content, user_id, date_creation)
        VALUES (?, ?, ?, ?, ?)`

	result, err := r.DB.Exec(query, post.Title, post.Category, post.TextContent, post.UserID, post.DateCreation)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// GetAllPosts retrieves all posts from the database
func (r *PostRepository) GetAllPosts() ([]models.Post, error) {
	query := `
        SELECT p.post_id, p.title, p.category, p.text_content, p.date_creation, p.user_id, p.likes, p.dislikes, u.nickname
        FROM posts p
        JOIN users u ON p.user_id = u.user_id
        ORDER BY p.date_creation DESC`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Category,
			&post.TextContent,
			&post.DateCreation,
			&post.UserID,
			&post.Likes,
			&post.Dislikes,
			&post.UserNickname,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

// GetPostByID retrieves a post by its ID
func (r *PostRepository) GetPostByID(id int) (models.Post, error) {
	var post models.Post
	query := `
        SELECT p.post_id, p.title, p.category, p.text_content, p.date_creation, p.user_id, p.likes, p.dislikes, u.nickname
        FROM posts p
        JOIN users u ON p.user_id = u.user_id
        WHERE p.post_id = ?`

	err := r.DB.QueryRow(query, id).Scan(
		&post.ID,
		&post.Title,
		&post.Category,
		&post.TextContent,
		&post.DateCreation,
		&post.UserID,
		&post.Likes,
		&post.Dislikes,
		&post.UserNickname,
	)

	if err != nil {
		return post, err
	}

	return post, nil
}
