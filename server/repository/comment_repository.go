package repository

import (
	"database/sql"
	"real-time-forum/models"
)

// CommentRepository provides methods to interact with comment data
type CommentRepository struct {
	DB *sql.DB
}

// NewCommentRepository creates a new CommentRepository
func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{DB: db}
}

// CreateComment inserts a new comment into the database
func (r *CommentRepository) CreateComment(comment models.Comment) (int64, error) {
	query := `
        INSERT INTO comments (text_content, user_id, post_id, create_date)
        VALUES (?, ?, ?, ?)`

	result, err := r.DB.Exec(query, comment.TextContent, comment.UserID, comment.PostID, comment.CreateDate)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// GetCommentsByPostID retrieves all comments for a specific post
func (r *CommentRepository) GetCommentsByPostID(postID int) ([]models.Comment, error) {
	query := `
        SELECT c.comments_id, c.text_content, c.create_date, c.user_id, c.post_id, c.likes, c.dislikes, u.nickname
        FROM comments c
        JOIN users u ON c.user_id = u.user_id
        WHERE c.post_id = ?
        ORDER BY c.create_date ASC`

	rows, err := r.DB.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.TextContent,
			&comment.CreateDate,
			&comment.UserID,
			&comment.PostID,
			&comment.Likes,
			&comment.Dislikes,
			&comment.UserNickname,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
